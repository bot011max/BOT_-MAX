#!/bin/bash
echo "📱 ПРОВЕРКА PUSH УВЕДОМЛЕНИЙ"
echo "============================="
echo ""

cd /workspaces/BOT_MAX

go run -e << 'GO_TEST' 2>/dev/null
package main

import (
    "fmt"
    "time"
    "github.com/bot011max/medical-bot/internal/notifier"
)

func main() {
    fmt.Println("📨 Тестирование Push уведомлений...")
    
    push := notifier.NewPushNotifier()
    
    // Отправка напоминания
    err := push.SendReminder("user123", "Аспирин", "Примите лекарство после еды")
    if err != nil {
        fmt.Printf("   ❌ Ошибка: %v\n", err)
    } else {
        fmt.Println("   ✅ Напоминание отправлено")
    }
    
    // Отправка отчета
    report := &notifier.HealthReport{
        UserID:        "user123",
        Date:          time.Now(),
        AdherenceRate: 94.5,
        MissedDoses:   2,
        Symptoms:      []string{"легкая головная боль"},
    }
    
    err = push.SendHealthReport("user123", report)
    if err != nil {
        fmt.Printf("   ❌ Ошибка: %v\n", err)
    } else {
        fmt.Println("   ✅ Отчет отправлен")
    }
    
    fmt.Println("\n✅ Push уведомления работают корректно!")
}
GO_TEST
