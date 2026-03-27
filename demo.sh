#!/bin/bash
echo "═══════════════════════════════════════════════════════════════"
echo "🏥 МЕДИЦИНСКИЙ БОТ - ДЕМОНСТРАЦИЯ"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "1. Health Check:"
curl -s http://localhost:8080/health | jq '.'
echo ""
echo "2. Получение JWT токена:"
LOGIN=$(curl -s -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d '{"email":"patient@example.com","password":"SecurePass123!"}')
echo $LOGIN | jq '.data | {token_preview: (.token[:50] + "..."), user: .user}'
TOKEN=$(echo $LOGIN | jq -r '.data.token')
echo ""
echo "3. Добавление лекарства:"
curl -s -X POST http://localhost:8080/api/medications -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{"name":"Аспирин","dosage":"100 мг","frequency":"1 раз в день"}' | jq '.data | {name: .name, dosage: .dosage}'
echo ""
echo "4. Список лекарств:"
curl -s -X GET http://localhost:8080/api/medications -H "Authorization: Bearer $TOKEN" | jq '.data[] | "   • \(.name) - \(.dosage) (\(.frequency))"'
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "✅ ДЕМОНСТРАЦИЯ ЗАВЕРШЕНА"
echo "═══════════════════════════════════════════════════════════════"
