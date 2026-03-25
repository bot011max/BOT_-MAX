#!/bin/bash

echo "🏥 ТЕСТИРОВАНИЕ НОВЫХ ФУНКЦИЙ МЕДИЦИНСКОГО БОТА"
echo "================================================"

# Проверяем доступность API
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "❌ API сервер не запущен! Запустите ./start.sh"
    exit 1
fi

echo "✅ API сервер работает"

# 1. Проверка статуса безопасности
echo -e "\n🔐 1. Статус безопасности:"
curl -s http://localhost:8080/security/status | jq '.'

# 2. Проверка HSM
echo -e "\n🛡️ 2. Аппаратное шифрование (HSM):"
curl -s http://localhost:8080/security/hsm 2>/dev/null || echo "Эндпоинт HSM пока не добавлен"

# 3. Создание бэкапа
echo -e "\n💾 3. Создание бэкапа:"
curl -s -X POST http://localhost:8080/security/backup \
  -H "Content-Type: application/json" \
  -d '{"description": "Test backup"}' 2>/dev/null || echo "Эндпоинт backup пока не добавлен"

# 4. Регистрация пользователя
echo -e "\n📝 4. Регистрация пользователя:"
REG_RESP=$(curl -s -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo_'$(date +%s)'@example.com",
    "password": "Test123!",
    "first_name": "Демо",
    "last_name": "Пользователь",
    "phone": "+79990001122"
  }')
echo $REG_RESP | jq '.'

# 5. Логин и получение токена
echo -e "\n🔑 5. Аутентификация:"
LOGIN_RESP=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"patient@example.com","password":"SecurePass123!"}')
  
TOKEN=$(echo $LOGIN_RESP | jq -r '.data.token')
if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
    echo "✅ Токен получен: ${TOKEN:0:50}..."
    
    # 6. Добавление лекарства
    echo -e "\n💊 6. Добавление лекарства:"
    curl -s -X POST http://localhost:8080/api/medications \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "name": "Витамин D3",
        "dosage": "2000 IU",
        "frequency": "1 раз в день",
        "instructions": "Принимать утром с едой",
        "start_date": "2026-03-25"
      }' | jq '.'
    
    # 7. Получение списка лекарств
    echo -e "\n📋 7. Список лекарств:"
    curl -s -X GET http://localhost:8080/api/medications \
      -H "Authorization: Bearer $TOKEN" | jq '.data[] | {name: .name, dosage: .dosage, frequency: .frequency}'
else
    echo "❌ Ошибка аутентификации"
fi

echo -e "\n✅ Тестирование завершено!"
