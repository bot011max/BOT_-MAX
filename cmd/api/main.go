package main

import (
    "log"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/bot011max/medical-bot/internal/api"
    "github.com/bot011max/medical-bot/internal/audit"
    "github.com/bot011max/medical-bot/internal/database"
    "github.com/bot011max/medical-bot/internal/middleware"
    "github.com/bot011max/medical-bot/internal/models"
    "github.com/bot011max/medical-bot/internal/repository"
    "github.com/bot011max/medical-bot/internal/security"
    "github.com/bot011max/medical-bot/internal/ocr"
    "github.com/bot011max/medical-bot/internal/security"
    "github.com/bot011max/medical-bot/internal/service"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️ .env file not found")
    }

    // Инициализация безопасности
    quantumCrypto := security.NewQuantumCrypto()
    log.Printf("🔐 Quantum key: %x", quantumCrypto.GetQuantumKeyPreview(5))
    
    intrusionDetection := security.NewIntrusionDetectionSystem()
    intrusionPrevention := security.NewIPS()
    adaptiveLimiter := security.NewAdaptiveRateLimiter()
    blockchainAudit := audit.NewBlockchain()
    
    blockchainAudit.AddEvent("SYSTEM_START", "system", "initialize", "Medical bot starting")

    // Подключение к БД
    if err := database.Connect(); err != nil {
        log.Println("⚠️ Database connection failed:", err)
    } else {
        defer database.Close()
        database.DB.AutoMigrate(&models.User{}, &models.Medication{}, &models.Reminder{})
    }

    // Репозитории и сервисы
    userRepo := repository.NewUserRepository(database.DB)
    medicationRepo := repository.NewMedicationRepository(database.DB)
    authService := service.NewAuthService(userRepo)
    authHandler := api.NewAuthHandler(authService)
    medicationHandler := api.NewMedicationHandler(medicationRepo)
    securityHandler := api.NewSecurityHandler()

    // Настройка Gin
    r := gin.Default()

    // Middleware защиты
    r.Use(func(c *gin.Context) {
        ip := c.ClientIP()
        if intrusionDetection.DetectPortScan(ip) {
            intrusionPrevention.BlockIP(ip, 30*time.Minute)
            c.AbortWithStatusJSON(429, gin.H{"error": "IP blocked"})
            return
        }
        if !adaptiveLimiter.Allow(ip) {
            c.AbortWithStatusJSON(429, gin.H{"error": "Rate limit exceeded"})
            return
        }
        c.Next()
    })

    // Маршруты
    r.POST("/api/register", authHandler.Register)
    r.POST("/api/login", authHandler.Login)
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok", "time": time.Now().Unix()})
    })
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Medical Bot API", "version": "2.0.0"})
    })

    // Защищенные маршруты
    authorized := r.Group("/api")
    authorized.Use(middleware.AuthMiddleware(authService))
    {
        authorized.GET("/profile", authHandler.Profile)
        authorized.POST("/medications", medicationHandler.Create)
        authorized.GET("/medications", medicationHandler.List)
        authorized.GET("/medications/:id", medicationHandler.Get)
        authorized.PUT("/medications/:id", medicationHandler.Update)
        authorized.DELETE("/medications/:id", medicationHandler.Delete)
    }

    // Метрики безопасности
    r.GET("/security/status", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "quantum":      true,
            "ids":          true,
            "ips":          true,
            "rate_limiter": true,
            "audit":        true,
            "blocks":       len(blockchainAudit.GetEvents()),
        })
    })

    // Запуск
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("🚀 Server starting on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal("❌ Failed to start server:", err)
    }
}
