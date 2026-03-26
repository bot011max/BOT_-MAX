#!/bin/bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"patient@example.com","password":"SecurePass123!"}' \
  | jq -r '.data.token')

if [ ${#TOKEN} -gt 50 ]; then
    echo "💊 Добавление лекарства..."
    curl -X POST http://localhost:8080/api/medications \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "name": "'"$1"'",
        "dosage": "'"$2"'",
        "frequency": "'"$3"'"
      }' | jq '.'
else
    echo "❌ Ошибка получения токена"
fi
