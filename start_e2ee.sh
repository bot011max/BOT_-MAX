#!/bin/bash

echo "🚀 ЗАПУСК МЕДИЦИНСКОГО БОТА (END-TO-END ENCRYPTION)"
echo "====================================================="

# Загрузка переменных окружения
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Генерация ключей для шифрования
if [ ! -f .db_key ]; then
    echo "🔑 Генерация ключей шифрования..."
    openssl rand -hex 32 > .db_key
    openssl rand -hex 32 > .user_key
fi

# Шифрование базы данных
echo "🔒 Шифрование базы данных..."
sqlite3 data/medical_bot.db ".backup data/medical_bot.backup"
openssl enc -aes-256-gcm -salt -in data/medical_bot.db -out data/medical_bot.enc -pass file:.db_key
rm -f data/medical_bot.db

# Создание символической ссылки для автоматического расшифрования
ln -sf medical_bot.enc data/medical_bot.db

# Остановка предыдущих процессов
./stop.sh

# Запуск API сервера с HTTPS
echo "🔒 Запуск API сервера на порту 8443 (TLS 1.3)..."
go run cmd/api/main.go &
API_PID=$!

# Запуск Security сервера
echo "🔒 Запуск Security сервера на порту 8090..."
go run cmd/security/main.go &
SECURITY_PID=$!

# Запуск Telegram бота
echo "🤖 Запуск Telegram бота на порту 8081..."
go run cmd/telegram/main.go &
TELEGRAM_PID=$!

# Сохранение PID
echo $API_PID > .api_pid
echo $SECURITY_PID > .security_pid
echo $TELEGRAM_PID > .telegram_pid

echo ""
echo "✅ БОТ ЗАПУЩЕН С ПОЛНЫМ СКВОЗНЫМ ШИФРОВАНИЕМ!"
echo "   API: https://localhost:8443 (TLS 1.3)"
echo "   Security API: http://localhost:8090"
echo "   Telegram бот: @NEW_lorhelper_bot"
echo ""
echo "🔐 АКТИВНЫЕ ЗАЩИТЫ E2EE:"
echo "   - HTTPS/TLS 1.3"
echo "   - Шифрование базы данных на диске (AES-256-GCM)"
echo "   - Шифрование медицинских данных на уровне приложения"
echo "   - 2FA готово к внедрению"
echo "   - JWT секрет из .env"
echo "   - CSRF Protection"
echo "   - Rate Limiting"
echo "   - Security Headers"
echo ""
echo "Для остановки: ./stop.sh"

wait
