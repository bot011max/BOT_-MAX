#!/bin/bash
echo "📊 ПРОВЕРКА РАБОТОСПОСОБНОСТИ БОТА"
echo "=================================="
echo ""

echo "1. Main API (порт 8080):"
curl -s http://localhost:8080/health | jq -r '.status' 2>/dev/null && echo "   ✅ API работает" || echo "   ❌ API не работает"

echo ""
echo "2. Security API (порт 8090):"
curl -s http://localhost:8090/security/hsm | jq -r '.data.mode' 2>/dev/null && echo "   ✅ Security API работает" || echo "   ❌ Security API не работает"

echo ""
echo "3. Telegram Bot (порт 8081):"
curl -s http://localhost:8081/health | jq -r '.status' 2>/dev/null && echo "   ✅ Telegram работает" || echo "   ❌ Telegram не работает"

echo ""
echo "4. JWT токен:"
TOKEN=$(curl -s -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d '{"email":"patient@example.com","password":"SecurePass123!"}' | jq -r '.data.token' 2>/dev/null)
if [ ${#TOKEN} -gt 50 ]; then
    echo "   ✅ JWT токен получен"
else
    echo "   ❌ Ошибка получения токена"
fi

echo ""
echo "✅ ПРОВЕРКА ЗАВЕРШЕНА"
