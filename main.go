// cmd/api/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/bot011max/BOT_MAX/internal/security"
    "github.com/bot011max/BOT_MAX/internal/monitoring"
    "github.com/bot011max/BOT_MAX/internal/middleware"
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    // 1. Инициализация логгера безопасности
    if err := security.InitAuditLogger(); err != nil {
        log.Fatalf("Ошибка инициализации аудита: %v", err)
    }

    // 2. Инициализация метрик
    monitoring.InitMetrics()

    // 3. Создание WAF
    wafConfig := security.WAFConfig{
        EnableSQLInjection:   true,
        EnableXSS:           true,
        EnablePathTraversal: true,
        EnableCommandInjection: true,
        EnableScannerDetection: true,
        BlockThreshold:      10,
    }
    
    waf, err := security.NewWAFMiddleware(wafConfig)
    if err != nil {
        log.Fatalf("Ошибка создания WAF: %v", err)
    }

    // 4. Rate limiter с Redis
    rateLimiter, err := security.NewAdaptiveRateLimiter("redis:6379", 
        &security.RateLimiterConfig{
            RequestsPerSecond: 10,
            Burst:             20,
            BlockDuration:     time.Hour,
            CleanupInterval:   time.Minute,
            EnableAdaptive:    true,
        })
    if err != nil {
        log.Fatalf("Ошибка создания rate limiter: %v", err)
    }

    // 5. Настройка Gin
    gin.SetMode(gin.ReleaseMode)
    r := gin.New()
    
    // Middleware в правильном порядке
    r.Use(gin.Recovery())
    r.Use(middleware.SecurityHeaders())
    r.Use(middleware.RequestID())
    r.Use(monitoring.MetricsMiddleware())
    r.Use(waf.Handler())
    r.Use(rateLimiter.Middleware())
    r.Use(middleware.AuditLogMiddleware())
    
    // 6. Метрики Prometheus
    r.GET("/metrics", gin.WrapH(promhttp.Handler()))
    
    // 7. Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "healthy",
            "time":   time.Now().Unix(),
        })
    })

    // 8. API routes
    api := r.Group("/api")
    {
        api.POST("/register", handlers.Register)
        api.POST("/login", rateLimiter.LoginMiddleware(), handlers.Login)
        
        // Защищенные маршруты
        authorized := api.Group("/")
        authorized.Use(middleware.AuthRequired())
        {
            authorized.GET("/profile", handlers.Profile)
            authorized.GET("/patients", middleware.RoleRequired("doctor", "admin"), handlers.GetPatients)
        }
    }

    // 9. HTTP сервер с таймаутами
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      r,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
        MaxHeaderBytes: 1 << 20, // 1MB
    }

    // 10. Graceful shutdown
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Ошибка сервера: %v", err)
        }
    }()

    log.Println("✅ Сервер запущен на :8080")

    // Ожидание сигнала завершения
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("🔄 Завершение работы...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Ошибка при завершении: %v", err)
    }

    log.Println("✅ Сервер остановлен")
}
