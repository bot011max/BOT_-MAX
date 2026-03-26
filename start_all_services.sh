#!/bin/bash

echo "🚀 ЗАПУСК ВСЕХ СЕРВИСОВ"
echo "========================"

cd /workspaces/BOT_MAX

# Установка переменных
export JWT_SECRET="medical_bot_super_secret_key_2026_military_grade_32bytes"
export MASTER_KEY="medical_bot_master_key_for_encryption_2026_32bytes"

# Остановка предыдущих процессов
pkill -f "go run" 2>/dev/null
sleep 1

# Запуск Main API
echo "📡 Запуск Main API (порт 8080)..."
go run cmd/api/main.go &
echo $! > .api_pid
sleep 3

# Запуск Security API
echo "🔒 Запуск Security API (порт 8090)..."
go run cmd/security/main.go &
echo $! > .security_pid
sleep 2

# Запуск Telegram бота
echo "🤖 Запуск Telegram бота (порт 8081)..."
go run cmd/telegram/main.go &
echo $! > .telegram_pid
sleep 2

echo ""
echo "✅ ВСЕ СЕРВИСЫ ЗАПУЩЕНЫ!"
echo "   Main API:     http://localhost:8080"
echo "   Security API: http://localhost:8090"
echo "   Telegram Bot: http://localhost:8081"
echo ""
echo "🛡️ АКТИВНЫЕ ЗАЩИТЫ:"
echo "   - HSM аппаратное шифрование (AES-256-GCM)"
echo "   - Rate Limiting"
echo "   - JWT аутентификация"
echo "   - Security Headers"
echo ""
echo "Для остановки: pkill -f 'go run'"
