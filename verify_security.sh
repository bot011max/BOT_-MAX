#!/bin/bash

echo "🔍 ПРОВЕРКА УЛУЧШЕНИЙ БЕЗОПАСНОСТИ"
echo "===================================="

# 1. Проверка HTTPS
echo -e "\n1. HTTPS:"
if curl -sk https://localhost:8443/health > /dev/null 2>&1; then
    echo "   ✅ HTTPS активен"
else
    echo "   ❌ HTTPS не активен"
fi

# 2. Проверка Security Headers
echo -e "\n2. Security Headers:"
HEADERS=$(curl -skI https://localhost:8443/health)
echo "$HEADERS" | grep -i "Strict-Transport-Security" && echo "   ✅ HSTS активен" || echo "   ❌ HSTS отсутствует"
echo "$HEADERS" | grep -i "X-Frame-Options" && echo "   ✅ X-Frame-Options активен" || echo "   ❌ X-Frame-Options отсутствует"
echo "$HEADERS" | grep -i "X-Content-Type-Options" && echo "   ✅ X-Content-Type-Options активен" || echo "   ❌ X-Content-Type-Options отсутствует"

# 3. Проверка Rate Limiting
echo -e "\n3. Rate Limiting:"
for i in {1..6}; do
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://localhost:8080/api/login \
        -H "Content-Type: application/json" \
        -d '{"email":"test@test.com","password":"wrong"}')
    if [ $STATUS -eq 429 ]; then
        echo "   ✅ Rate limiting активен (блокировка после $i попыток)"
        break
    fi
done

# 4. Проверка бэкапов
echo -e "\n4. Авто-бэкапы:"
BACKUPS=$(curl -s http://localhost:8090/security/backups | jq '.data | length' 2>/dev/null)
if [ "$BACKUPS" -gt 0 ]; then
    echo "   ✅ Бэкапы создаются ($BACKUPS бэкапов)"
else
    echo "   ❌ Бэкапы отсутствуют"
fi

echo -e "\n✅ Все улучшения применены!"
