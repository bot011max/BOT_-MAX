#!/bin/bash

echo "🚀 ЗАПУСК ВСЕХ СЕРВИСОВ"
echo "========================"

cd /workspaces/BOT_MAX

# Переменные окружения
export JWT_SECRET="medical_bot_super_secret_key_2026_military_grade_32bytes"
export MASTER_KEY="medical_bot_master_key_for_encryption_2026_32bytes"
export CSRF_SECRET="csrf_protection_secret_key_2026_32bytes"

# Остановка предыдущих
pkill -f "go run" 2>/dev/null
rm -f .api_pid .security_pid .telegram_pid

# Запуск Main API
echo "📡 Запуск Main API (порт 8080)..."
go run cmd/api/main.go &
echo $! > .api_pid
sleep 2

# Запуск Security API
echo "🔒 Запуск Security API (порт 8090)..."
go run cmd/security/main.go &
echo $! > .security_pid
sleep 2

# Запуск Telegram бота
echo "🤖 Запуск Telegram бота (порт 8081)..."
go run cmd/telegram/main.go &
echo $! > .telegram_pid
sleep 3

echo ""
echo "✅ СЕРВИСЫ ЗАПУЩЕНЫ!"
echo "   Main API:     http://localhost:8080"
echo "   Security API: http://localhost:8090"
echo "   Telegram Bot: http://localhost:8081"
echo ""
echo "🛡️ АКТИВНЫЕ ЗАЩИТЫ:"
echo "   - HSM аппаратное шифрование (hardware mode)"
echo "   - Rate Limiting"
echo "   - JWT аутентификация"
echo "   - Security Headers"
echo "   - SQL Injection Protection"
echo "   - XSS Protection"
echo ""
echo "Для проверки: ./check_status.sh"
echo "Для остановки: ./stop_bot.sh"
