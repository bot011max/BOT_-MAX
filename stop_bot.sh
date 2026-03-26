#!/bin/bash
echo "🛑 Остановка медицинского бота..."

for pid_file in .api_pid .security_pid .telegram_pid; do
    if [ -f "$pid_file" ]; then
        kill $(cat "$pid_file") 2>/dev/null
        rm -f "$pid_file"
    fi
done

pkill -f "go run" 2>/dev/null
echo "✅ Бот остановлен"
