#!/bin/bash

echo "🚀 ЗАПУСК МЕДИЦИНСКОГО БОТА (SIMPLE)"
echo "====================================="

cd /workspaces/BOT_MAX

# Переменные окружения
export JWT_SECRET="medical_bot_super_secret_key_2026_military_grade_32bytes"
export MASTER_KEY="medical_bot_master_key_for_encryption_2026_32bytes"
export CSRF_SECRET="csrf_protection_secret_key_2026_32bytes"

# Остановка предыдущих
pkill -f "go run" 2>/dev/null
sleep 1

# Запуск в фоне с перенаправлением вывода
echo "📡 Запуск Main API (порт 8080)..."
nohup go run cmd/api/main.go > logs/api.log 2>&1 &
echo $! > .api_pid

echo "🔒 Запуск Security API (порт 8090)..."
nohup go run cmd/security/main.go > logs/security.log 2>&1 &
echo $! > .security_pid

echo "🤖 Запуск Telegram бота (порт 8081)..."
nohup go run cmd/telegram/main.go > logs/telegram.log 2>&1 &
echo $! > .telegram_pid

echo ""
echo "Ждем запуска сервисов..."
sleep 5

# Проверка
echo ""
echo "📊 ПРОВЕРКА:"
curl -s http://localhost:8080/health && echo " ✅ API работает" || echo " ❌ API не работает"
curl -s http://localhost:8090/security/hsm | jq -r '.data.mode' 2>/dev/null && echo " ✅ Security API работает" || echo " ❌ Security API не работает"
curl -s http://localhost:8081/health && echo " ✅ Telegram работает" || echo " ❌ Telegram не работает"

echo ""
echo "📝 Логи:"
echo "   API: tail -f logs/api.log"
echo "   Security: tail -f logs/security.log"
echo "   Telegram: tail -f logs/telegram.log"
echo ""
echo "Для остановки: pkill -f 'go run'"
