#!/bin/bash
echo "💾 СОХРАНЕНИЕ ПРОЕКТА НА GITHUB"
echo "================================"

# Добавляем файлы
git add .
git rm --cached .env 2>/dev/null

# Создаем коммит
git commit -m "🎉 Стабильная версия медицинского бота

- API сервер с JWT
- Telegram бот @NEW_lorhelper_bot
- CRUD для лекарств
- PostgreSQL база данных"

# Отправляем
git push -u origin main

# Создаем тег
git tag -a v1.0.0 -m "Первый стабильный релиз"
git push origin v1.0.0

echo "✅ Проект сохранен на GitHub!"
