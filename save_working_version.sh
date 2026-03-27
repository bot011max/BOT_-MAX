#!/bin/bash

echo "📦 СОХРАНЕНИЕ РАБОТОСПОСОБНОЙ ВЕРСИИ В GITHUB"
echo "=============================================="

cd /workspaces/BOT_MAX

# Проверка git
if ! command -v git &> /dev/null; then
    echo "❌ Git не установлен"
    exit 1
fi

echo "✅ Git: $(git --version)"

# Проверка изменений
if [ -z "$(git status --porcelain)" ]; then
    echo "⚠️  Нет изменений для сохранения"
    exit 0
fi

# Показать изменения
echo ""
echo "📝 ИЗМЕНЕНИЯ:"
git status --short
echo ""

# Добавление файлов
echo "📁 Добавление файлов..."
git add .

# Статистика
FILES=$(git status --porcelain | wc -l)
echo "✅ Добавлено файлов: $FILES"

# Создание коммита
echo "📝 Создание коммита..."
git commit -m "feat: FINAL WORKING VERSION v5.0 - Medical Bot with Military Grade Security

✅ Все функции работают:
- Main API (8080)
- Security API (8090)
- Telegram Bot (8081)
- JWT аутентификация
- CRUD операции с лекарствами
- База данных SQLite
- Полный аудит

🔑 Тестовые данные: patient@example.com / SecurePass123!

🚀 Запуск: ./start_bot_full.sh"

# Создание тега
echo "🏷️ Создание тега..."
git tag -a v5.0.0 -m "Release v5.0.0 - Working Medical Bot"

# Отправка
echo "📤 Отправка в GitHub..."
git push origin main
git push --tags

echo ""
echo "✅ СОХРАНЕНИЕ ЗАВЕРШЕНО!"
echo "🔗 Репозиторий: https://github.com/bot011max/BOT_MAX"
echo "🏷️ Тег: v5.0.0"
echo ""
echo "🚀 КЛОНИРОВАНИЕ НА ДРУГОМ КОМПЬЮТЕРЕ:"
echo "   git clone https://github.com/bot011max/BOT_MAX.git"
echo "   cd BOT_MAX"
echo "   chmod +x start_bot_full.sh"
echo "   ./start_bot_full.sh"
