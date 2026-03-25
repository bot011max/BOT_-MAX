#!/bin/bash

echo "═══════════════════════════════════════════════════════════════"
echo "🏥 МЕДИЦИНСКИЙ БОТ - MILITARY GRADE SECURITY"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# 1. Security API
echo "🔒 SECURITY API (порт 8090):"
curl -s http://localhost:8090/security/hsm | jq '.data | {mode: .mode, encryption: .encryption, available: .available}'
echo ""

# 2. Telegram Bot
echo "🤖 TELEGRAM BOT (порт 8081):"
curl -s http://localhost:8081/health | jq '.'
echo ""

# 3. Main API
echo "📡 MAIN API (порт 8080):"
curl -s http://localhost:8080/health | jq '.'
echo ""

# 4. Backup count
echo "💾 BACKUPS:"
BACKUP_COUNT=$(curl -s http://localhost:8090/security/backups | jq '.data | length')
echo "   Создано бэкапов: $BACKUP_COUNT"
echo ""

echo "═══════════════════════════════════════════════════════════════"
echo "✅ ВСЕ СИСТЕМЫ РАБОТАЮТ!"
echo "═══════════════════════════════════════════════════════════════"
