#!/bin/bash

echo "🔒 ТЕСТИРОВАНИЕ SECURITY API"
echo "=============================="

# Проверяем доступность Security API
if ! curl -s http://localhost:8090/security/hsm > /dev/null 2>&1; then
    echo "❌ Security API не запущен! Запустите ./start.sh"
    exit 1
fi

echo "✅ Security API работает"

# 1. Проверка HSM
echo -e "\n🛡️ 1. Аппаратное шифрование (HSM):"
curl -s http://localhost:8090/security/hsm | jq '.'

# 2. Создание бэкапа
echo -e "\n💾 2. Создание бэкапа:"
curl -s -X POST http://localhost:8090/security/backup \
  -H "Content-Type: application/json" \
  -d '{"description": "Test backup"}' | jq '.'

# 3. Список бэкапов
echo -e "\n📋 3. Список бэкапов:"
curl -s http://localhost:8090/security/backups | jq '.'

# 4. Тестирование OCR
echo -e "\n📸 4. Распознавание рецепта:"
# Создаем тестовое изображение
if command -v convert &> /dev/null; then
    convert -size 800x600 xc:white -font Arial -pointsize 20 \
      -draw "text 50,50 'Амоксициллин 500 мг'" \
      -draw "text 50,80 'По 1 капсуле 3 раза в день'" \
      -draw "text 50,110 'Врач: Иванова М.А.'" \
      -draw "text 50,140 'Курс: 7 дней'" \
      /tmp/prescription.jpg 2>/dev/null
    
    if [ -f /tmp/prescription.jpg ]; then
        curl -s -X POST http://localhost:8090/api/prescription/scan \
          -F "prescription=@/tmp/prescription.jpg" | jq '.'
        rm /tmp/prescription.jpg
    else
        echo "Не удалось создать тестовое изображение"
    fi
else
    echo "ImageMagick не установлен. Установите: sudo apt-get install imagemagick"
    # Имитируем тест
    curl -s -X POST http://localhost:8090/api/prescription/scan \
      -F "prescription=@/dev/null" 2>&1 | head -5
fi

echo -e "\n✅ Тестирование завершено!"
