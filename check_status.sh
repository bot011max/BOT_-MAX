#!/bin/bash
echo "📊 ПРОВЕРКА СТАТУСА СЕРВИСОВ"
echo "============================"
echo ""

echo "1. Main API (порт 8080):"
curl -s http://localhost:8080/health | jq '.' 2>/dev/null || echo "   ❌ Не отвечает"

echo ""
echo "2. Security API (порт 8090):"
curl -s http://localhost:8090/security/hsm | jq '.data.mode' 2>/dev/null || echo "   ❌ Не отвечает"

echo ""
echo "3. Telegram Bot (порт 8081):"
curl -s http://localhost:8081/health | jq '.' 2>/dev/null || echo "   ❌ Не отвечает"

echo ""
echo "✅ Проверка завершена"
