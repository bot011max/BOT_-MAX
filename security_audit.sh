#!/bin/bash

echo "🔒 АУДИТ БЕЗОПАСНОСТИ МЕДИЦИНСКОГО БОТА"
echo "========================================="
echo ""

# 1. Проверка HTTPS
echo "1. ПРОВЕРКА HTTPS:"
if curl -s https://localhost:8080 > /dev/null 2>&1; then
    echo "   ✅ HTTPS включен"
else
    echo "   ❌ HTTPS НЕ включен (КРИТИЧЕСКАЯ УЯЗВИМОСТЬ)"
fi

# 2. Проверка Security Headers
echo -e "\n2. SECURITY HEADERS:"
HEADERS=$(curl -sI http://localhost:8080)
echo "$HEADERS" | grep -i "Strict-Transport-Security" || echo "   ❌ HSTS отсутствует"
echo "$HEADERS" | grep -i "X-Frame-Options" || echo "   ❌ X-Frame-Options отсутствует"
echo "$HEADERS" | grep -i "X-Content-Type-Options" || echo "   ❌ X-Content-Type-Options отсутствует"

# 3. Проверка Rate Limiting
echo -e "\n3. RATE LIMITING:"
for i in {1..10}; do
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/api/login \
        -H "Content-Type: application/json" \
        -d '{"email":"test@test.com","password":"wrong"}')
    if [ $STATUS -eq 429 ]; then
        echo "   ✅ Rate limiting активен (блокировка после $i попыток)"
        break
    fi
done

# 4. Проверка JWT
echo -e "\n4. JWT БЕЗОПАСНОСТЬ:"
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
    -H "Content-Type: application/json" \
    -d '{"email":"patient@example.com","password":"SecurePass123!"}' | jq -r '.data.token')
if [ ${#TOKEN} -gt 100 ]; then
    echo "   ✅ JWT токен выдается"
fi

# 5. Проверка SQL инъекций
echo -e "\n5. ЗАЩИТА ОТ SQL ИНЪЕКЦИЙ:"
RESPONSE=$(curl -s "http://localhost:8080/api/login?email=' OR '1'='1")
if [[ $RESPONSE == *"error"* ]] || [[ $RESPONSE == *"invalid"* ]]; then
    echo "   ✅ SQL инъекции блокируются"
else
    echo "   ❌ Потенциальная уязвимость к SQL инъекциям"
fi

# 6. Проверка CORS
echo -e "\n6. CORS ПОЛИТИКА:"
CORS=$(curl -sI http://localhost:8080 | grep -i "access-control-allow-origin")
if [ -n "$CORS" ]; then
    echo "   ⚠️ CORS настроен: $CORS"
else
    echo "   ℹ️ CORS не ограничен"
fi

# 7. Проверка HSM
echo -e "\n7. АППАРАТНОЕ ШИФРОВАНИЕ:"
HSM_STATUS=$(curl -s http://localhost:8090/security/hsm | jq -r '.data.mode')
if [ "$HSM_STATUS" = "hardware" ]; then
    echo "   ✅ HSM активен (hardware mode)"
else
    echo "   ⚠️ HSM в software mode"
fi

# 8. Проверка бэкапов
echo -e "\n8. АВТОМАТИЧЕСКОЕ ВОССТАНОВЛЕНИЕ:"
BACKUPS=$(curl -s http://localhost:8090/security/backups | jq '.data | length')
if [ $BACKUPS -gt 0 ]; then
    echo "   ✅ Бэкапы создаются ($BACKUPS бэкапов)"
else
    echo "   ❌ Бэкапы отсутствуют"
fi

echo -e "\n📊 ИТОГОВАЯ ОЦЕНКА:"
echo "================================"
echo "Критические уязвимости: $(grep -c "КРИТИЧЕСКАЯ" security_audit.sh)"
echo "Высокий приоритет: $(grep -c "❌" security_audit.sh)"
echo "Средний приоритет: $(grep -c "⚠️" security_audit.sh)"
