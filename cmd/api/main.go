package main

import (
    "log"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // Загрузка .env
    godotenv.Load()

    r := gin.Default()

    // Health check
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

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("API server starting on port %s", port)
    r.Run(":" + port)
}
