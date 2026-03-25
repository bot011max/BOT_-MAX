#!/bin/bash

echo "🛑 Остановка медицинского бота"

# Останавливаем все процессы
if [ -f .api_pid ]; then
    kill $(cat .api_pid) 2>/dev/null
    rm .api_pid
fi

if [ -f .security_pid ]; then
    kill $(cat .security_pid) 2>/dev/null
    rm .security_pid
fi

if [ -f .telegram_pid ]; then
    kill $(cat .telegram_pid) 2>/dev/null
    rm .telegram_pid
fi

# Дополнительная очистка
pkill -f "go run cmd/api/main.go" 2>/dev/null
pkill -f "go run cmd/security/main.go" 2>/dev/null
pkill -f "go run cmd/telegram/main.go" 2>/dev/null

echo "✅ Бот остановлен"
