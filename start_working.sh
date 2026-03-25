#!/bin/bash

echo "🚀 ЗАПУСК МЕДИЦИНСКОГО БОТА (WORKING MODE)"
echo "============================================"

# Останавливаем предыдущие процессы
./stop.sh

# Запускаем оригинальный API (рабочую версию)
echo "📡 Запуск API сервера на порту 8080..."
cd /workspaces/BOT_MAX
go run cmd/api/main.go &
API_PID=$!

# Запускаем Security сервер
echo "🔒 Запуск Security сервера на порту 8090..."
go run cmd/security/main.go &
SECURITY_PID=$!

# Запускаем Telegram бота
echo "🤖 Запуск Telegram бота на порту 8081..."
go run cmd/telegram/main.go &
TELEGRAM_PID=$!

# Сохраняем PID
echo $API_PID > .api_pid
echo $SECURITY_PID > .security_pid
echo $TELEGRAM_PID > .telegram_pid

echo ""
echo "✅ БОТ ЗАПУЩЕН!"
echo "   API: http://localhost:8080"
echo "   Security API: http://localhost:8090"
echo "   Telegram бот: http://localhost:8081"
echo ""
echo "Для проверки: ./check_all.sh"
echo "Для остановки: ./stop.sh"

wait
