package main

import (
    "log"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    godotenv.Load()

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

    port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }

    log.Printf("Telegram bot starting on port %s", port)
    r.Run(":" + port)
}
