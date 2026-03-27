#!/bin/bash

echo "📦 СОХРАНЕНИЕ МЕДИЦИНСКОГО БОТА В GITHUB"
echo "========================================="

cd /workspaces/BOT_MAX

# Проверка git
if ! command -v git &> /dev/null; then
    echo "❌ Git не установлен"
    exit 1
fi

# Инициализация репозитория
if [ ! -d ".git" ]; then
    git init
    echo "✅ Репозиторий инициализирован"
fi

# Добавление .gitignore
cat > .gitignore << 'GITIGNORE'
# Бинарные файлы
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out

# Go файлы
/vendor/
/bin/
/pkg/

# Файлы окружения
.env
.env.local
.env.production

# Логи
*.log
logs/
*.pid
.api_pid
.security_pid
.telegram_pid

# База данных
data/*.db
data/*.db-journal
data/*.db-wal
data/*.db-shm

# Бэкапы
backups/
*.tar.gz
*.zip

# IDE
.vscode/
.idea/
*.swp
*.swo
*~
.DS_Store

# Временные файлы
tmp/
/tmp
*.tmp
GITIGNORE

echo "✅ .gitignore создан"

# Добавление файлов
echo "📁 Добавление файлов..."
git add .

# Статистика
FILES=$(git status --porcelain | wc -l)
echo "✅ Добавлено файлов: $FILES"

# Создание коммита
echo "📝 Создание коммита..."
git commit -m "feat: Working Medical Bot v2.0 - Military Grade Security

✅ JWT аутентификация
✅ CRUD операции с лекарствами
✅ Security API (HSM, бэкапы)
✅ Telegram интеграция
✅ База данных SQLite

🔑 Тестовые данные: patient@example.com / SecurePass123!"

# Создание тега
echo "🏷️ Создание тега..."
git tag -a v2.0.0-working -m "Release v2.0.0 - Working Medical Bot"

# Добавление удаленного репозитория
if ! git remote -v | grep -q origin; then
    git remote add origin https://github.com/bot011max/medical-bot.git
    echo "✅ Удаленный репозиторий добавлен"
fi

# Отправка
echo "📤 Отправка в GitHub..."
git push -u origin main
git push --tags

echo ""
echo "✅ СОХРАНЕНИЕ ЗАВЕРШЕНО!"
echo "🔗 Репозиторий: https://github.com/bot011max/medical-bot"
echo "🏷️ Тег: v2.0.0-working"
