package middleware

import (
    "crypto/rand"
    "encoding/base64"
    "net/http"
    "sync"
    "time"
    
    "github.com/gin-gonic/gin"
)

type CSRFManager struct {
    tokens     map[string]csrfToken
    mu         sync.RWMutex
    expiration time.Duration
}

type csrfToken struct {
    token     string
    userID    string
    expiresAt time.Time
}

var globalCSRF = &CSRFManager{
    tokens:     make(map[string]csrfToken),
    expiration: 30 * time.Minute,
}

// GenerateCSRFToken генерирует новый CSRF токен
func GenerateCSRFToken(userID string) string {
    bytes := make([]byte, 32)
    rand.Read(bytes)
    token := base64.URLEncoding.EncodeToString(bytes)
    
    globalCSRF.mu.Lock()
    defer globalCSRF.mu.Unlock()
    
    globalCSRF.tokens[token] = csrfToken{
        token:     token,
        userID:    userID,
        expiresAt: time.Now().Add(globalCSRF.expiration),
    }
    
    return token
}

// ValidateCSRFToken проверяет CSRF токен
func ValidateCSRFToken(token, userID string) bool {
    globalCSRF.mu.RLock()
    defer globalCSRF.mu.RUnlock()
    
    csrf, exists := globalCSRF.tokens[token]
    if !exists {
        return false
    }
    
    if time.Now().After(csrf.expiresAt) {
        globalCSRF.mu.RUnlock()
        globalCSRF.mu.Lock()
        delete(globalCSRF.tokens, token)
        globalCSRF.mu.Unlock()
        globalCSRF.mu.RLock()
        return false
    }
    
    return csrf.userID == userID
}

// CSRFProtection middleware для защиты от CSRF атак
func CSRFProtection() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Только для методов, изменяющих состояние
        if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
            c.Next()
            return
        }
        
        // Для неавторизованных запросов пропускаем
        userID, exists := c.Get("user_id")
        if !exists {
            c.Next()
            return
        }
        
        token := c.GetHeader("X-CSRF-Token")
        if token == "" {
            token = c.PostForm("csrf_token")
        }
        
        if token == "" || !ValidateCSRFToken(token, userID.(string)) {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Invalid CSRF token",
                "code":  "CSRF_TOKEN_INVALID",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// GetCSRFTokenHandler возвращает CSRF токен для клиента
func GetCSRFTokenHandler(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    token := GenerateCSRFToken(userID.(string))
    c.JSON(http.StatusOK, gin.H{
        "csrf_token": token,
        "expires_in": globalCSRF.expiration.Seconds(),
    })
}
