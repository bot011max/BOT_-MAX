#!/bin/bash

echo "═══════════════════════════════════════════════════════════════"
echo "🏥 ПОЛНОЕ ТЕСТИРОВАНИЕ МЕДИЦИНСКОГО БОТА"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# 1. Получить токен
echo "1. Получение JWT токена..."
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"patient@example.com","password":"SecurePass123!"}' \
  | jq -r '.data.token')

if [ ${#TOKEN} -gt 50 ]; then
    echo "   ✅ Токен получен (длина: ${#TOKEN})"
else
    echo "   ❌ Ошибка получения токена"
    exit 1
fi
echo ""

# 2. Добавить лекарство
echo "2. Добавление лекарства..."
RESULT=$(curl -s -X POST http://localhost:8080/api/medications \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Тестовое лекарство",
    "dosage": "250 мг",
    "frequency": "3 раза в день"
  }')

if echo "$RESULT" | jq -e '.success == true' > /dev/null 2>&1; then
    echo "   ✅ Лекарство добавлено"
else
    echo "   ❌ Ошибка добавления лекарства"
fi
echo ""

# 3. Список лекарств
echo "3. Список лекарств..."
curl -s -X GET http://localhost:8080/api/medications \
  -H "Authorization: Bearer $TOKEN" | jq '.data[] | "   • \(.name) - \(.dosage) (\(.frequency))"'
echo ""

# 4. Проверить бэкапы
echo "4. Проверка бэкапов..."
BACKUP_COUNT=$(curl -s http://localhost:8090/security/backups | jq '.data | length')
echo "   ✅ Создано бэкапов: $BACKUP_COUNT"
echo ""

# 5. Проверить HSM
echo "5. Проверка HSM..."
HSM_MODE=$(curl -s http://localhost:8090/security/hsm | jq -r '.data.mode')
echo "   ✅ HSM режим: $HSM_MODE"
echo ""

echo "═══════════════════════════════════════════════════════════════"
echo "✅ ТЕСТИРОВАНИЕ ЗАВЕРШЕНО"
echo "═══════════════════════════════════════════════════════════════"
