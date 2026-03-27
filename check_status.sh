#!/bin/bash
echo "📊 СТАТУС СЕРВИСОВ:"
curl -s http://localhost:8080/health > /dev/null && echo "   Main API: ✅" || echo "   Main API: ❌"
curl -s http://localhost:8090/security/hsm > /dev/null && echo "   Security API: ✅" || echo "   Security API: ❌"
curl -s http://localhost:8081/health > /dev/null && echo "   Telegram Bot: ✅" || echo "   Telegram Bot: ❌"
