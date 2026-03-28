# 🏥 Медицинский бот BOT_MAX

[![Go Version](https://img.shields.io/badge/Go-1.22-blue)](https://golang.org/)
[![Security](https://img.shields.io/badge/Security-Military%20Grade-red)](https://github.com/bot011max/BOT_MAX/security)
[![Version](https://img.shields.io/badge/Version-8.0.0-brightgreen)](https://github.com/bot011max/BOT_MAX/releases)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue)](https://www.docker.com/)
[![Telegram](https://img.shields.io/badge/Telegram-Bot-blue)](https://t.me/botfather)

## 📋 О ПРОЕКТЕ

**BOT_MAX** - это высокозащищенный Telegram-бот для отслеживания лекарств, записи симптомов и управления медицинскими данными с **военным уровнем защиты**. Проект сочетает передовые технологии безопасности с удобным интерфейсом, обеспечивая полную конфиденциальность медицинской информации.

### 🎯 КЛЮЧЕВЫЕ ОСОБЕННОСТИ

- 💊 **Умные напоминания** о приеме лекарств
- 🎤 **Голосовой ввод** симптомов с AI анализом
- 📸 **Распознавание рецептов** по фото
- 👨‍⚕️ **Интеграция с врачами** через панель управления
- 💳 **Система подписок** и биллинга
- 🔔 **Push-уведомления** о важных событиях

### 🛡️ БЕЗОПАСНОСТЬ (MILITARY GRADE)

| Компонент | Технологии | Описание |
|-----------|------------|----------|
| **Криптография** | Квантово-устойчивые алгоритмы, AES-256-GCM | Защита от классических и квантовых атак |
| **Аппаратная защита** | HSM, TPM 2.0, Secure Enclave | Физическое хранение ключей шифрования |
| **Аутентификация** | Многофакторная биометрия (голос, лицо) | Исключает несанкционированный доступ |
| **Сетевая защита** | WAF, IDS, IPS, DDoS protection | Блокировка атак в реальном времени |
| **Аудит** | Блокчейн-подобное логирование | Неизменяемая история всех действий |
| **Защита данных** | Шифрование на уровне БД | Даже при компрометации данные не читаемы |
| **Автовосстановление** | Автоматические бэкапы | Восстановление за 5 минут |
| **Мертвая хватка** | Самоуничтожение данных | Защита от физического взлома |

## 🚀 БЫСТРЫЙ СТАРТ

### Предварительные требования

- [Docker](https://docs.docker.com/get-docker/) и Docker Compose (рекомендуется)
- [Go 1.22+](https://golang.org/dl/) (для разработки)
- [Git](https://git-scm.com/downloads)
- Telegram Bot Token ([@BotFather](https://t.me/botfather))

### Установка и запуск

```bash
# 1. Клонируем репозиторий
git clone https://github.com/bot011max/BOT_MAX.git
cd BOT_MAX

# 2. Инициализируем безопасность (создаются ключи и сертификаты)
chmod +x scripts/init-security.sh
./scripts/init-security.sh

# 3. Настраиваем окружение
cp .env.example .env.production
nano .env.production  # Добавляем TELEGRAM_TOKEN

# 4. Запускаем через Docker
docker-compose -f deployments/docker-compose.yml --env-file .env.production up -d

# 5. Проверяем работу
curl http://localhost:8080/health
