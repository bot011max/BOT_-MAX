#!/bin/bash

echo "💾 СОХРАНЕНИЕ ПРОЕКТА"
echo "====================="
echo ""

# Создаем директорию для бэкапов
mkdir -p backups

# Генерируем имя файла
BACKUP_NAME="medical_bot_$(date +%Y%m%d_%H%M%S)"

echo "1. Сохраняем в Git..."
git add .
git commit -m "Backup: $BACKUP_NAME" 2>/dev/null || echo "   Нет изменений для коммита"

echo "2. Создаем архив проекта..."
tar -czf "backups/${BACKUP_NAME}.tar.gz" \
  --exclude="*.log" \
  --exclude=".api_pid" \
  --exclude=".security_pid" \
  --exclude=".telegram_pid" \
  --exclude="backups" \
  --exclude=".git" \
  .

echo "3. Сохраняем зависимости..."
go list -m all > "backups/${BACKUP_NAME}_deps.txt"

echo "4. Сохраняем конфигурацию..."
cp .env "backups/${BACKUP_NAME}_env.txt" 2>/dev/null || echo "   .env не найден"

echo ""
echo "✅ Проект сохранен в: backups/${BACKUP_NAME}.tar.gz"
echo "   Размер: $(du -h backups/${BACKUP_NAME}.tar.gz | cut -f1)"
echo ""
echo "📋 Список бэкапов:"
ls -lh backups/
