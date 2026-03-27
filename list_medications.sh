#!/bin/bash
echo "📋 СПИСОК ЛЕКАРСТВ:"
echo "===================="
TOKEN=$(curl -s -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d '{"email":"patient@example.com","password":"SecurePass123!"}' | jq -r '.data.token')
if [ ${#TOKEN} -gt 50 ]; then
    curl -s -X GET http://localhost:8080/api/medications \
      -H "Authorization: Bearer $TOKEN" | jq '.data[] | "   • \(.name) - \(.dosage) (\(.frequency))"'
else
    echo "❌ Ошибка получения токена"
fi
