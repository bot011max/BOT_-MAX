#!/bin/bash
# ============================================
# MEDICAL BOT - UNIVERSAL LAUNCHER
# Military Grade Security Bot
# Версия: 4.2 - Исправлены пути
# ============================================

# Цвета
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# Определяем директорию скрипта
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

print_banner() {
    echo -e "${CYAN}"
    echo "╔═══════════════════════════════════════════════════════════════════════╗"
    echo "║  🏥 МЕДИЦИНСКИЙ БОТ - MILITARY GRADE SECURITY                        ║"
    echo "║  🚀 УНИВЕРСАЛЬНЫЙ ЗАПУСК ВСЕХ СЕРВИСОВ                                ║"
    echo "╚═══════════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

log_info() { echo -e "${CYAN}[$(date '+%H:%M:%S')] 📌 $1${NC}"; }
log_success() { echo -e "${GREEN}[$(date '+%H:%M:%S')] ✅ $1${NC}"; }
log_warning() { echo -e "${YELLOW}[$(date '+%H:%M:%S')] ⚠️  $1${NC}"; }

# Освобождение портов
free_ports() {
    for port in 8080 8081 8090; do
        if lsof -i:$port > /dev/null 2>&1; then
            log_warning "Порт $port занят, освобождаю..."
            lsof -ti:$port | xargs kill -9 2>/dev/null || true
        fi
    done
    pkill -f "go run" 2>/dev/null || true
    rm -f .api_pid .security_pid .telegram_pid
}

# Создание базы данных
setup_database() {
    mkdir -p data logs
    if [ ! -f "data/medical_bot.db" ]; then
        sqlite3 data/medical_bot.db << 'SQL'
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    role TEXT DEFAULT 'patient',
    phone TEXT,
    is_active BOOLEAN DEFAULT 1,
    telegram_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);
INSERT OR IGNORE INTO users (id, email, password_hash, first_name, last_name, role, is_active)
VALUES (
    'test-user-id',
    'patient@example.com',
    '$2a$10$3xxr7wmozPGp1MdZXT1eAeqeI6QJY29NuWU2GBMBgXdSfEchs00cK',
    'Иван',
    'Петров',
    'patient',
    1
);
SQL
        log_success "База данных создана"
    else
        log_success "База данных существует"
    fi
}

# Запуск Main API
start_main_api() {
    log_info "Запуск Main API (порт 8080)..."
    export JWT_SECRET="medical_bot_super_secret_key_2026_military_grade_32bytes"
    export MASTER_KEY="medical_bot_master_key_for_encryption_2026_32bytes"
    nohup go run cmd/api/main.go > logs/api.log 2>&1 &
    echo $! > .api_pid
    sleep 3
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        log_success "Main API запущен"
    else
        log_warning "Main API не отвечает"
    fi
}

# Запуск Security API
start_security_api() {
    log_info "Запуск Security API (порт 8090)..."
    export JWT_SECRET="medical_bot_super_secret_key_2026_military_grade_32bytes"
    export MASTER_KEY="medical_bot_master_key_for_encryption_2026_32bytes"
    nohup go run cmd/security/main.go > logs/security.log 2>&1 &
    echo $! > .security_pid
    sleep 2
    log_success "Security API запущен"
}

# Запуск Telegram бота
start_telegram_bot() {
    log_info "Запуск Telegram бота (порт 8081)..."
    export JWT_SECRET="medical_bot_super_secret_key_2026_military_grade_32bytes"
    export MASTER_KEY="medical_bot_master_key_for_encryption_2026_32bytes"
    nohup go run cmd/telegram/main.go > logs/telegram.log 2>&1 &
    echo $! > .telegram_pid
    sleep 2
    log_success "Telegram бот запущен"
}

# Проверка статуса
check_status() {
    echo ""
    echo -e "${CYAN}════════════════════════════════════════════════════════════════${NC}"
    echo -e "${CYAN}📊 СТАТУС СЕРВИСОВ:${NC}"
    echo -e "${CYAN}════════════════════════════════════════════════════════════════${NC}"
    
    curl -s http://localhost:8080/health > /dev/null && echo "   📡 Main API:     ✅ РАБОТАЕТ" || echo "   📡 Main API:     ❌ НЕ РАБОТАЕТ"
    curl -s http://localhost:8090/security/hsm > /dev/null && echo "   🔒 Security API: ✅ РАБОТАЕТ" || echo "   🔒 Security API: ❌ НЕ РАБОТАЕТ"
    curl -s http://localhost:8081/health > /dev/null && echo "   🤖 Telegram Bot: ✅ РАБОТАЕТ" || echo "   🤖 Telegram Bot: ❌ НЕ РАБОТАЕТ"
    
    echo -e "${CYAN}════════════════════════════════════════════════════════════════${NC}"
}

# Показ информации
show_info() {
    echo ""
    echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}✅ МЕДИЦИНСКИЙ БОТ ЗАПУЩЕН!${NC}"
    echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
    echo ""
    echo -e "${CYAN}🌐 ДОСТУПНЫЕ СЕРВИСЫ:${NC}"
    echo "   📡 Main API:     http://localhost:8080"
    echo "   🔒 Security API: http://localhost:8090"
    echo "   🤖 Telegram Bot: http://localhost:8081"
    echo ""
    echo -e "${CYAN}🔑 ТЕСТОВЫЕ ДАННЫЕ:${NC}"
    echo "   👤 Email:    patient@example.com"
    echo "   🔐 Пароль:   SecurePass123!"
    echo ""
    echo -e "${CYAN}🚀 ПОЛУЧИТЬ ТОКЕН:${NC}"
    echo "   curl -X POST http://localhost:8080/api/login -H 'Content-Type: application/json' -d '{\"email\":\"patient@example.com\",\"password\":\"SecurePass123!\"}'"
    echo ""
    echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
}

# Основная функция
main() {
    print_banner
    log_info "ЗАПУСК МЕДИЦИНСКОГО БОТА"
    echo ""
    
    free_ports
    setup_database
    start_main_api
    start_security_api
    start_telegram_bot
    
    sleep 2
    check_status
    show_info
    
    log_info "Бот работает. Нажмите Ctrl+C для остановки"
    wait
}

# Обработка Ctrl+C
trap 'echo ""; log_info "Остановка..."; free_ports; exit 0' INT TERM

main "$@"
