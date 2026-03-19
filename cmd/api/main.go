package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/bot011max/medical-bot/internal/api"
    "github.com/bot011max/medical-bot/internal/security"
    "github.com/bot011max/medical-bot/internal/middleware"
)

func main() {
    // 1. Инициализация безопасности (ВОЕННЫЙ УРОВЕНЬ)
    armor := security.NewAbsoluteArmor()
    
    // 2. Настройка роутера
    r := gin.New()
    r.Use(middleware.SecurityHeaders())
    r.Use(armor.ProtectRequest()) // WAF + Rate Limiting + IDS
    
    // 3. Маршруты
    api.SetupRoutes(r)
    
    // 4. Запуск с graceful shutdown
    srv := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }
    
    go func() {
        if err := srv.ListenAndServe(); err != nil {
            log.Printf("Server error: %v", err)
        }
    }()
    
    // Ожидание сигнала завершения
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    srv.Shutdown(ctx)
}
