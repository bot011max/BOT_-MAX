#!/bin/bash

echo "🚀 НАГРУЗОЧНОЕ ТЕСТИРОВАНИЕ МЕДИЦИНСКОГО БОТА"
echo "=============================================="

# Тестовые данные
cat > targets.txt << 'TARGETS'
POST http://localhost:8080/api/register
Content-Type: application/json
@register.json

POST http://localhost:8080/api/login
Content-Type: application/json
@login.json

GET http://localhost:8080/api/medications
Authorization: Bearer {token}

POST http://localhost:8080/api/medications
Content-Type: application/json
Authorization: Bearer {token}
@medication.json
TARGETS

# Тестовые JSON файлы
cat > register.json << 'JSON'
{"email":"test@example.com","password":"Test123!","first_name":"Test","last_name":"User"}
JSON

cat > login.json << 'JSON'
{"email":"patient@example.com","password":"SecurePass123!"}
JSON

cat > medication.json << 'JSON'
{"name":"Test Medication","dosage":"100 mg","frequency":"1 time per day"}
JSON

echo ""
echo "1. Легкая нагрузка (50 запросов/сек, 30 секунд):"
vegeta attack -rate=50 -duration=30s -targets=targets.txt | vegeta report

echo ""
echo "2. Средняя нагрузка (100 запросов/сек, 30 секунд):"
vegeta attack -rate=100 -duration=30s -targets=targets.txt | vegeta report

echo ""
echo "3. Высокая нагрузка (200 запросов/сек, 30 секунд):"
vegeta attack -rate=200 -duration=30s -targets=targets.txt | vegeta report

echo ""
echo "4. Пиковая нагрузка (500 запросов/сек, 30 секунд):"
vegeta attack -rate=500 -duration=30s -targets=targets.txt | vegeta report

echo ""
echo "✅ Нагрузочное тестирование завершено"
