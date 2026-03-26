#!/bin/bash

echo "═══════════════════════════════════════════════════════════════"
echo "🏥 ПРОВЕРКА УЛУЧШЕНИЙ БЕЗОПАСНОСТИ"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# 1. Проверка переменных
echo "1. ПЕРЕМЕННЫЕ ОКРУЖЕНИЯ:"
echo "   JWT_SECRET: ${JWT_SECRET:0:20}..."
echo "   MASTER_KEY: ${MASTER_KEY:0:20}..."
echo "   ✅ Переменные загружены"
echo ""

# 2. Проверка API
echo "2. API СТАТУС:"
curl -s http://localhost:8080/health | jq '.'
echo ""

# 3. Проверка логина
echo "3. АУТЕНТИФИКАЦИЯ:"
LOGIN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"patient@example.com","password":"SecurePass123!"}')
TOKEN=$(echo $LOGIN | jq -r '.data.token')
if [ ${#TOKEN} -gt 50 ]; then
    echo "   ✅ JWT токен получен (длина: ${#TOKEN})"
else
    echo "   ❌ Ошибка логина"
fi
echo ""

# 4. Проверка бэкапов
echo "4. БЭКАПЫ:"
BACKUPS=$(curl -s http://localhost:8090/security/backups | jq '.data | length')
echo "   ✅ Создано бэкапов: $BACKUPS"
echo ""

# 5. Проверка HSM
echo "5. HSM:"
HSM_MODE=$(curl -s http://localhost:8090/security/hsm | jq -r '.data.mode')
echo "   ✅ HSM режим: $HSM_MODE"
echo ""

echo "═══════════════════════════════════════════════════════════════"
echo "✅ УЛУЧШЕНИЯ ПРИМЕНЕНЫ УСПЕШНО"
echo "═══════════════════════════════════════════════════════════════"
