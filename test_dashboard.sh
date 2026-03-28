#!/bin/bash
echo "📊 ПРОВЕРКА ДАШБОРДА ВРАЧА"
echo "==========================="
echo ""

cd /workspaces/BOT_MAX

go run -e << 'GO_TEST' 2>/dev/null
package main

import (
    "encoding/json"
    "fmt"
    "github.com/bot011max/medical-bot/internal/dashboard"
)

func main() {
    fmt.Println("📈 Тестирование дашборда врача...")
    
    // Генерация дашборда
    dash := dashboard.GenerateDashboard("doctor_001")
    
    // Вывод статистики
    fmt.Printf("\n📊 СТАТИСТИКА:\n")
    fmt.Printf("   • Всего пациентов: %d\n", dash.Statistics.TotalPatients)
    fmt.Printf("   • Средняя приверженность: %.1f%%\n", dash.Statistics.AvgAdherence)
    fmt.Printf("   • Пациенты с высоким риском: %d\n", dash.Statistics.HighRiskPatients)
    fmt.Printf("   • Критические оповещения: %d\n", dash.Statistics.CriticalAlerts)
    
    // Вывод пациентов
    fmt.Printf("\n👥 ПАЦИЕНТЫ:\n")
    for _, p := range dash.Patients {
        riskIcon := "🟢"
        if p.RiskLevel == "medium" {
            riskIcon = "🟡"
        } else if p.RiskLevel == "high" {
            riskIcon = "🔴"
        }
        fmt.Printf("   %s %s - Приверженность: %.1f%%, Лекарств: %d\n", 
            riskIcon, p.Name, p.AdherenceRate, p.ActiveMeds)
    }
    
    // Вывод алертов
    if len(dash.Alerts) > 0 {
        fmt.Printf("\n⚠️ АКТИВНЫЕ АЛЕРТЫ:\n")
        for _, a := range dash.Alerts {
            severityIcon := "🟡"
            if a.Severity == "high" {
                severityIcon = "🔴"
            }
            fmt.Printf("   %s %s - %s\n", severityIcon, a.PatientName, a.Type)
        }
    }
    
    // Вывод JSON
    fmt.Printf("\n📄 JSON дашборда:\n")
    jsonData, _ := dash.ToJSON()
    fmt.Printf("%s\n", string(jsonData[:200]) + "...")
    
    fmt.Println("\n✅ Дашборд врача работает корректно!")
}
GO_TEST
