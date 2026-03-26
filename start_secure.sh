#!/bin/bash

echo "🚀 ЗАПУСК МЕДИЦИНСКОГО БОТА (SECURE MODE)"
echo "============================================"

# Загрузка переменных окружения
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Остановка предыдущих процессов
./stop.sh

# Запуск API сервера с HTTPS
echo "🔒 Запуск API сервера на порту 8443 (HTTPS)..."
go run cmd/api/main.go &
API_PID=$!

# Запуск Security сервера
echo "🔒 Запуск Security сервера на порту 8090..."
go run cmd/security/main.go &
SECURITY_PID=$!

# Запуск Telegram бота
echo "🤖 Запуск Telegram бота на порту 8081..."
go run cmd/telegram/main.go &
TELEGRAM_PID=$!

# Сохранение PID
echo $API_PID > .api_pid
echo $SECURITY_PID > .security_pid
echo $TELEGRAM_PID > .telegram_pid

echo ""
echo "✅ БОТ ЗАПУЩЕН В ЗАЩИЩЕННОМ РЕЖИМЕ!"
echo "   API: https://localhost:8443"
echo "   Security API: http://localhost:8090"
echo "   Telegram бот: @NEW_lorhelper_bot"
echo ""
echo "🔐 АКТИВНЫЕ ЗАЩИТЫ:"
echo "   - HTTPS/TLS 1.3"
echo "   - JWT секрет из .env"
echo "   - CSRF Protection"
echo "   - Шифрование базы данных"
echo "   - 2FA готово"
echo "   - Security Headers"
echo "   - Rate Limiting"
echo ""
echo "Для остановки: ./stop.sh"

wait
