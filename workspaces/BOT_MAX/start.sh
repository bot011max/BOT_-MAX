#!/bin/bash

echo "====================================="
echo "🚀 ЗАПУСК TELEGRAM-БОТА"
echo "====================================="

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Функция для проверки ошибок
check_error() {
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Ошибка!${NC}"
        exit 1
    fi
}

# Шаг 1: Проверяем наличие PostgreSQL
echo -e "${YELLOW}📦 Шаг 1: Проверка PostgreSQL...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker не установлен!${NC}"
    exit 1
fi

# Проверяем, запущен ли контейнер с PostgreSQL
if ! docker ps | grep -q "bot_max-postgres-1"; then
    echo -e "${YELLOW}⚠️  PostgreSQL не запущен. Запускаем через docker-compose...${NC}"
    docker-compose up -d postgres
    check_error
    echo -e "${GREEN}✅ PostgreSQL запущен${NC}"
    # Ждем 5 секунд, чтобы PostgreSQL полностью инициализировался
    echo -e "${YELLOW}⏳ Ожидание инициализации PostgreSQL...${NC}"
    sleep 5
else
    echo -e "${GREEN}✅ PostgreSQL уже запущен${NC}"
fi

# Шаг 2: Проверяем наличие файла .env
echo -e "${YELLOW}🔑 Шаг 2: Проверка файла .env...${NC}"
if [ ! -f .env ]; then
    echo -e "${RED}❌ Файл .env не найден!${NC}"
    echo -e "${YELLOW}📝 Создайте файл .env из .env.example:${NC}"
    echo "cp .env.example .env"
    exit 1
else
    echo -e "${GREEN}✅ Файл .env найден${NC}"
    
    # Проверяем, есть ли TELEGRAM_TOKEN в .env
    if ! grep -q "TELEGRAM_TOKEN" .env; then
        echo -e "${RED}❌ TELEGRAM_TOKEN не найден в .env!${NC}"
        echo -e "${YELLOW}📝 Добавьте строку: TELEGRAM_TOKEN=ваш_токен_от_BotFather${NC}"
        exit 1
    fi
fi

# Шаг 3: Проверяем зависимости Go
echo -e "${YELLOW}📚 Шаг 3: Проверка зависимостей Go...${NC}"
if [ ! -f go.mod ]; then
    echo -e "${RED}❌ go.mod не найден!${NC}"
    exit 1
fi

echo -e "${YELLOW}📦 Установка зависимостей (если нужно)...${NC}"
go mod download
check_error
echo -e "${GREEN}✅ Зависимости Go в порядке${NC}"

# Шаг 4: Проверяем миграции базы данных
echo -e "${YELLOW}🗄️ Шаг 4: Проверка миграций БД...${NC}"
if [ -f migrations/002_telegram_users.sql ]; then
    echo -e "${GREEN}✅ Файл миграции найден${NC}"
    
    # Проверяем, есть ли таблица telegram_users
    if ! docker exec bot_max-postgres-1 psql -U postgres -d medical_bot -c "\dt" | grep -q "telegram_users"; then
        echo -e "${YELLOW}📊 Таблица telegram_users не найдена. Применяем миграцию...${NC}"
        docker exec -i bot_max-postgres-1 psql -U postgres -d medical_bot < migrations/002_telegram_users.sql
        check_error
        echo -e "${GREEN}✅ Миграция применена${NC}"
    else
        echo -e "${GREEN}✅ Таблица telegram_users уже существует${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Файл миграции не найден. Используем AutoMigrate...${NC}"
fi

# Шаг 5: Запуск бота
echo -e "${GREEN}=====================================${NC}"
echo -e "${GREEN}🚀 Шаг 5: Запуск бота...${NC}"
echo -e "${GREEN}=====================================${NC}"

# Запускаем бота и перенаправляем вывод
go run cmd/telegram-bot/main.go 2>&1 | tee bot.log

# Эта команда будет выполняться, пока бот работает
# Для остановки нажми Ctrl+C
