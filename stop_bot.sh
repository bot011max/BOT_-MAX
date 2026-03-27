#!/bin/bash
for f in .api_pid .security_pid .telegram_pid; do
    [ -f "$f" ] && kill $(cat "$f") 2>/dev/null
    rm -f "$f"
done
pkill -f "go run" 2>/dev/null
echo "✅ Бот остановлен"
