#!/bin/bash

echo "🚀 ЗАПУСК МЕДИЦИНСКОГО БОТА (FINAL)"
echo "====================================="

# Останавливаем все
pkill -f "go run" 2>/dev/null
sleep 1

# Запускаем Security API (порт 8090)
echo "🔒 Запуск Security API на порту 8090..."
cd /workspaces/BOT_MAX
go run cmd/security/main.go &
SECURITY_PID=$!
echo $SECURITY_PID > .security_pid

sleep 2

# Запускаем Telegram бота (порт 8081)
echo "🤖 Запуск Telegram бота на порту 8081..."
go run cmd/telegram/main.go &
TELEGRAM_PID=$!
echo $TELEGRAM_PID > .telegram_pid

sleep 2

# Запускаем основной API (порт 8080)
echo "📡 Запуск API на порту 8080..."
go run cmd/api/main.go &
API_PID=$!
echo $API_PID > .api_pid

echo ""
echo "✅ ВСЕ СЕРВИСЫ ЗАПУЩЕНЫ!"
echo "   API: http://localhost:8080"
echo "   Security API: http://localhost:8090"
echo "   Telegram бот: http://localhost:8081"
echo ""
echo "🛡️ АКТИВНЫЕ ЗАЩИТЫ:"
echo "   - HSM аппаратное шифрование (hardware mode)"
echo "   - Автоматическое восстановление (бэкапы)"
echo "   - Распознавание рецептов OCR"
echo "   - JWT аутентификация"
echo "   - Rate Limiting"
echo ""

wait
