#!/bin/bash

echo "🚀 ЗАПУСК МЕДИЦИНСКОГО БОТА"
echo "============================"

cd /workspaces/BOT_MAX/BOT_MAX

# Переменные окружения
export JWT_SECRET="medical_bot_super_secret_key_2026_military_grade_32bytes"
export MASTER_KEY="medical_bot_master_key_for_encryption_2026_32bytes"

# Создание директорий
mkdir -p data logs

# Остановка предыдущих процессов
pkill -f "go run" 2>/dev/null
sleep 1

# Запуск Main API
echo "📡 Запуск Main API (порт 8080)..."
nohup go run cmd/api/main.go > logs/api.log 2>&1 &
echo $! > .api_pid

# Запуск Security API
echo "🔒 Запуск Security API (порт 8090)..."
nohup go run cmd/security/main.go > logs/security.log 2>&1 &
echo $! > .security_pid

# Запуск Telegram бота
echo "🤖 Запуск Telegram бота (порт 8081)..."
nohup go run cmd/telegram/main.go > logs/telegram.log 2>&1 &
echo $! > .telegram_pid

sleep 5

echo ""
echo "✅ БОТ ЗАПУЩЕН!"
echo "   Main API:     http://localhost:8080"
echo "   Security API: http://localhost:8090"
echo "   Telegram Bot: http://localhost:8081"
echo ""
echo "🔑 Тестовые данные: patient@example.com / SecurePass123!"
echo ""
echo "📝 Логи: tail -f logs/api.log"
echo "🛑 Остановка: ./stop_bot.sh"
