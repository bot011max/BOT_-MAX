package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/bot011max/medical-bot/internal/api"
)

func main() {
    router := gin.Default()
    securityHandler := api.NewSecurityHandler()
    
    router.GET("/security/hsm", securityHandler.GetHSMInfo)
    router.POST("/security/backup", securityHandler.CreateBackup)
    router.GET("/security/backups", securityHandler.ListBackups)
    router.POST("/security/restore/:id", securityHandler.Rollback)
    router.POST("/api/prescription/scan", securityHandler.ProcessPrescription)
    
    log.Println("🔒 Security API server starting on port 8090")
    router.Run(":8090")
}
