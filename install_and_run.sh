#!/bin/bash
# ============================================
# MEDICAL BOT - UNIVERSAL INSTALLER
# Military Grade Security Bot
# Полная автоматическая установка и запуск
# Версия: 3.0 - FINAL
# ============================================

set -e

# Цвета
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

print_banner() {
    echo -e "${CYAN}"
    echo "╔═══════════════════════════════════════════════════════════════════════╗"
    echo "║  🏥 МЕДИЦИНСКИЙ БОТ - MILITARY GRADE SECURITY                        ║"
    echo "║  🚀 АВТОМАТИЧЕСКАЯ УСТАНОВКА И ЗАПУСК                                 ║"
    echo "║  🔒 End-to-End Encryption | HSM | Quantum Crypto                     ║"
    echo "╚═══════════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

log_info() { echo -e "${CYAN}[$(date '+%H:%M:%S')] 📌 $1${NC}"; }
log_success() { echo -e "${GREEN}[$(date '+%H:%M:%S')] ✅ $1${NC}"; }
log_error() { echo -e "${RED}[$(date '+%H:%M:%S')] ❌ $1${NC}"; }

# Проверка и установка Go
check_go() {
    if ! command -v go &> /dev/null; then
        log_info "Установка Go..."
        sudo apt-get update && sudo apt-get install -y golang-go
    fi
    log_success "Go: $(go version | cut -d' ' -f3)"
}

# Проверка и установка Git
check_git() {
    if ! command -v git &> /dev/null; then
        log_info "Установка Git..."
        sudo apt-get install -y git
    fi
    log_success "Git: $(git --version | cut -d' ' -f3)"
}

# Проверка и установка SQLite
check_sqlite() {
    if ! command -v sqlite3 &> /dev/null; then
        log_info "Установка SQLite..."
        sudo apt-get install -y sqlite3
    fi
    log_success "SQLite: $(sqlite3 --version | cut -d' ' -f1)"
}

# Установка зависимостей Go
install_deps() {
    log_info "Установка зависимостей Go..."
    go mod download
    go mod tidy
    log_success "Зависимости установлены"
}

# Создание базы данных
setup_database() {
    log_info "Создание базы данных..."
    mkdir -p data logs
    sqlite3 data/medical_bot.db << SQL
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
    log_success "База данных создана"
}

# Настройка переменных окружения
setup_env() {
    log_info "Настройка переменных окружения..."
    export JWT_SECRET="medical_bot_super_secret_key_2026_military_grade_32bytes"
    export MASTER_KEY="medical_bot_master_key_for_encryption_2026_32bytes"
    log_success "Переменные окружения настроены"
}

# Запуск сервисов
start_services() {
    log_info "Запуск сервисов..."
    pkill -f "go run" 2>/dev/null || true
    
    go run cmd/api/main.go > logs/api.log 2>&1 &
    echo $! > .api_pid
    
    go run cmd/security/main.go > logs/security.log 2>&1 &
    echo $! > .security_pid
    
    go run cmd/telegram/main.go > logs/telegram.log 2>&1 &
    echo $! > .telegram_pid
    
    sleep 5
    log_success "Сервисы запущены"
}

# Создание вспомогательных скриптов
create_scripts() {
    log_info "Создание вспомогательных скриптов..."
    
    cat > run_working.sh << 'RUN'
#!/bin/bash
cd /workspaces/BOT_MAX/BOT_MAX
export JWT_SECRET="medical_bot_super_secret_key_2026_military_grade_32bytes"
export MASTER_KEY="medical_bot_master_key_for_encryption_2026_32bytes"
pkill -f "go run" 2>/dev/null || true
go run cmd/api/main.go > logs/api.log 2>&1 &
echo $! > .api_pid
go run cmd/security/main.go > logs/security.log 2>&1 &
echo $! > .security_pid
go run cmd/telegram/main.go > logs/telegram.log 2>&1 &
echo $! > .telegram_pid
sleep 3
echo "✅ Бот запущен!"
RUN
    chmod +x run_working.sh

    cat > stop_bot.sh << 'STOP'
#!/bin/bash
for f in .api_pid .security_pid .telegram_pid; do
    [ -f "$f" ] && kill $(cat "$f") 2>/dev/null
    rm -f "$f"
done
pkill -f "go run" 2>/dev/null
echo "✅ Бот остановлен"
STOP
    chmod +x stop_bot.sh

    cat > check_status.sh << 'CHECK'
#!/bin/bash
echo "📊 СТАТУС СЕРВИСОВ:"
curl -s http://localhost:8080/health > /dev/null && echo "   Main API: ✅" || echo "   Main API: ❌"
curl -s http://localhost:8090/security/hsm > /dev/null && echo "   Security API: ✅" || echo "   Security API: ❌"
curl -s http://localhost:8081/health > /dev/null && echo "   Telegram Bot: ✅" || echo "   Telegram Bot: ❌"
CHECK
    chmod +x check_status.sh

    cat > demo.sh << 'DEMO'
#!/bin/bash
echo "═══════════════════════════════════════════════════════════════"
echo "🏥 МЕДИЦИНСКИЙ БОТ - ДЕМОНСТРАЦИЯ"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "1. Health Check:"
curl -s http://localhost:8080/health | jq '.'
echo ""
echo "2. Получение JWT токена:"
LOGIN=$(curl -s -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d '{"email":"patient@example.com","password":"SecurePass123!"}')
echo $LOGIN | jq '.data | {token_preview: (.token[:50] + "..."), user: .user}'
TOKEN=$(echo $LOGIN | jq -r '.data.token')
echo ""
echo "3. Добавление лекарства:"
curl -s -X POST http://localhost:8080/api/medications -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{"name":"Аспирин","dosage":"100 мг","frequency":"1 раз в день"}' | jq '.data | {name: .name, dosage: .dosage}'
echo ""
echo "4. Список лекарств:"
curl -s -X GET http://localhost:8080/api/medications -H "Authorization: Bearer $TOKEN" | jq '.data[] | "   • \(.name) - \(.dosage) (\(.frequency))"'
echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "✅ ДЕМОНСТРАЦИЯ ЗАВЕРШЕНА"
echo "═══════════════════════════════════════════════════════════════"
DEMO
    chmod +x demo.sh
    
    log_success "Скрипты созданы"
}

# Финальная проверка
final_check() {
    log_info "Финальная проверка..."
    sleep 2
    
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        log_success "Main API работает"
    else
        log_error "Main API не отвечает"
    fi
    
    if curl -s http://localhost:8090/security/hsm > /dev/null 2>&1; then
        log_success "Security API работает"
    else
        log_error "Security API не отвечает"
    fi
    
    if curl -s http://localhost:8081/health > /dev/null 2>&1; then
        log_success "Telegram бот работает"
    else
        log_error "Telegram бот не отвечает"
    fi
}

# Главная функция
main() {
    print_banner
    log_info "НАЧАЛО УСТАНОВКИ НА НОВОМ КОМПЬЮТЕРЕ"
    echo ""
    
    check_go
    check_git
    check_sqlite
    install_deps
    setup_database
    setup_env
    start_services
    create_scripts
    final_check
    
    echo ""
    echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}✅ МЕДИЦИНСКИЙ БОТ УСПЕШНО УСТАНОВЛЕН И ЗАПУЩЕН!${NC}"
    echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
    echo ""
    echo -e "${CYAN}🌐 ДОСТУПНЫЕ СЕРВИСЫ:${NC}"
    echo "   📡 Main API:     http://localhost:8080"
    echo "   🔒 Security API: http://localhost:8090"
    echo "   🤖 Telegram Bot: http://localhost:8081"
    echo ""
    echo -e "${CYAN}🔑 ТЕСТОВЫЕ ДАННЫЕ:${NC}"
    echo "   👤 Email: patient@example.com"
    echo "   🔐 Пароль: SecurePass123!"
    echo ""
    echo -e "${CYAN}🚀 БЫСТРЫЕ КОМАНДЫ:${NC}"
    echo "   ./run_working.sh   - быстрый запуск"
    echo "   ./demo.sh          - полная демонстрация"
    echo "   ./check_status.sh  - проверка статуса"
    echo "   ./stop_bot.sh      - остановка бота"
    echo ""
}

main "$@"
