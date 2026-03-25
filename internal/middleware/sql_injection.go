package middleware

import (
    "regexp"
    "github.com/gin-gonic/gin"
)

var sqlPatterns = []*regexp.Regexp{
    regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE|DROP|UNION|--|;|'|\|)`),
    regexp.MustCompile(`(?i)(OR|AND)\s+\d+=\d+`),
    regexp.MustCompile(`(?i)WAITFOR\s+DELAY`),
    regexp.MustCompile(`(?i)EXEC\s+XP_`),
}

func SQLInjectionProtection() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Проверка параметров URL
        for _, values := range c.Request.URL.Query() {
            for _, value := range values {
                if isSQLInjection(value) {
                    c.JSON(400, gin.H{"error": "Invalid request parameters"})
                    c.Abort()
                    return
                }
            }
        }
        
        // Проверка заголовков
        for _, values := range c.Request.Header {
            for _, value := range values {
                if isSQLInjection(value) {
                    c.JSON(400, gin.H{"error": "Invalid request headers"})
                    c.Abort()
                    return
                }
            }
        }
        
        c.Next()
    }
}

func isSQLInjection(input string) bool {
    for _, pattern := range sqlPatterns {
        if pattern.MatchString(input) {
            return true
        }
    }
    return false
}
