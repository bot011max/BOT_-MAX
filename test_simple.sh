#!/bin/bash
echo "⚡ ПРОСТАЯ ПРОВЕРКА МЕДИЦИНСКОГО БОТА"
echo "====================================="
echo ""

cd /workspaces/BOT_MAX

# 1. Проверка API
echo "1. Проверка API (порт 8080):"
curl -s http://localhost:8080/health
echo ""
echo ""

# 2. Проверка Security API
echo "2. Проверка Security API (порт 8090):"
curl -s http://localhost:8090/security/hsm
echo ""
echo ""

# 3. Проверка Telegram бота
echo "3. Проверка Telegram бота (порт 8081):"
curl -s http://localhost:8081/health
echo ""
echo ""

# 4. Получение токена
echo "4. Получение JWT токена:"
LOGIN_RESP=$(curl -s -X POST http://localhost:8080/api/login \
    -H "Content-Type: application/json" \
    -d '{"email":"patient@example.com","password":"SecurePass123!"}')
TOKEN=$(echo "$LOGIN_RESP" | grep -o token:[^]*' | head -1 | cut -d' -f4)

if [ -n "$TOKEN" ]; then
    echo "   ✅ Токен получен (первые 50 символов): ${TOKEN:0:50}..."
else
    echo "   ❌ Ошибка получения токена"
    echo "   Ответ: $LOGIN_RESP"
fi
echo ""

# 5. Добавление лекарства
echo "5. Добавление лекарства:"
if [ -n "$TOKEN" ]; then
    curl -s -X POST http://localhost:8080/api/medications \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{"name":"Тест","dosage":"100 мг","frequency":"1 раз"}'
else
    echo "   ❌ Нет токена"
fi
echo ""
echo ""

echo "✅ Проверка завершена"
