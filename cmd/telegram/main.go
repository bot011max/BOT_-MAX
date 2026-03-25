package main

import (
    "log"
    "os"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
    
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok", "service": "telegram-bot"})
    })
    
    router.POST("/webhook/telegram", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    port := os.Getenv("TELEGRAM_PORT")
    if port == "" {
        port = "8081"
    }
    
    log.Printf("🤖 Telegram bot webhook server starting on port %s", port)
    router.Run(":" + port)
}
