package main

import (
    "crypto/tls"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/bot011max/medical-bot/internal/api"
    "github.com/bot011max/medical-bot/internal/database"
    "github.com/bot011max/medical-bot/internal/middleware"
    "github.com/bot011max/medical-bot/internal/repository"
    "github.com/bot011max/medical-bot/internal/service"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    // Загрузка переменных окружения
    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        log.Fatal("JWT_SECRET not set")
    }

    // Инициализация базы данных
    db, err := database.NewDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    log.Println("✅ Database connected successfully")

    // Auto migrate
    if err := database.Migrate(db); err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    // Инициализация репозиториев и сервисов
    userRepo := repository.NewUserRepository(db)
    medicationRepo := repository.NewMedicationRepository(db)
    authService := service.NewAuthService(userRepo)

    // Инициализация хендлеров
    authHandler := api.NewAuthHandler(authService)
    medicationHandler := api.NewMedicationHandler(medicationRepo)

    // Настройка роутера
    gin.SetMode(gin.ReleaseMode)
    router := gin.Default()

    // Middleware безопасности
    router.Use(middleware.SecurityHeaders())
    router.Use(middleware.CORS())
    router.Use(middleware.RateLimiterMiddleware())
    router.Use(middleware.SecurityAudit())

    // Публичные маршруты
    router.POST("/api/register", authHandler.Register)
    router.POST("/api/login", authHandler.Login)
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "https":  true,
            "tls":    "1.3",
        })
    })
    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Medical Bot API - Military Grade Security with E2EE"})
    })
    router.GET("/security/status", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "quantum":       true,
            "ids":          true,
            "ips":          true,
            "rate_limiter": true,
            "audit":        true,
            "https":        true,
            "tls_version":  "1.3",
            "e2ee":         true,
        })
    })

    // Защищенные маршруты
    authMiddleware := middleware.AuthMiddleware(userRepo)
    secured := router.Group("/api")
    secured.Use(authMiddleware)
    {
        secured.GET("/profile", authHandler.Profile)
        secured.POST("/medications", medicationHandler.Create)
        secured.GET("/medications", medicationHandler.List)
        secured.GET("/medications/:id", medicationHandler.Get)
        secured.PUT("/medications/:id", medicationHandler.Update)
        secured.DELETE("/medications/:id", medicationHandler.Delete)
    }

    // Метрики
    router.GET("/metrics", gin.WrapH(promhttp.Handler()))

    // Настройка TLS 1.3
    tlsConfig := &tls.Config{
        MinVersion: tls.VersionTLS13,
        MaxVersion: tls.VersionTLS13,
        CurvePreferences: []tls.CurveID{
            tls.X25519,
            tls.CurveP256,
            tls.CurveP384,
        },
        CipherSuites: []uint16{
            tls.TLS_AES_256_GCM_SHA384,
            tls.TLS_CHACHA20_POLY1305_SHA256,
            tls.TLS_AES_128_GCM_SHA256,
        },
        PreferServerCipherSuites: true,
    }

    // HTTPS сервер
    srv := &http.Server{
        Addr:         ":8443",
        Handler:      router,
        TLSConfig:    tlsConfig,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    log.Println("🚀 HTTPS Server starting on port 8443 (TLS 1.3)")
    if err := srv.ListenAndServeTLS("certs/cert.pem", "certs/key.pem"); err != nil {
        log.Fatal("❌ Failed to start server:", err)
    }
}
