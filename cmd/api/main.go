package main

import (
    "log"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/bot011max/medical-bot/internal/api"
    "github.com/bot011max/medical-bot/internal/database"
    "github.com/bot011max/medical-bot/internal/middleware"
    "github.com/bot011max/medical-bot/internal/models"
    "github.com/bot011max/medical-bot/internal/repository"
    "github.com/bot011max/medical-bot/internal/service"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️ .env file not found, using environment variables")
    }

    if err := database.Connect(); err != nil {
        log.Println("⚠️ Database connection failed:", err)
        log.Println("📝 Continuing without database...")
    } else {
        log.Println("✅ Database connected successfully")
        defer database.Close()

        if err := database.DB.AutoMigrate(
            &models.User{},
            &models.Medication{},
            &models.Reminder{},
        ); err != nil {
            log.Println("⚠️ Migration warning:", err)
        }
    }

    userRepo := repository.NewUserRepository(database.DB)
    medicationRepo := repository.NewMedicationRepository(database.DB)

    authService := service.NewAuthService(userRepo)
    authHandler := api.NewAuthHandler(authService)
    medicationHandler := api.NewMedicationHandler(medicationRepo)

    r := gin.Default()

    // Публичные маршруты
    r.POST("/api/register", authHandler.Register)
    r.POST("/api/login", authHandler.Login)

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "time":   time.Now().Unix(),
        })
    })

    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Medical Bot API",
            "version": "1.0.0",
        })
    })

    // Защищенные маршруты
    authorized := r.Group("/api")
    authorized.Use(middleware.AuthMiddleware(authService))
    {
        authorized.GET("/profile", authHandler.Profile)
        authorized.GET("/ping", func(c *gin.Context) {
            c.JSON(200, gin.H{"message": "pong"})
        })

        authorized.POST("/medications", medicationHandler.Create)
        authorized.GET("/medications", medicationHandler.List)
        authorized.GET("/medications/:id", medicationHandler.Get)
        authorized.PUT("/medications/:id", medicationHandler.Update)
        authorized.DELETE("/medications/:id", medicationHandler.Delete)
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("🚀 API server starting on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal("❌ Failed to start server:", err)
    }
}
