#!/bin/bash

echo "🚀 ЗАПУСК МЕДИЦИНСКОГО БОТА (SECURE MODE)"
echo "============================================"

# Останавливаем предыдущие процессы
./stop.sh

# Запускаем API сервер с HTTPS
echo "📡 Запуск API сервера на порту 8443 (HTTPS)..."
cd certs && go run ../cmd/api/main.go &
API_PID=$!
cd ..

# Запускаем Security сервер
echo "🔒 Запуск Security сервера на порту 8090..."
go run cmd/security/main.go &
SECURITY_PID=$!

# Запускаем Telegram бота
echo "🤖 Запуск Telegram бота на порту 8081..."
go run cmd/telegram/main.go &
TELEGRAM_PID=$!

# Запускаем планировщик бэкапов
(
    while true; do
        sleep 3600  # Каждый час
        curl -s -X POST http://localhost:8090/security/backup \
            -H "Content-Type: application/json" \
            -d '{"description": "Auto backup - hourly"}'
    done
) &
BACKUP_PID=$!

# Сохраняем PID
echo $API_PID > .api_pid
echo $SECURITY_PID > .security_pid
echo $TELEGRAM_PID > .telegram_pid
echo $BACKUP_PID > .backup_pid

echo ""
echo "✅ БОТ ЗАПУЩЕН В ЗАЩИЩЕННОМ РЕЖИМЕ!"
echo "   API: https://localhost:8443"
echo "   Security API: http://localhost:8090"
echo "   Telegram бот: @NEW_lorhelper_bot"
echo ""
echo "🔐 АКТИВНЫЕ ЗАЩИТЫ:"
echo "   - HTTPS/TLS"
echo "   - Security Headers"
echo "   - Rate Limiting"
echo "   - SQL Injection Protection"
echo "   - CORS Policy"
echo "   - Security Audit"
echo "   - Auto Backups"
echo ""
echo "Для остановки: ./stop.sh"

wait
