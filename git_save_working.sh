#!/bin/bash
# ============================================
# GIT SAVE - Medical Bot Working Version
# Сохранение работоспособной версии в GitHub
# Версия: 5.0 - FINAL WORKING
# ============================================

# Цвета
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

# Конфигурация
REPO_URL="https://github.com/bot011max/medical-bot.git"
BRANCH="main"
VERSION="v2.0.0-working"
COMMIT_DATE=$(date '+%Y-%m-%d %H:%M:%S')

print_banner() {
    echo -e "${CYAN}"
    cat << "EOF"
╔═══════════════════════════════════════════════════════════════════════╗
║  📦 GIT SAVE - WORKING VERSION                                       ║
║  🚀 Сохранение работоспособной версии медицинского бота              ║
║  🔒 Military Grade Security - Fully Functional                       ║
╚═══════════════════════════════════════════════════════════════════════╝
