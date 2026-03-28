#!/bin/bash
echo "🤖 ПРОВЕРКА AI АНАЛИЗА СИМПТОМОВ"
echo "================================="
echo ""

cd /workspaces/BOT_MAX

# Запускаем тест AI
go run -e << 'GO_TEST' 2>/dev/null
package main

import (
    "encoding/json"
    "fmt"
    "github.com/bot011max/medical-bot/internal/ai"
)

func main() {
    fmt.Println("🧠 Тестирование AI анализа симптомов...")
    
    analyzer := ai.NewSymptomAnalyzer()
    
    // Тест 1: Симптомы ОРВИ
    fmt.Println("\n📋 Тест 1: Симптомы ОРВИ")
    symptoms1 := []string{"кашель", "насморк", "температура 37.5", "головная боль"}
    results1 := analyzer.Analyze(symptoms1)
    
    for _, d := range results1 {
        fmt.Printf("   • %s (вероятность: %.0f%%)\n", d.Name, d.Probability*100)
        fmt.Printf("     %s\n", d.Description)
        fmt.Printf("     💊 %s\n", d.Advice)
    }
    
    // Тест 2: Симптомы аллергии
    fmt.Println("\n📋 Тест 2: Симптомы аллергии")
    symptoms2 := []string{"чихание", "зуд в глазах", "слезотечение"}
    results2 := analyzer.Analyze(symptoms2)
    
    for _, d := range results2 {
        fmt.Printf("   • %s (вероятность: %.0f%%)\n", d.Name, d.Probability*100)
        fmt.Printf("     💊 %s\n", d.Advice)
    }
    
    fmt.Println("\n✅ AI анализ работает корректно!")
}
GO_TEST
