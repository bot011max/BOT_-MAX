#!/bin/bash

echo "🔍 ПРОВЕРКА ВСЕХ СЕРВИСОВ"
echo "=========================="
echo ""

# 1. Проверка Security API
echo "1. SECURITY API (порт 8090):"
HSM_STATUS=$(curl -s http://localhost:8090/security/hsm 2>/dev/null | jq -r '.data.mode')
if [ "$HSM_STATUS" = "hardware" ]; then
    echo "   ✅ HSM: HARDWARE режим"
else
    echo "   ⚠️ HSM: $HSM_STATUS"
fi

BACKUPS=$(curl -s http://localhost:8090/security/backups 2>/dev/null | jq '.data | length')
echo "   💾 Бэкапов: $BACKUPS"
echo ""

# 2. Проверка Telegram бота
echo "2. TELEGRAM БОТ (порт 8081):"
TG_STATUS=$(curl -s http://localhost:8081/health 2>/dev/null | jq -r '.status')
if [ "$TG_STATUS" = "ok" ]; then
    echo "   ✅ Telegram бот: работает"
else
    echo "   ❌ Telegram бот: не отвечает"
fi
echo ""

# 3. Проверка основного API
echo "3. ОСНОВНОЙ API (порт 8080):"
API_STATUS=$(curl -s http://localhost:8080/health 2>/dev/null | jq -r '.status')
if [ "$API_STATUS" = "ok" ]; then
    echo "   ✅ API: работает"
    
    # Тест логина
    echo -e "\n4. ТЕСТ АУТЕНТИФИКАЦИИ:"
    LOGIN=$(curl -s -X POST http://localhost:8080/api/login \
        -H "Content-Type: application/json" \
        -d '{"email":"patient@example.com","password":"SecurePass123!"}' 2>/dev/null)
    
    TOKEN=$(echo $LOGIN | jq -r '.data.token')
    if [ ${#TOKEN} -gt 50 ]; then
        echo "   ✅ JWT токен получен"
        
        # Тест профиля
        PROFILE=$(curl -s -X GET http://localhost:8080/api/profile \
            -H "Authorization: Bearer $TOKEN" 2>/dev/null)
        USER_NAME=$(echo $PROFILE | jq -r '.data.first_name')
        if [ "$USER_NAME" != "null" ]; then
            echo "   ✅ Профиль: $USER_NAME"
        fi
    else
        echo "   ❌ Ошибка аутентификации"
    fi
else
    echo "   ❌ API: не отвечает"
fi

echo ""
echo "✅ Проверка завершена!"
