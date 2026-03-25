#!/bin/bash

echo "📦 СОХРАНЕНИЕ ПРОЕКТА В GIT"
echo "============================"
echo ""

# Проверяем наличие .git
if [ ! -d .git ]; then
    echo "Инициализация Git репозитория..."
    git init
fi

# Проверяем текущую ветку
BRANCH=$(git branch --show-current)
if [ -z "$BRANCH" ]; then
    BRANCH="main"
    git checkout -b $BRANCH
fi

echo "Текущая ветка: $BRANCH"

# Добавляем все изменения
echo ""
echo "Добавление файлов..."
git add .

# Показываем что будет закоммичено
echo ""
echo "Файлы для коммита:"
git status -s

# Запрашиваем описание коммита
echo ""
read -p "Введите описание коммита: " COMMIT_MSG

if [ -z "$COMMIT_MSG" ]; then
    COMMIT_MSG="Update: $(date '+%Y-%m-%d %H:%M:%S')"
fi

# Создаем коммит
echo ""
echo "Создание коммита..."
git commit -m "$COMMIT_MSG"

# Показываем последний коммит
echo ""
echo "Последний коммит:"
git log -1 --oneline

# Спрашиваем о push
echo ""
read -p "Отправить на удаленный репозиторий? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    # Проверяем наличие удаленного репозитория
    if ! git remote -v | grep -q origin; then
        echo "Удаленный репозиторий не настроен"
        read -p "Введите URL удаленного репозитория: " REMOTE_URL
        git remote add origin $REMOTE_URL
    fi
    
    echo "Отправка изменений..."
    git push -u origin $BRANCH
    git push --tags
    echo "✅ Изменения отправлены"
fi

echo ""
echo "✅ Готово!"
echo "История коммитов:"
git log --oneline -5
