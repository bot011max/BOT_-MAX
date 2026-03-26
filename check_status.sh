#!/bin/bash

echo "📊 ПРОВЕРКА СТАТУСА СЕРВИСОВ"
echo "============================"
echo ""

# 1. Main API
echo "1. Main API (порт 8080):"
API_RESP=$(curl -s http://localhost:8080/health 2>/dev/null)
if echo "$API_RESP" | grep -q "ok"; then
    echo "   ✅ API работает: $API_RESP"
else
    echo "   ❌ API не отвечает"
fi
echo ""

# 2. Security API
echo "2. Security API (порт 8090):"
SEC_RESP=$(curl -s http://localhost:8090/security/hsm 2>/dev/null)
if echo "$SEC_RESP" | grep -q "success"; then
    MODE=$(echo "$SEC_RESP" | jq -r '.data.mode' 2>/dev/null)
    echo "   ✅ Security API работает (режим: $MODE)"
else
    echo "   ❌ Security API не отвечает"
fi
echo ""

# 3. Telegram Bot
echo "3. Telegram Bot (порт 8081):"
TG_RESP=$(curl -s http://localhost:8081/health 2>/dev/null)
if echo "$TG_RESP" | grep -q "ok"; then
    echo "   ✅ Telegram бот работает"
else
    echo "   ❌ Telegram бот не отвечает"
fi
echo ""

# 4. JWT токен
echo "4. Получение JWT токена:"
LOGIN_RESP=$(curl -s -X POST http://localhost:8080/api/login \
    -H "Content-Type: application/json" \
    -d '{"email":"patient@example.com","password":"SecurePass123!"}' 2>/dev/null)

if echo "$LOGIN_RESP" | grep -q "token"; then
    echo "   ✅ JWT токен получен"
else
    echo "   ❌ Ошибка получения токена"
fi
echo ""

echo "✅ Проверка завершена"
