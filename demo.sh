#!/bin/bash
echo "🏥 МЕДИЦИНСКИЙ БОТ - ДЕМОНСТРАЦИЯ"
echo "===================================="
echo ""

# Функция для красивого вывода
print_success() { echo -e "\033[32m✅ $1\033[0m"; }
print_error() { echo -e "\033[31m❌ $1\033[0m"; }
print_info() { echo -e "\033[36m📌 $1\033[0m"; }

# 1. Проверка сервисов
print_info "1. Проверка работы сервисов..."
curl -s http://localhost:8080/health > /dev/null && print_success "API сервер работает" || print_error "API сервер не отвечает"
curl -s http://localhost:8081/health > /dev/null && print_success "Telegram бот работает" || print_error "Telegram бот не отвечает"
echo ""

# 2. Регистрация нового пользователя
print_info "2. Регистрация пользователя..."
REG_RESP=$(curl -s -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo_'$(date +%s)'@example.com",
    "password": "SecurePass123!",
    "first_name": "Демо",
    "last_name": "Пользователь",
    "phone": "+79990001122"
  }')
  
if echo $REG_RESP | jq -e '.success' > /dev/null; then
    print_success "Пользователь зарегистрирован"
else
    print_error "Ошибка регистрации"
fi
echo ""

# 3. Аутентификация
print_info "3. Аутентификация..."
LOGIN_RESP=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo_'$(date +%s)'@example.com","password":"SecurePass123!"}')
  
TOKEN=$(echo $LOGIN_RESP | jq -r '.data.token')
if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
    print_success "JWT токен получен"
else
    print_error "Ошибка аутентификации"
    TOKEN=""
fi
echo ""

# 4. Создание лекарства
if [ -n "$TOKEN" ]; then
    print_info "4. Добавление лекарства..."
    curl -s -X POST http://localhost:8080/api/medications \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d '{
        "name": "Витамин D3",
        "dosage": "2000 IU",
        "frequency": "1 раз в день",
        "instructions": "Принимать утром с едой",
        "start_date": "2026-03-25"
      }' | jq '.data.name' && print_success "Лекарство добавлено"
    echo ""
fi

# 5. Получение списка лекарств
if [ -n "$TOKEN" ]; then
    print_info "5. Список лекарств:"
    curl -s -X GET http://localhost:8080/api/medications \
      -H "Authorization: Bearer $TOKEN" | jq '.data[] | {name: .name, dosage: .dosage, frequency: .frequency}'
    echo ""
fi

# 6. Безопасность
print_info "6. Статус безопасности:"
curl -s http://localhost:8080/security/status | jq '.'

echo ""
print_success "Демонстрация завершена!"
print_info "API доступен по адресу: http://localhost:8080"
print_info "Telegram бот: @NEW_lorhelper_bot"
