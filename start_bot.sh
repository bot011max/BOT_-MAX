#!/bin/bash
echo "🚀 ЗАПУСК МЕДИЦИНСКОГО БОТА"
echo "============================"

cd /workspaces/BOT_MAX

# Переменные окружения
export JWT_SECRET="medical_bot_super_secret_key_2026_military_grade_32bytes"
export MASTER_KEY="medical_bot_master_key_for_encryption_2026_32bytes"

# Создание директорий
mkdir -p data logs

# Проверка базы данных
if [ ! -f "data/medical_bot.db" ]; then
    echo "📦 Создание базы данных..."
    sqlite3 data/medical_bot.db << 'SQL'
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    role TEXT DEFAULT 'patient',
    phone TEXT,
    is_active BOOLEAN DEFAULT 1,
    telegram_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);
INSERT OR IGNORE INTO users (id, email, password_hash, first_name, last_name, role, is_active)
VALUES (
    'test-user-id',
    'patient@example.com',
    '$2a$10$3xxr7wmozPGp1MdZXT1eAeqeI6QJY29NuWU2GBMBgXdSfEchs00cK',
    'Иван',
    'Петров',
    'patient',
    1
);
SQL
fi

# Проверка, что сервисы уже запущены
if pgrep -f "cmd/api/main.go" > /dev/null; then
    echo "✅ Main API уже запущен"
else
    echo "📡 Запуск Main API (порт 8080)..."
    go run cmd/api/main.go > logs/api.log 2>&1 &
    echo $! > .api_pid
    sleep 2
fi

if pgrep -f "cmd/security/main.go" > /dev/null; then
    echo "✅ Security API уже запущен"
else
    echo "🔒 Запуск Security API (порт 8090)..."
    go run cmd/security/main.go > logs/security.log 2>&1 &
    echo $! > .security_pid
    sleep 2
fi

if pgrep -f "cmd/telegram/main.go" > /dev/null; then
    echo "✅ Telegram бот уже запущен"
else
    echo "🤖 Запуск Telegram бота (порт 8081)..."
    go run cmd/telegram/main.go > logs/telegram.log 2>&1 &
    echo $! > .telegram_pid
    sleep 2
fi

echo ""
echo "✅ БОТ ЗАПУЩЕН!"
echo "   Main API:     http://localhost:8080"
echo "   Security API: http://localhost:8090"
echo "   Telegram Bot: http://localhost:8081"
echo ""
echo "🔑 Тестовые данные: patient@example.com / SecurePass123!"
echo ""
echo "📝 Логи: tail -f logs/api.log"
echo "🛑 Остановка: pkill -f 'go run'"
