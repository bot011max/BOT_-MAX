package middleware

import (
    "github.com/gin-gonic/gin"
)

// SecurityHeaders добавляет заголовки безопасности
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        // HSTS - принудительное использование HTTPS
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
        
        // Защита от clickjacking
        c.Header("X-Frame-Options", "DENY")
        
        // Защита от MIME sniffing
        c.Header("X-Content-Type-Options", "nosniff")
        
        // Защита от XSS
        c.Header("X-XSS-Protection", "1; mode=block")
        
        // Content Security Policy
        c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:;")
        
        // Referrer Policy
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // Permissions Policy - ограничение доступа к API
        c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=()")
        
        c.Next()
    }
}
