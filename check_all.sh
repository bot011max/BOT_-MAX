#!/bin/bash

echo "═══════════════════════════════════════════════════════════════"
echo "🏥 ПРОВЕРКА МЕДИЦИНСКОГО БОТА - MILITARY GRADE SECURITY"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# 1. Проверка Main API
echo "1. 📡 MAIN API (порт 8080):"
API_RESP=$(curl -s http://localhost:8080/health 2>/dev/null)
if echo "$API_RESP" | grep -q "ok"; then
    echo "   ✅ API работает: $API_RESP"
else
    echo "   ❌ API не отвечает"
fi
echo ""

# 2. Проверка Security API
echo "2. 🔒 SECURITY API (порт 8090):"
SEC_RESP=$(curl -s http://localhost:8090/security/hsm 2>/dev/null)
if echo "$SEC_RESP" | grep -q "success"; then
    MODE=$(echo "$SEC_RESP" | grep -o '"mode":"[^"]*"' | cut -d'"' -f4)
    echo "   ✅ Security API работает (режим: $MODE)"
else
    echo "   ❌ Security API не отвечает"
fi
echo ""

# 3. Проверка Telegram бота
echo "3. 🤖 TELEGRAM БОТ (порт 8081):"
TG_RESP=$(curl -s http://localhost:8081/health 2>/dev/null)
if echo "$TG_RESP" | grep -q "ok"; then
    echo "   ✅ Telegram бот работает: $TG_RESP"
else
    echo "   ❌ Telegram бот не отвечает"
fi
echo ""

# 4. Получение JWT токена
echo "4. 🔑 АУТЕНТИФИКАЦИЯ (JWT):"
LOGIN=$(curl -s -X POST http://localhost:8080/api/login \
    -H "Content-Type: application/json" \
    -d '{"email":"patient@example.com","password":"SecurePass123!"}' 2>/dev/null)

if echo "$LOGIN" | grep -q "token"; then
    TOKEN=$(echo "$LOGIN" | jq -r '.data.token')
    echo "   ✅ JWT токен получен (длина: ${#TOKEN})"
    echo "   👤 Пользователь: $(echo "$LOGIN" | jq -r '.data.user.email')"
else
    echo "   ❌ Ошибка аутентификации"
    TOKEN=""
fi
echo ""

# 5. Проверка лекарств
if [ -n "$TOKEN" ]; then
    echo "5. 💊 УПРАВЛЕНИЕ ЛЕКАРСТВАМИ:"
    
    # Добавление лекарства
    echo "   ➕ Добавление лекарства:"
    ADD=$(curl -s -X POST http://localhost:8080/api/medications \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{"name":"Тестовое лекарство","dosage":"250 мг","frequency":"2 раза в день"}' 2>/dev/null)
    
    if echo "$ADD" | grep -q "success"; then
        echo "      ✅ Лекарство добавлено: $(echo "$ADD" | jq -r '.data.name')"
    else
        echo "      ❌ Ошибка добавления"
    fi
    
    # Список лекарств
    echo "   📋 Список лекарств:"
    LIST=$(curl -s -X GET http://localhost:8080/api/medications \
        -H "Authorization: Bearer $TOKEN" 2>/dev/null)
    
    if echo "$LIST" | grep -q "data"; then
        echo "$LIST" | jq -r '.data[] | "      • \(.name) - \(.dosage) (\(.frequency))"' 2>/dev/null || echo "      ⚠️ Нет лекарств"
    else
        echo "      ❌ Ошибка получения списка"
    fi
fi
echo ""

# 6. Проверка бэкапов
echo "6. 💾 СИСТЕМА БЭКАПОВ:"
BACKUPS=$(curl -s http://localhost:8090/security/backups 2>/dev/null)
COUNT=$(echo "$BACKUPS" | jq '.data | length' 2>/dev/null)
if [ -n "$COUNT" ] && [ "$COUNT" -gt 0 ]; then
    echo "   ✅ Создано бэкапов: $COUNT"
    echo "$BACKUPS" | jq -r '.data[-1] | "   📁 Последний: \(.id) (\(.timestamp | .[0:19]))"' 2>/dev/null
else
    echo "   ⚠️ Бэкапы не найдены"
fi
echo ""

# 7. Проверка базы данных
echo "7. 🗄️ БАЗА ДАННЫХ:"
if [ -f "data/medical_bot.db" ]; then
    SIZE=$(du -h data/medical_bot.db | cut -f1)
    USERS=$(sqlite3 data/medical_bot.db "SELECT COUNT(*) FROM users;" 2>/dev/null)
    MEDS=$(sqlite3 data/medical_bot.db "SELECT COUNT(*) FROM medications;" 2>/dev/null)
    echo "   ✅ База данных: $SIZE"
    echo "   👥 Пользователей: $USERS"
    echo "   💊 Лекарств: $MEDS"
else
    echo "   ❌ База данных не найдена"
fi
echo ""

# 8. Проверка процессов
echo "8. 🔄 ПРОЦЕССЫ:"
if pgrep -f "cmd/api/main.go" > /dev/null; then
    echo "   ✅ Main API: запущен"
else
    echo "   ❌ Main API: не запущен"
fi
if pgrep -f "cmd/security/main.go" > /dev/null; then
    echo "   ✅ Security API: запущен"
else
    echo "   ❌ Security API: не запущен"
fi
if pgrep -f "cmd/telegram/main.go" > /dev/null; then
    echo "   ✅ Telegram бот: запущен"
else
    echo "   ❌ Telegram бот: не запущен"
fi
echo ""

# 9. Итоговая оценка
echo "═══════════════════════════════════════════════════════════════"
echo "📊 ИТОГОВАЯ ОЦЕНКА:"
echo "═══════════════════════════════════════════════════════════════"

SCORE=0
TOTAL=8

# Подсчет баллов
[ -n "$API_RESP" ] && SCORE=$((SCORE + 1))
[ -n "$SEC_RESP" ] && SCORE=$((SCORE + 1))
[ -n "$TG_RESP" ] && SCORE=$((SCORE + 1))
[ -n "$TOKEN" ] && SCORE=$((SCORE + 1))
[ -n "$COUNT" ] && SCORE=$((SCORE + 1))
[ -f "data/medical_bot.db" ] && SCORE=$((SCORE + 1))
pgrep -f "cmd/api/main.go" > /dev/null && SCORE=$((SCORE + 1))
pgrep -f "cmd/security/main.go" > /dev/null && SCORE=$((SCORE + 1))

PERCENT=$((SCORE * 100 / TOTAL))

echo "   ✅ Проверено: $SCORE из $TOTAL"
echo "   📈 Уровень: $PERCENT%"

if [ $PERCENT -ge 90 ]; then
    echo "   🏆 Статус: ВЫСОКИЙ - Все системы работают идеально"
elif [ $PERCENT -ge 70 ]; then
    echo "   🟢 Статус: ХОРОШИЙ - Основные функции работают"
elif [ $PERCENT -ge 50 ]; then
    echo "   🟡 Статус: СРЕДНИЙ - Требуется проверка некоторых компонентов"
else
    echo "   🔴 Статус: КРИТИЧЕСКИЙ - Бот не работает"
fi

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "✅ ПРОВЕРКА ЗАВЕРШЕНА"
echo "═══════════════════════════════════════════════════════════════"
