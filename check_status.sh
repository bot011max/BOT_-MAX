#!/bin/bash
echo "📊 СТАТУС СЕРВИСОВ:"
echo -n "   Main API: "
curl -s http://localhost:8080/health > /dev/null 2>&1 && echo "✅" || echo "❌"
echo -n "   Security API: "
curl -s http://localhost:8090/security/hsm > /dev/null 2>&1 && echo "✅" || echo "❌"
echo -n "   Telegram Bot: "
curl -s http://localhost:8081/health > /dev/null 2>&1 && echo "✅" || echo "❌"
