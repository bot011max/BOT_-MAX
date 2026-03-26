#!/bin/bash

echo "🔍 ПРОВЕРКА УЛУЧШЕНИЙ БЕЗОПАСНОСТИ"
echo "===================================="
echo ""

# 1. Проверка HTTPS
echo "1. HTTPS/TLS:"
if curl -sk https://localhost:8443/health > /dev/null 2>&1; then
    echo "   ✅ HTTPS активен (порт 8443)"
    TLS_VERSION=$(curl -skv https://localhost:8443/health 2>&1 | grep "TLSv" | head -1)
    echo "   $TLS_VERSION"
else
    echo "   ❌ HTTPS не активен"
fi
echo ""

# 2. Проверка JWT секрета
echo "2. JWT СЕКРЕТ:"
if [ -n "$JWT_SECRET" ]; then
    echo "   ✅ JWT секрет загружен из .env"
    echo "   Длина: ${#JWT_SECRET}"
else
    echo "   ❌ JWT секрет не найден"
fi
echo ""

# 3. Проверка CSRF
echo "3. CSRF ЗАЩИТА:"
CSRF_HEADER=$(curl -sI https://localhost:8443/health 2>/dev/null | grep -i "x-csrf-token")
if [ -n "$CSRF_HEADER" ]; then
    echo "   ✅ CSRF защита активна"
else
    echo "   ⚠️ CSRF защита требует авторизации"
fi
echo ""

# 4. Проверка шифрования БД
echo "4. ШИФРОВАНИЕ БАЗЫ ДАННЫХ:"
if [ -n "$DB_ENCRYPTION_KEY" ]; then
    echo "   ✅ Ключ шифрования БД установлен"
else
    echo "   ⚠️ Ключ шифрования БД не установлен"
fi
echo ""

# 5. Проверка Security Headers
echo "5. SECURITY HEADERS:"
HEADERS=$(curl -skI https://localhost:8443/health 2>/dev/null)
echo "$HEADERS" | grep -i "Strict-Transport-Security" && echo "   ✅ HSTS активен" || echo "   ❌ HSTS отсутствует"
echo "$HEADERS" | grep -i "X-Frame-Options" && echo "   ✅ X-Frame-Options активен" || echo "   ❌ X-Frame-Options отсутствует"
echo "$HEADERS" | grep -i "X-Content-Type-Options" && echo "   ✅ X-Content-Type-Options активен" || echo "   ❌ X-Content-Type-Options отсутствует"
echo ""

echo "✅ Проверка завершена"
