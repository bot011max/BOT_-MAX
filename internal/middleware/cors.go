package middleware

import (
    "github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Разрешенные источники
        allowedOrigins := map[string]bool{
            "http://localhost:3000":   true,
            "https://localhost:3000":  true,
            "http://localhost:8080":   true,
            "https://localhost:8443":  true,
            "http://localhost:8090":   true,
            "https://localhost:8090":  true,
        }
        
        origin := c.GetHeader("Origin")
        if allowedOrigins[origin] {
            c.Header("Access-Control-Allow-Origin", origin)
        } else {
            c.Header("Access-Control-Allow-Origin", "")
        }
        
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token, X-Request-ID")
        c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type, X-Request-ID")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Max-Age", "86400")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
