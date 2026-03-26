#!/bin/bash

echo "🚀 ЗАПУСК МЕДИЦИНСКОГО БОТА"
echo "============================"

# Загрузка переменных окружения
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Проверка JWT_SECRET
if [ -z "$JWT_SECRET" ]; then
    echo "⚠️ JWT_SECRET не установлен, использую значение по умолчанию"
    export JWT_SECRET="medical_bot_super_secret_key_2026_military_grade_32bytes"
fi

# Остановка предыдущих процессов
pkill -f "go run" 2>/dev/null
rm -f .api_pid .security_pid .telegram_pid 2>/dev/null

# Запуск сервисов
echo "📡 Запуск API на порту 8080..."
cd /workspaces/BOT_MAX
go run cmd/api/main.go &
echo $! > .api_pid
sleep 2

echo "🔒 Запуск Security API на порту 8090..."
go run cmd/security/main.go &
echo $! > .security_pid
sleep 2

echo "🤖 Запуск Telegram бота на порту 8081..."
go run cmd/telegram/main.go &
echo $! > .telegram_pid

sleep 3
echo ""
echo "✅ БОТ ЗАПУЩЕН!"
echo "   API: http://localhost:8080"
echo "   Security API: http://localhost:8090"
echo "   Telegram бот: @NEW_lorhelper_bot"
echo ""
echo "🛡️ АКТИВНЫЕ ЗАЩИТЫ:"
echo "   - HSM аппаратное шифрование (hardware mode)"
echo "   - Rate Limiting"
echo "   - JWT аутентификация"
echo "   - Security Headers"
echo "   - SQL Injection Protection"
echo "   - XSS Protection"
echo ""
echo "Для остановки: ./stop.sh"
