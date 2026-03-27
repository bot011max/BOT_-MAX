#!/bin/bash
cd /workspaces/BOT_MAX/BOT_MAX
export JWT_SECRET="medical_bot_super_secret_key_2026_military_grade_32bytes"
export MASTER_KEY="medical_bot_master_key_for_encryption_2026_32bytes"
pkill -f "go run" 2>/dev/null || true
go run cmd/api/main.go > logs/api.log 2>&1 &
echo $! > .api_pid
go run cmd/security/main.go > logs/security.log 2>&1 &
echo $! > .security_pid
go run cmd/telegram/main.go > logs/telegram.log 2>&1 &
echo $! > .telegram_pid
sleep 3
echo "✅ Бот запущен!"
