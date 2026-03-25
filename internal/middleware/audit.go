package middleware

import (
    "log"
    "time"
    "github.com/gin-gonic/gin"
)

func SecurityAudit() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        // Логируем запрос
        log.Printf("[AUDIT] %s %s - IP: %s, User-Agent: %s",
            c.Request.Method,
            c.Request.URL.Path,
            c.ClientIP(),
            c.GetHeader("User-Agent"))
        
        c.Next()
        
        // Логируем ответ
        duration := time.Since(start)
        log.Printf("[AUDIT] Response: %d - Duration: %v",
            c.Writer.Status(),
            duration)
        
        // Логируем ошибки аутентификации
        if c.Writer.Status() == 401 || c.Writer.Status() == 403 {
            log.Printf("[SECURITY] Authentication failure - IP: %s, Path: %s",
                c.ClientIP(),
                c.Request.URL.Path)
        }
    }
}
