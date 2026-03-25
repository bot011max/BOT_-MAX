package main

import (
    "log"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/bot011max/medical-bot/internal/database"
    "github.com/bot011max/medical-bot/internal/repository"
    "github.com/bot011max/medical-bot/internal/telegram"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️ .env file not found, using environment variables")
    }

    token := os.Getenv("TELEGRAM_BOT_TOKEN")
    if token == "" {
        log.Println("⚠️ TELEGRAM_BOT_TOKEN not set, bot will not start")
    }

    if err := database.Connect(); err != nil {
        log.Println("⚠️ Database connection failed:", err)
    } else {
        defer database.Close()
    }

    userRepo := repository.NewUserRepository(database.DB)
    medicationRepo := repository.NewMedicationRepository(database.DB)

    // Запускаем Telegram бота
    if token != "" {
        bot, err := telegram.NewBot(token, userRepo, medicationRepo)
        if err != nil {
            log.Printf("⚠️ Failed to create Telegram bot: %v", err)
        } else {
            go bot.Start()
            log.Println("✅ Telegram bot started")
        }
    }

    // Запускаем HTTP сервер для вебхуков
    r := gin.Default()

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "time":   time.Now().Unix(),
        })
    })

    r.POST("/webhook/telegram", func(c *gin.Context) {
        c.JSON(200, gin.H{"ok": true})
    })

    port := os.Getenv("TELEGRAM_PORT")
    if port == "" {
        port = "8081"
    }

    log.Printf("🤖 Telegram bot webhook server starting on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatal("❌ Failed to start server:", err)
    }
}
