#!/bin/bash

echo "🚀 ЗАПУСК МЕДИЦИНСКОГО БОТА"
echo "============================"

# Останавливаем предыдущие процессы
./stop.sh

# Запускаем API сервер в фоне
echo "📡 Запуск API сервера на порту 8080..."
go run cmd/api/main.go &
API_PID=$!

# Ждем запуска API
sleep 2

# Запускаем Security сервер в фоне
echo "🔒 Запуск Security сервера на порту 8090..."
go run cmd/security/main.go &
SECURITY_PID=$!

# Запускаем Telegram бота в фоне
echo "🤖 Запуск Telegram бота на порту 8081..."
go run cmd/telegram/main.go &
TELEGRAM_PID=$!

# Сохраняем PID для остановки
echo $API_PID > .api_pid
echo $SECURITY_PID > .security_pid
echo $TELEGRAM_PID > .telegram_pid

echo ""
echo "✅ БОТ ЗАПУЩЕН!"
echo "   API: http://localhost:8080"
echo "   Security API: http://localhost:8090"
echo "   Telegram бот: @NEW_lorhelper_bot"
echo ""
echo "Для остановки: ./stop.sh"
echo ""

# Ждем завершения
wait
