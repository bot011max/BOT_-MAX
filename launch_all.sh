#!/bin/bash
# ============================================
# MEDICAL BOT - UNIVERSAL LAUNCHER
# Military Grade Security Bot
# Полный запуск всех сервисов
# Версия: 4.0
# ============================================

# Цвета
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m'

print_banner() {
    echo -e "${CYAN}"
    echo "╔═══════════════════════════════════════════════════════════════════════╗"
    echo "║  🏥 МЕДИЦИНСКИЙ БОТ - MILITARY GRADE SECURITY                        ║"
    echo "║  🚀 УНИВЕРСАЛЬНЫЙ ЗАПУСК ВСЕХ СЕРВИСОВ                                ║"
    echo "╚═══════════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

log() { echo -e "${BLUE}[$(date '+%H:%M:%S')]${NC} $1"; }
log_success() { echo -e "${GREEN}[$(date '+%H:%M:%S')] ✅ $1${NC}"; }
log_error() { echo -e "${RED}[$(date '+%H:%M:%S')] ❌ $1${NC}"; }
log_warning() { echo -e "${YELLOW}[$(date '+%H:%M:%S')] ⚠️  $1${NC}"; }
log_info() { echo -e "${CYAN}[$(date '+%H:%M:%S')] 📌 $1${NC}"; }

# Создание базы данных
setup_database() {
    log_info "Настройка базы данных..."
    
    mkdir -p data logs
    
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

CREATE TABLE IF NOT EXISTS medications (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    dosage TEXT,
    frequency TEXT,
    instructions TEXT,
    start_date DATETIME,
    end_date DATETIME,
    is_active BOOLEAN DEFAULT 1,
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
    log_success "Переменные окружения настроены"
}

# Запуск Main API
start_main_api() {
    log_info "Запуск Main API (порт 8080)..."
    cd /workspaces/BOT_MAX/BOT_MAX
    nohup go run cmd/api/main.go > logs/api.log 2>&1 &
    echo $! > .api_pid
    sleep 3
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        log_success "Main API запущен"
        return 0
    else
        log_warning "Main API не отвечает, проверьте логи: tail -f logs/api.log"
        return 1
    fi
}

# Запуск Security API
start_security_api() {
    log_info "Запуск Security API (порт 8090)..."
    cd /workspaces/BOT_MAX/BOT_MAX
    nohup go run cmd/security/main.go > logs/security.log 2>&1 &
    echo $! > .security_pid
    sleep 2
    log_success "Security API запущен"
}

# Запуск Telegram бота
start_telegram_bot() {
    log_info "Запуск Telegram бота (порт 8081)..."
    cd /workspaces/BOT_MAX/BOT_MAX
    nohup go run cmd/telegram/main.go > logs/telegram.log 2>&1 &
    echo $! > .telegram_pid
    sleep 2
    log_success "Telegram бот запущен"
}

# Финальная проверка
final_check() {
    echo ""
    echo -e "${CYAN}════════════════════════════════════════════════════════════════${NC}"
    echo -e "${CYAN}📊 ПРОВЕРКА СЕРВИСОВ:${NC}"
    echo -e "${CYAN}════════════════════════════════════════════════════════════════${NC}"
    
    # Main API
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo -e "   📡 Main API:     ${GREEN}✅ РАБОТАЕТ${NC}"
    else
        echo -e "   📡 Main API:     ${RED}❌ НЕ РАБОТАЕТ${NC}"
    fi
    
    # Security API
    if curl -s http://localhost:8090/security/hsm > /dev/null 2>&1; then
        echo -e "   🔒 Security API: ${GREEN}✅ РАБОТАЕТ${NC}"
    else
        echo -e "   🔒 Security API: ${YELLOW}⚠️  НЕ ОТВЕЧАЕТ${NC}"
    fi
    
    # Telegram Bot
    if curl -s http://localhost:8081/health > /dev/null 2>&1; then
        echo -e "   🤖 Telegram Bot: ${GREEN}✅ РАБОТАЕТ${NC}"
    else
        echo -e "   🤖 Telegram Bot: ${YELLOW}⚠️  НЕ ОТВЕЧАЕТ${NC}"
    fi
    
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
    echo -e "${CYAN}🚀 БЫСТРЫЕ КОМАНДЫ:${NC}"
    echo "   🔐 Получить токен:"
    echo "      curl -X POST http://localhost:8080/api/login -H 'Content-Type: application/json' -d '{\"email\":\"patient@example.com\",\"password\":\"SecurePass123!\"}'"
    echo ""
    echo -e "${CYAN}📝 ЛОГИ:${NC}"
    echo "   tail -f logs/api.log"
    echo "   tail -f logs/security.log"
    echo "   tail -f logs/telegram.log"
    echo ""
    echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
}

# Остановка
stop_services() {
    echo ""
    echo "🛑 Остановка сервисов..."
    for f in .api_pid .security_pid .telegram_pid; do
        [ -f "$f" ] && kill $(cat "$f") 2>/dev/null && rm -f "$f"
    done
    pkill -f "go run" 2>/dev/null
    echo "✅ Бот остановлен"
}

# Основная функция
main() {
    print_banner
    log_info "ЗАПУСК МЕДИЦИНСКОГО БОТА"
    echo ""
    
    # Остановка предыдущих
    stop_services
    sleep 1
    
    # Подготовка
    setup_database
    setup_env
    
    # Запуск
    start_main_api
    start_security_api
    start_telegram_bot
    
    # Проверка
    sleep 3
    final_check
    show_info
    
    echo ""
    log_info "Бот работает. Нажмите Ctrl+C для остановки"
    
    # Ожидание
    wait
}

# Обработка Ctrl+C
trap 'stop_services; exit 0' INT TERM

main "$@"
