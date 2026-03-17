package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/bot011max/BOT_MAX/internal/handlers"
    "github.com/bot011max/BOT_MAX/internal/security"
    "github.com/bot011max/BOT_MAX/internal/monitoring"
    "github.com/bot011max/BOT_MAX/internal/middleware"
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/joho/godotenv"
)

func main() {
    // Загружаем .env
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️ .env файл не найден, использую переменные окружения")
    }

    // 1. Инициализация логгера безопасности
    if err := security.InitAuditLogger(); err != nil {
        log.Printf("⚠️ Ошибка инициализации аудита: %v", err)
    }

    // 2. Инициализация метрик
    monitoring.InitMetrics()

    // 3. Проверяем Redis (но не останавливаем приложение если его нет)
    redisAddr := os.Getenv("REDIS_HOST")
    if redisAddr == "" {
        redisAddr = "redis"  // default for docker-compose
    }
    redisAddr += ":6379"
    
    log.Printf("🔄 Подключение к Redis: %s", redisAddr)

    // 4. Создаем rate limiter с правильной конфигурацией
    rateLimiter, err := security.NewAdvancedRateLimiter(redisAddr, 
        &security.AdvancedConfig{
            RequestsPerSecond: 2,        // 2 запроса в секунду
            Burst:             5,        // максимум 5 подряд
            BlockDuration:     15 * time.Minute,
            CleanupInterval:   time.Minute,
            EnableBehavioral:  true,
            TorBlock:          false,
            VPNBlock:          false,
            DatacenterBlock:   false,
        })
    
    if err != nil {
        log.Printf("❌ Ошибка создания rate limiter: %v", err)
        log.Printf("⚠️ Продолжаем без rate limiter...")
        rateLimiter = nil
    } else {
        log.Println("✅ Rate limiter создан")
    }

    // 5. Создаем WAF
    waf, err := security.NewWAFMiddleware(security.WAFConfig{
        EnableSQLInjection:   true,
        EnableXSS:           true,
        EnablePathTraversal: true,
        EnableCommandInjection: true,
        EnableScannerDetection: true,
        BlockThreshold:      10,
    })
    
    if err != nil {
        log.Printf("⚠️ Ошибка создания WAF: %v", err)
        waf = nil
    }

    // 6. Настройка Gin
    gin.SetMode(gin.ReleaseMode)
    r := gin.New()
    r.Use(gin.Recovery())
    r.Use(gin.Logger())

    // 7. Подключаем middleware если они созданы
    if waf != nil {
        r.Use(waf.Handler())
        log.Println("✅ WAF подключен")
    }
    
    if rateLimiter != nil {
        r.Use(rateLimiter.Middleware())
        log.Println("✅ Rate limiter подключен")
    }
    
    r.Use(monitoring.MetricsMiddleware())

    // 8. Метрики Prometheus
    r.GET("/metrics", gin.WrapH(promhttp.Handler()))

    // 9. Тестовый endpoint
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "time":   time.Now().Unix(),
        })
    })

    r.GET("/test-rate", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Rate limiter test endpoint",
            "ip":      c.ClientIP(),
        })
    })

    // 10. API маршруты
    api := r.Group("/api")
    {
        // Для логина используем специальный middleware если он есть
        loginHandler := handlers.Login
        if rateLimiter != nil {
            api.POST("/login", rateLimiter.LoginMiddleware(), loginHandler)
        } else {
            api.POST("/login", loginHandler)
        }
        
        api.POST("/register", handlers.Register)
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("✅ Сервер запущен на порту %s", port)
    log.Printf("🌐 http://localhost:%s", port)
    log.Printf("📊 Метрики: http://localhost:%s/metrics", port)
    
    if err := r.Run(":" + port); err != nil {
        log.Fatalf("❌ Ошибка запуска сервера: %v", err)
    }
}
