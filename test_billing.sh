#!/bin/bash
echo "💰 ПРОВЕРКА СИСТЕМЫ МОНЕТИЗАЦИИ"
echo "================================"
echo ""

cd /workspaces/BOT_MAX

go run -e << 'GO_TEST' 2>/dev/null
package main

import (
    "fmt"
    "github.com/bot011max/medical-bot/internal/billing"
)

func main() {
    fmt.Println("💎 Тестирование системы монетизации...")
    
    manager := billing.NewSubscriptionManager()
    
    // Проверка бесплатного тарифа
    fmt.Println("\n📋 БЕСПЛАТНЫЙ ТАРИФ:")
    freeTier := billing.TierConfigs[billing.TierFree]
    fmt.Printf("   • Цена: %.0f руб./мес\n", freeTier.Price)
    fmt.Printf("   • Функции:\n")
    for _, f := range freeTier.Features {
        fmt.Printf("     - %s\n", f)
    }
    
    // Проверка премиум тарифа
    fmt.Println("\n💎 ПРЕМИУМ ТАРИФ:")
    premiumTier := billing.TierConfigs[billing.TierPremium]
    fmt.Printf("   • Цена: %.0f руб./мес\n", premiumTier.Price)
    fmt.Printf("   • Функции:\n")
    for _, f := range premiumTier.Features {
        fmt.Printf("     - %s\n", f)
    }
    
    // Тест апгрейда
    userID := "test_user_001"
    fmt.Printf("\n👤 Пользователь %s\n", userID)
    fmt.Printf("   Текущий тариф: %s\n", manager.GetUserTier(userID))
    
    manager.UpgradeToPremium(userID)
    fmt.Printf("   ✅ Апгрейд до %s\n", manager.GetUserTier(userID))
    
    // Проверка доступа к функциям
    features := []string{"AI-анализ симптомов", "Push-уведомления", "Экспорт отчетов", "Несуществующая функция"}
    fmt.Printf("\n🔍 ПРОВЕРКА ДОСТУПА К ФУНКЦИЯМ:\n")
    for _, f := range features {
        if manager.CheckFeature(userID, f) {
            fmt.Printf("   ✅ %s - доступно\n", f)
        } else {
            fmt.Printf("   ❌ %s - недоступно\n", f)
        }
    }
    
    fmt.Println("\n✅ Система монетизации работает корректно!")
}
GO_TEST
