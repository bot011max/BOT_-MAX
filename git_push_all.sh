#!/bin/bash

echo "📦 СОХРАНЕНИЕ ПРОЕКТА В GITHUB"
echo "==============================="
echo ""

# Проверяем настройки
echo "1. Проверка git..."
git status

echo ""
echo "2. Добавление всех файлов..."
git add .

echo ""
echo "3. Создание коммита..."
git commit -m "feat: Medical Bot with Military Grade Security - final version

- HSM hardware encryption (AES-256-GCM)
- Quantum cryptography
- Multi-factor biometrics (voice, face)
- Self-healing backup system
- OCR prescription recognition
- JWT authentication
- Rate limiting & DDoS protection
- Full audit logging
- SQL injection protection
- Security headers (XSS, clickjacking)
- CORS policy"

echo ""
echo "4. Создание тэга v1.0.0..."
git tag -a v1.0.0 -m "Release v1.0.0 - Military Grade Security"

echo ""
echo "5. Отправка на GitHub..."
git push -u origin main
git push --tags

echo ""
echo "✅ Проект сохранен в GitHub!"
echo "   https://github.com/$(git config user.name)/medical-bot"
