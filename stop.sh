#!/bin/bash
echo "🛑 Остановка медицинского бота"

# Останавливаем процессы на портах
fuser -k 8080/tcp 2>/dev/null
fuser -k 8081/tcp 2>/dev/null

# Убиваем процессы go
pkill -f "go run cmd/api/main.go" 2>/dev/null
pkill -f "go run cmd/telegram/main.go" 2>/dev/null

echo "✅ Бот остановлен"
