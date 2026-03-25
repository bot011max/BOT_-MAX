#!/bin/bash

echo "═══════════════════════════════════════════════════════════════"
echo "🏥 МЕДИЦИНСКИЙ БОТ - MILITARY GRADE SECURITY"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# 1. HSM Status
echo "🔒 1. АППАРАТНОЕ ШИФРОВАНИЕ (HSM):"
curl -s http://localhost:8090/security/hsm | jq '.data | {mode: .mode, encryption: .encryption, device: .device_path}'
echo ""

# 2. Login
echo "🔑 2. АУТЕНТИФИКАЦИЯ:"
LOGIN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"patient@example.com","password":"SecurePass123!"}')
TOKEN=$(echo $LOGIN | jq -r '.data.token')
echo $LOGIN | jq '.data | {token_preview: (.token[:50] + "..."), user: .user}'
echo ""

# 3. Create Medication
echo "💊 3. ДОБАВЛЕНИЕ ЛЕКАРСТВА:"
curl -s -X POST http://localhost:8080/api/medications \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Витамин D3",
    "dosage": "2000 IU",
    "frequency": "1 раз в день",
    "instructions": "Принимать утром с едой"
  }' | jq '.data | {name: .name, dosage: .dosage, frequency: .frequency}'
echo ""

# 4. List Medications
echo "📋 4. СПИСОК ЛЕКАРСТВ:"
curl -s -X GET http://localhost:8080/api/medications \
  -H "Authorization: Bearer $TOKEN" | jq '.data[] | "   • \(.name) - \(.dosage) (\(.frequency))"'
echo ""

# 5. Security Status
echo "🛡️ 5. СТАТУС БЕЗОПАСНОСТИ:"
curl -s http://localhost:8080/security/status | jq '.'
echo ""

# 6. Backups
echo "💾 6. АВТОМАТИЧЕСКИЕ БЭКАПЫ:"
BACKUP_COUNT=$(curl -s http://localhost:8090/security/backups | jq '.data | length')
echo "   Создано бэкапов: $BACKUP_COUNT"
curl -s http://localhost:8090/security/backups | jq '.data[-1] | {id: .id, timestamp: .timestamp, size: .size}'
echo ""

echo "═══════════════════════════════════════════════════════════════"
echo "✅ ДЕМОНСТРАЦИЯ ЗАВЕРШЕНА"
echo "═══════════════════════════════════════════════════════════════"
