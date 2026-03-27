#!/bin/bash
# ============================================
# MEDICAL BOT - UNIVERSAL AUTO-START SCRIPT
# Исправленная версия - без set -e
# ============================================

# Цвета
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m'

# Конфигурация
API_PORT=8080
SECURITY_PORT=8090
TELEGRAM_PORT=8081
LOG_DIR="logs"
DATA_DIR="data"

print_banner() {
    echo -e "${CYAN}"
    echo "╔═══════════════════════════════════════════════════════════════════════╗"
    echo "║  🏥 МЕДИЦИНСКИЙ БОТ - MILITARY GRADE SECURITY                        ║"
    echo "║  🚀 АВТОМАТИЧЕСКИЙ ЗАПУСК                                             ║"
    echo "╚═══════════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

log() { echo -e "${BLUE}[$(date '+%H:%M:%S')]${NC} $1"; }
log_success() { echo -e "${GREEN}[$(date '+%H:%M:%S')] ✅ $1${NC}"; }
log_error() { echo -e "${RED}[$(date '+%H:%M:%S')] ❌ $1${NC}"; }
log_info() { echo -e "${CYAN}[$(date '+%H:%M:%S')] 📌 $1${NC}"; }

# Определение директории
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR" || cd /workspaces/BOT_MAX

# Остановка процессов
stop_services() {
    log_info "Остановка сервисов..."
    pkill -f "go run" 2>/dev/null || true
    for port in 8080 8081 8090; do
        lsof -ti:$port | xargs kill -9 2>/dev/null || true
    done
    rm -f .api_pid .security_pid .telegram_pid
    log_success "Очистка завершена"
}

# Создание базы данных
create_database() {
    log_info "Создание базы данных..."
    mkdir -p "$DATA_DIR"
    
    sqlite3 "$DATA_DIR/medical_bot.db" << 'SQL'
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
    log_success "База данных готова"
}

# Настройка переменных
setup_env() {
    export JWT_SECRET="medical_bot_super_secret_key_2026_military_grade_32bytes"
    export MASTER_KEY="medical_bot_master_key_for_encryption_2026_32bytes"
}

# Запуск Main API
start_main_api() {
    log_info "Запуск Main API (порт $API_PORT)..."
    cd "$SCRIPT_DIR"
    nohup go run cmd/api/main.go > "$LOG_DIR/api.log" 2>&1 &
    echo $! > .api_pid
    sleep 3
    if curl -s http://localhost:$API_PORT/health > /dev/null 2>&1; then
        log_success "Main API запущен"
        return 0
    else
        log_warning "Main API не отвечает, проверьте логи"
        return 1
    fi
}

# Запуск Security API
start_security_api() {
    log_info "Запуск Security API (порт $SECURITY_PORT)..."
    cd "$SCRIPT_DIR"
    nohup go run cmd/security/main.go > "$LOG_DIR/security.log" 2>&1 &
    echo $! > .security_pid
    sleep 2
    log_success "Security API запущен"
}

# Запуск Telegram бота
start_telegram_bot() {
    log_info "Запуск Telegram бота (порт $TELEGRAM_PORT)..."
    cd "$SCRIPT_DIR"
    nohup go run cmd/telegram/main.go > "$LOG_DIR/telegram.log" 2>&1 &
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
    
    curl -s http://localhost:$API_PORT/health > /dev/null && echo "   📡 Main API:     ✅ РАБОТАЕТ" || echo "   📡 Main API:     ❌ НЕ РАБОТАЕТ"
    curl -s http://localhost:$SECURITY_PORT/security/hsm > /dev/null && echo "   🔒 Security API: ✅ РАБОТАЕТ" || echo "   🔒 Security API: ❌ НЕ РАБОТАЕТ"
    curl -s http://localhost:$TELEGRAM_PORT/health > /dev/null && echo "   🤖 Telegram Bot: ✅ РАБОТАЕТ" || echo "   🤖 Telegram Bot: ❌ НЕ РАБОТАЕТ"
    
    echo -e "${CYAN}════════════════════════════════════════════════════════════════${NC}"
}

# Демонстрация
run_demo() {
    echo ""
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${CYAN}🏥 ДЕМОНСТРАЦИЯ${NC}"
    echo -e "${CYAN}═══════════════════════════════════════════════════════════════${NC}"
    
    echo "1. Health Check:"
    curl -s http://localhost:$API_PORT/health | jq '.'
    
    echo ""
    echo "2. JWT токен:"
    LOGIN=$(curl -s -X POST http://localhost:$API_PORT/api/login -H "Content-Type: application/json" -d '{"email":"patient@example.com","password":"SecurePass123!"}')
    echo $LOGIN | jq '.data | {token_preview: (.token[:50] + "..."), user: .user}'
    TOKEN=$(echo $LOGIN | jq -r '.data.token')
    
    echo ""
    echo "3. Добавление лекарства:"
    curl -s -X POST http://localhost:$API_PORT/api/medications -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{"name":"Аспирин","dosage":"100 мг","frequency":"1 раз в день"}' | jq '.data | {name: .name, dosage: .dosage}'
    
    echo ""
    echo "4. Список лекарств:"
    curl -s -X GET http://localhost:$API_PORT/api/medications -H "Authorization: Bearer $TOKEN" | jq '.data[] | "   • \(.name) - \(.dosage) (\(.frequency))"'
    
    echo ""
    echo -e "${GREEN}✅ ДЕМОНСТРАЦИЯ ЗАВЕРШЕНА${NC}"
}

# Показ информации
show_info() {
    echo ""
    echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}✅ МЕДИЦИНСКИЙ БОТ ЗАПУЩЕН!${NC}"
    echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
    echo ""
    echo -e "${CYAN}🌐 СЕРВИСЫ:${NC}"
    echo "   📡 Main API:     http://localhost:$API_PORT"
    echo "   🔒 Security API: http://localhost:$SECURITY_PORT"
    echo "   🤖 Telegram Bot: http://localhost:$TELEGRAM_PORT"
    echo ""
    echo -e "${CYAN}🔑 ТЕСТОВЫЕ ДАННЫЕ:${NC}"
    echo "   👤 Email: patient@example.com"
    echo "   🔐 Пароль: SecurePass123!"
    echo ""
    echo -e "${CYAN}🛑 ОСТАНОВКА:${NC}"
    echo "   pkill -f 'go run'"
    echo "   или ./stop_bot.sh"
    echo ""
}

# Главная функция
main() {
    print_banner
    log_info "ЗАПУСК МЕДИЦИНСКОГО БОТА"
    echo ""
    
    mkdir -p "$LOG_DIR"
    stop_services
    create_database
    setup_env
    
    start_main_api
    start_security_api
    start_telegram_bot
    
    sleep 3
    check_status
    run_demo
    show_info
    
    log_info "Бот работает. Нажмите Ctrl+C для остановки"
    wait
}

# Обработка сигналов
trap 'echo ""; stop_services; exit 0' INT TERM

main "$@"
