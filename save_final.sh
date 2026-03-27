#!/bin/bash

echo "📦 СОХРАНЕНИЕ ФИНАЛЬНОЙ ВЕРСИИ В GITHUB"
echo "========================================"

cd /workspaces/BOT_MAX/BOT_MAX

# Проверка git
if ! command -v git &> /dev/null; then
    echo "❌ Git не установлен"
    exit 1
fi

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

# Создание коммита
echo "📝 Создание коммита..."
git commit -m "feat: FINAL WORKING VERSION v3.0 - Medical Bot with Military Grade Security

## 🚀 РАБОТАЮЩИЕ СЕРВИСЫ
- Main API (порт 8080)
- Security API (порт 8090)
- Telegram Bot (порт 8081)
- База данных SQLite

## 🔧 ФУНКЦИОНАЛ
- JWT аутентификация
- Регистрация и логин
- CRUD операции с лекарствами
- Security API (HSM, бэкапы, OCR)
- Полный аудит
- Автоматическое восстановление

## 🛡️ АКТИВНЫЕ ЗАЩИТЫ
- Квантово-устойчивая криптография
- HSM шифрование (AES-256-GCM)
- Rate Limiting
- SQL Injection Protection
- XSS Protection
- Security Headers

## 📁 СОЗДАННЫЕ СКРИПТЫ
- install_and_run.sh - полная установка
- run_working.sh - быстрый запуск
- demo.sh - демонстрация
- check_status.sh - проверка статуса
- stop_bot.sh - остановка

## 🔑 ТЕСТОВЫЕ ДАННЫЕ
- Email: patient@example.com
- Password: SecurePass123!"

# Создание тега
echo "🏷️ Создание тега..."
git tag -a v3.0.0 -m "Release v3.0.0 - FINAL WORKING VERSION"

# Отправка
echo "📤 Отправка в GitHub..."
git push origin main
git push --tags

echo ""
echo "✅ СОХРАНЕНИЕ ЗАВЕРШЕНО!"
echo "🔗 Репозиторий: https://github.com/bot011max/BOT_MAX"
echo "🏷️ Тег: v3.0.0"
echo ""
echo "🚀 КЛОНИРОВАНИЕ НА НОВОМ КОМПЬЮТЕРЕ:"
echo "   git clone https://github.com/bot011max/BOT_MAX.git"
echo "   cd BOT_MAX"
echo "   chmod +x install_and_run.sh"
echo "   ./install_and_run.sh"
