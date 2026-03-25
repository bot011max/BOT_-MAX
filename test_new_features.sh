#!/bin/bash
echo "🧪 ТЕСТИРОВАНИЕ НОВЫХ ФУНКЦИЙ"
echo "================================"

# 1. Тест HSM
echo -e "\n1️⃣ Аппаратное шифрование (HSM):"
curl -s http://localhost:8080/security/hsm | jq '.'

# 2. Создание бэкапа
echo -e "\n2️⃣ Создание бэкапа:"
curl -s -X POST http://localhost:8080/security/backup \
  -H "Content-Type: application/json" \
  -d '{"description": "Test backup"}' | jq '.'

# 3. Список бэкапов
echo -e "\n3️⃣ Список бэкапов:"
curl -s http://localhost:8080/security/backups | jq '.'

# 4. Тест OCR рецепта
echo -e "\n4️⃣ Анализ фото рецепта:"
# Создаем тестовое изображение
convert -size 800x600 xc:white -font Arial -pointsize 20 \
  -draw "text 50,50 'РЕЦЕПТ № 12345'" \
  -draw "text 50,80 'Врач: Иванова М.А.'" \
  -draw "text 50,110 'Амоксициллин 500 мг'" \
  -draw "text 50,140 'По 1 капсуле 3 раза в день'" \
  /tmp/test_prescription.jpg 2>/dev/null

if [ -f /tmp/test_prescription.jpg ]; then
    curl -s -X POST http://localhost:8080/api/prescription/scan \
      -F "prescription=@/tmp/test_prescription.jpg" | jq '.'
    rm /tmp/test_prescription.jpg
else
    echo "ImageMagick not installed, skipping image test"
    curl -s -X POST http://localhost:8080/api/prescription/scan \
      -F "prescription=@/dev/null" 2>/dev/null | jq '.'
fi

echo -e "\n✅ Тестирование завершено!"
