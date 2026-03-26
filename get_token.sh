#!/bin/bash
echo "🔐 Получение JWT токена..."
curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"patient@example.com","password":"SecurePass123!"}' | jq '.'
