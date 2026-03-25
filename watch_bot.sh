#!/bin/bash
while true; do
    clear
    echo "=== МОНИТОРИНГ МЕДИЦИНСКОГО БОТА ==="
    echo "Время: $(date '+%H:%M:%S')"
    echo ""
    
    # API статус
    API_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)
    if [ "$API_STATUS" = "200" ]; then
        echo "✅ API: OK (http://localhost:8080)"
    else
        echo "❌ API: Error ($API_STATUS)"
    fi
    
    # Telegram статус
    TG_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8081/health)
    if [ "$TG_STATUS" = "200" ]; then
        echo "✅ Telegram Bot: OK (http://localhost:8081)"
    else
        echo "❌ Telegram Bot: Error ($TG_STATUS)"
    fi
    
    # Безопасность
    echo ""
    echo "🔐 Security Status:"
    curl -s http://localhost:8080/security/status | jq -c '.'
    
    # Последние логи
    echo ""
    echo "📝 Последние логи:"
    tail -3 logs/bot.log 2>/dev/null || echo "Логи не найдены"
    
    sleep 3
done
