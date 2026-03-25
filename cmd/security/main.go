package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/bot011max/medical-bot/internal/api"
)

func main() {
    router := gin.Default()
    securityHandler := api.NewSecurityHandler()
    
    // Новые эндпоинты безопасности
    router.GET("/security/hsm", securityHandler.GetHSMInfo)
    router.POST("/security/backup", securityHandler.CreateBackup)
    router.GET("/security/backups", securityHandler.ListBackups)
    router.POST("/security/restore/:id", securityHandler.Rollback)
    router.POST("/api/prescription/scan", securityHandler.ProcessPrescription)
    
    log.Println("🔒 Security API server starting on port 8090")
    if err := router.Run(":8090"); err != nil {
        log.Fatal("Failed to start security server:", err)
    }
}
