cat > start.sh << 'EOF'
#!/bin/bash

# =====================================
# АВТОМАТИЧЕСКИЙ ЗАПУСК TELEGRAM-БОТА
# =====================================

# Цвета для красивого вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}=====================================${NC}"
echo -e "${BLUE}ЗАПУСК TELEGRAM-БОТА${NC}"
echo -e "${BLUE}=====================================${NC}"

# Функция проверки ошибок
check_error() {
    if [ $? -ne 0 ]; then
        echo -e "${RED}Ошибка!${NC}"
        exit 1
    fi
}

# =====================================
# Шаг 1: Проверка PostgreSQL
# =====================================
echo -e "\n${YELLOW}Шаг 1: Проверка PostgreSQL...${NC}"

if ! docker ps | grep -q "postgres"; then
    echo -e "${YELLOW}PostgreSQL не запущен. Запускаем...${NC}"
    docker-compose down 2>/dev/null
    docker-compose up -d postgres
    check_error
    echo -e "${GREEN}PostgreSQL запущен${NC}"
    echo -e "${YELLOW}Ожидание инициализации (5 сек)...${NC}"
    sleep 5
else
    echo -e "${GREEN}PostgreSQL уже запущен${NC}"
fi

# =====================================
# Шаг 2: Проверка .env
# =====================================
echo -e "\n${YELLOW}Шаг 2: Проверка .env...${NC}"

if [ ! -f .env ]; then
    echo -e "${RED}Файл .env не найден!${NC}"
    echo -e "${YELLOW}Создаю .env из .env.example...${NC}"
    cp .env.example .env
    check_error
    echo -e "${GREEN}Файл .env создан${NC}"
    echo -e "${RED}Нужно отредактировать .env и добавить TELEGRAM_TOKEN!${NC}"
    exit 1
fi

# Проверяем наличие TELEGRAM_TOKEN
if ! grep -q "TELEGRAM_TOKEN" .env; then
    echo -e "${RED}TELEGRAM_TOKEN не найден в .env!${NC}"
    echo -e "${YELLOW}Добавьте строку: TELEGRAM_TOKEN=ваш_токен_от_BotFather${NC}"
    exit 1
fi

# Загружаем токен для проверки
source .env
if [ -z "$TELEGRAM_TOKEN" ] || [ "$TELEGRAM_TOKEN" = "your-telegram-bot-token" ]; then
    echo -e "${RED}TELEGRAM_TOKEN не установлен или имеет значение по умолчанию!${NC}"
    exit 1
fi
echo -e "${GREEN}TELEGRAM_TOKEN найден${NC}"

# =====================================
# Шаг 3: Проверка зависимостей Go
# =====================================
echo -e "\n${YELLOW}Шаг 3: Проверка зависимостей Go...${NC}"

if [ ! -f go.mod ]; then
    echo -e "${RED}go.mod не найден!${NC}"
    exit 1
fi

# Проверяем наличие всех необходимых пакетов
echo -e "${YELLOW}Проверка и установка зависимостей...${NC}"

# Список всех необходимых пакетов
PACKAGES=(
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "github.com/joho/godotenv"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

# Проверяем и устанавливаем каждый пакет
for pkg in "${PACKAGES[@]}"; do
    if ! go list "$pkg" > /dev/null 2>&1; then
        echo -e "${YELLOW}Устанавливаю $pkg...${NC}"
        go get "$pkg"
        check_error
    fi
done

# Синхронизируем зависимости
echo -e "${YELLOW}Синхронизация зависимостей (go mod tidy)...${NC}"
go mod tidy
check_error
echo -e "${GREEN}Все зависимости установлены${NC}"

# =====================================
# Шаг 4: Исправление базы данных
# =====================================
echo -e "\n${YELLOW}Шаг 4: Исправление базы данных...${NC}"

# Проверяем, существует ли таблица telegram_users
TABLE_EXISTS=$(docker exec bot_max-postgres-1 psql -U postgres -d medical_bot -t -c "
SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'telegram_users');" 2>/dev/null | xargs)

if [ "$TABLE_EXISTS" = "t" ]; then
    echo -e "${GREEN}Таблица telegram_users существует${NC}"
    
    # === ИСПРАВЛЕНИЕ ТИПА user_id ===
    echo -e "${YELLOW}Проверяю тип поля user_id...${NC}"
    
    # Проверяем тип поля user_id
    USER_ID_TYPE=$(docker exec bot_max-postgres-1 psql -U postgres -d medical_bot -t -c "
    SELECT data_type FROM information_schema.columns 
    WHERE table_name = 'telegram_users' AND column_name = 'user_id';" 2>/dev/null | xargs)
    
    if [ "$USER_ID_TYPE" = "uuid" ]; then
        echo -e "${YELLOW}Поле user_id имеет тип uuid. Меняю на text...${NC}"
        
        # Удаляем foreign key и меняем тип
        docker exec bot_max-postgres-1 psql -U postgres -d medical_bot -c "
        BEGIN;
        ALTER TABLE telegram_users DROP CONSTRAINT IF EXISTS telegram_users_user_id_fkey;
        ALTER TABLE telegram_users ALTER COLUMN user_id TYPE text;
        COMMIT;
        " > /dev/null 2>&1
        
        echo -e "${GREEN}Поле user_id изменено на text${NC}"
    else
        echo -e "${GREEN}Поле user_id уже имеет тип text${NC}"
    fi
    
    # === ДОБАВЛЕНИЕ ПОЛЯ email ===
    echo -e "${YELLOW}Проверяю наличие поля email...${NC}"
    
    EMAIL_EXISTS=$(docker exec bot_max-postgres-1 psql -U postgres -d medical_bot -t -c "
    SELECT COUNT(*) FROM information_schema.columns 
    WHERE table_name = 'telegram_users' AND column_name = 'email';" 2>/dev/null | xargs)
    
    if [ "$EMAIL_EXISTS" -eq 0 ]; then
        echo -e "${YELLOW}Добавляю поле email...${NC}"
        docker exec bot_max-postgres-1 psql -U postgres -d medical_bot -c "
        ALTER TABLE telegram_users ADD COLUMN email VARCHAR(255);" > /dev/null 2>&1
        echo -e "${GREEN}Поле email добавлено${NC}"
    else
        echo -e "${GREEN}Поле email существует${NC}"
    fi
    
    # === ДОБАВЛЕНИЕ ДРУГИХ ПОЛЕЙ ===
    for field in language_code auth_token token_expires; do
        FIELD_EXISTS=$(docker exec bot_max-postgres-1 psql -U postgres -d medical_bot -t -c "
        SELECT COUNT(*) FROM information_schema.columns 
        WHERE table_name = 'telegram_users' AND column_name = '$field';" 2>/dev/null | xargs)
        
        if [ "$FIELD_EXISTS" -eq 0 ]; then
            echo -e "${YELLOW}Добавляю поле $field...${NC}"
            case $field in
                language_code) docker exec bot_max-postgres-1 psql -U postgres -d medical_bot -c "ALTER TABLE telegram_users ADD COLUMN language_code VARCHAR(10);" > /dev/null 2>&1 ;;
                auth_token) docker exec bot_max-postgres-1 psql -U postgres -d medical_bot -c "ALTER TABLE telegram_users ADD COLUMN auth_token VARCHAR(50);" > /dev/null 2>&1 ;;
                token_expires) docker exec bot_max-postgres-1 psql -U postgres -d medical_bot -c "ALTER TABLE telegram_users ADD COLUMN token_expires TIMESTAMP;" > /dev/null 2>&1 ;;
            esac
            echo -e "${GREEN}Поле $field добавлено${NC}"
        else
            echo -e "${GREEN}Поле $field существует${NC}"
        fi
    done
else
    echo -e "${YELLOW}Таблица telegram_users будет создана автоматически при запуске${NC}"
fi

# =====================================
# Шаг 5: Запуск бота
# =====================================
echo -e "\n${GREEN}=====================================${NC}"
echo -e "${GREEN}ВСЕ ПРОВЕРКИ ПРОЙДЕНЫ!${NC}"
echo -e "${GREEN}=====================================${NC}"
echo -e "${GREEN}Запуск бота...${NC}"
echo -e "${GREEN}=====================================${NC}"
echo -e "${YELLOW}Для остановки нажмите Ctrl+C${NC}\n"

go run cmd/telegram-bot/main.go

echo -e "\n${YELLOW}Бот остановлен${NC}"
EOF

chmod +x start.sh

echo -e "\n${GREEN}Скрипт start.sh успешно создан!${NC}"
echo -e "${YELLOW}Для запуска выполните: ./start.sh${NC}"
