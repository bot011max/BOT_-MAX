#!/bin/bash
echo "🚀 ЗАПУСК МЕДИЦИНСКОГО БОТА"
echo "============================"

# Останавливаем старые процессы
./stop.sh 2>/dev/null

# Ждем освобождения портов
sleep 2

# Запускаем API на порту 8080
echo "📡 Запуск API сервера на порту 8080..."
go run cmd/api/main.go &
API_PID=$!

# Ждем 2 секунды
sleep 2

# Запускаем Telegram бота на порту 8081
echo "🤖 Запуск Telegram бота на порту 8081..."
go run cmd/telegram/main.go &
BOT_PID=$!

# Сохраняем PID
echo $API_PID > .api.pid
echo $BOT_PID > .bot.pid

echo ""
echo "✅ БОТ ЗАПУЩЕН!"
echo "   API: http://localhost:8080"
echo "   Telegram бот: @NEW_lorhelper_bot"
echo ""
echo "Для остановки: ./stop.sh"
echo ""

wait
