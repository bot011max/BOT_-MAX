#!/bin/bash
echo "🧪 Тестирование API"
echo "=================="

echo ""
echo "1. Проверка здоровья:"
curl -s http://localhost:8080/health | python3 -m json.tool 2>/dev/null || curl -s http://localhost:8080/health

echo ""
echo "2. Регистрация пользователя:"
curl -s -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test_'$(date +%s)'@example.com","password":"password123","first_name":"Тест","last_name":"Пользователь","phone":"+79991234567"}' \
  | python3 -m json.tool 2>/dev/null || echo "API не отвечает"

echo ""
echo "✅ Тест завершен"
