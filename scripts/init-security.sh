#!/bin/bash
# Инициализация безопасности - ЗАПУСКАТЬ ПЕРВЫМ!

set -e

echo "🔐 ИНИЦИАЛИЗАЦИЯ ВОЕННОЙ ЗАЩИТЫ"
echo "================================"

# Цвета
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

check_error() {
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Ошибка!${NC}"
        exit 1
    fi
}

# 1. Создание директории для секретов
echo -e "\n${YELLOW}1. Создание секретов...${NC}"
mkdir -p secrets
chmod 700 secrets

# Генерация мастер-ключа (32 байта)
MASTER_KEY=$(openssl rand -base64 32)
echo -n $MASTER_KEY > secrets/master_key.txt
echo -e "${GREEN}✅ Мастер-ключ создан${NC}"

# Генерация JWT секрета
JWT_SECRET=$(openssl rand -base64 64)
echo -n $JWT_SECRET > secrets/jwt_secret.txt
echo -e "${GREEN}✅ JWT секрет создан${NC}"

# Генерация пароля PostgreSQL
POSTGRES_PASSWORD=$(openssl rand -base64 24 | tr -d "=+/" | cut -c1-24)
echo -n $POSTGRES_PASSWORD > secrets/postgres_password.txt
echo -e "${GREEN}✅ PostgreSQL пароль создан${NC}"

# Генерация пароля Redis
REDIS_PASSWORD=$(openssl rand -base64 24 | tr -d "=+/" | cut -c1-24)
echo -n $REDIS_PASSWORD > secrets/redis_password.txt
echo -e "${GREEN}✅ Redis пароль создан${NC}"

# Генерация ключа для аудита
AUDIT_KEY=$(openssl rand -base64 64)
echo -n $AUDIT_KEY > secrets/audit_key.txt
echo -e "${GREEN}✅ Ключ аудита создан${NC}"

# 2. Создание SSL сертификатов
echo -e "\n${YELLOW}2. Генерация SSL сертификатов...${NC}"
mkdir -p config/nginx/ssl

openssl req -x509 -nodes -days 365 -newkey rsa:4096 \
    -keyout config/nginx/ssl/privkey.pem \
    -out config/nginx/ssl/fullchain.pem \
    -subj "/C=RU/ST=Moscow/L=Moscow/O=MedicalBot/CN=localhost" \
    -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

check_error
echo -e "${GREEN}✅ SSL сертификаты созданы${NC}"

# 3. Настройка прав доступа
echo -e "\n${YELLOW}3. Настройка прав доступа...${NC}"
chmod 600 secrets/*.txt
chmod 600 config/nginx/ssl/*.pem

echo -e "${GREEN}✅ Права установлены${NC}"

# 4. Создание .env файла
echo -e "\n${YELLOW}4. Создание .env.production...${NC}"

cat > .env.production << EOF
# ==========================================
# PRODUCTION КОНФИГУРАЦИЯ
# ==========================================

# База данных
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=${POSTGRES_PASSWORD}
DB_NAME=medical_bot

# Redis
REDIS_HOST=redis
REDIS_PASSWORD=${REDIS_PASSWORD}

# JWT
JWT_SECRET=${JWT_SECRET}
JWT_EXPIRY=900

# Telegram
TELEGRAM_TOKEN=your-telegram-bot-token
WEBHOOK_URL=https://your-domain.com
WEBHOOK_SECRET=$(openssl rand -base64 32)

# Мониторинг
GRAFANA_PASSWORD=$(openssl rand -base64 16)

# Криптография
MASTER_KEY=${MASTER_KEY}

# AWS (опционально)
AWS_REGION=eu-central-1
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=

# SIEM (опционально)
SIEM_URL=
SIEM_API_KEY=
EOF

check_error
echo -e "${GREEN}✅ .env.production создан${NC}"

# 5. Создание бэкапа ключей
echo -e "\n${YELLOW}5. Создание резервной копии ключей...${NC}"

BACKUP_FILE="secrets-backup-$(date +%Y%m%d-%H%M%S).enc"
tar czf - secrets/ 2>/dev/null | \
    openssl enc -aes-256-cbc -salt -pbkdf2 \
    -pass pass:"$(openssl rand -base64 32)" \
    -out $BACKUP_FILE

echo -e "${GREEN}✅ Резервная копия создана: $BACKUP_FILE${NC}"
echo -e "${YELLOW}⚠️ ВАЖНО: Сохраните пароль от бэкапа в безопасном месте!${NC}"

# 6. Проверка безопасности
echo -e "\n${YELLOW}6. Проверка безопасности...${NC}"

# Проверка Docker security
docker info --format '{{.SecurityOptions}}' | grep -q "name=seccomp"
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Seccomp включен${NC}"
fi

# Итог
echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}✅ ИНИЦИАЛИЗАЦИЯ ЗАВЕРШЕНА!${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "\n📁 Секреты сохранены в директории: secrets/"
echo -e "📝 Production конфиг: .env.production"
echo -e "🔐 Бэкап ключей: $BACKUP_FILE"
echo -e "\n${YELLOW}⚠️ НЕОБХОДИМЫЕ ДЕЙСТВИЯ:${NC}"
echo "1. Отредактируйте .env.production, добавьте TELEGRAM_TOKEN"
echo "2. Сохраните пароль от бэкапа в менеджере паролей"
echo "3. Запустите: docker-compose -f deployments/docker-compose.yml up -d"
