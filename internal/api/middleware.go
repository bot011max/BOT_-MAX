package api

import (
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/redis/go-redis/v9"

    "github.com/bot011max/medical-bot/internal/models"
    "github.com/bot011max/medical-bot/internal/repository"
    "github.com/bot011max/medical-bot/internal/security"
)

// AuthMiddleware - middleware для JWT аутентификации
func AuthMiddleware(jwtSecret string, audit *security.AuditLogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            audit.LogSecurity("MISSING_AUTH", "", c.ClientIP(), nil)
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "authorization header required",
            })
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "invalid authorization header format",
            })
            return
        }

        tokenStr := parts[1]

        // Парсинг и валидация JWT
        token, err := jwt.ParseWithClaims(tokenStr, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(jwtSecret), nil
        })

        if err != nil {
            audit.LogSecurity("INVALID_TOKEN", "", c.ClientIP(), map[string]interface{}{
                "error": err.Error(),
            })
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "invalid or expired token",
            })
            return
        }

        claims, ok := token.Claims.(*AuthClaims)
        if !ok || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "invalid token claims",
            })
            return
        }

        // Проверка blacklist (если токен отозван)
        // TODO: проверить в Redis

        c.Set("user_id", claims.UserID)
        c.Set("user_role", claims.Role)
        c.Set("subscription", claims.Subscription)
        c.Set("twofa_completed", claims.TwoFACompleted)
        
        c.Next()
    }
}

// AuthClaims - кастомные claims
type AuthClaims struct {
    UserID         string `json:"user_id"`
    Role           string `json:"role"`
    Subscription   string `json:"subscription"`
    TwoFACompleted bool   `json:"twofa_completed"`
    jwt.RegisteredClaims
}

// RequireRole - middleware для проверки роли
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("user_role")
        
        for _, role := range allowedRoles {
            if userRole == role {
                c.Next()
                return
            }
        }
        
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
            "error": "access denied",
            "required": allowedRoles,
        })
    }
}

// RequireSubscription - middleware для проверки подписки и лимитов
func RequireSubscription(repo *repository.SubscriptionRepository, cache *redis.Client, audit *security.AuditLogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")
        
        // Получаем из кэша (Redis)
        var sub *models.Subscription
        cacheKey := "subscription:" + userID
        
        cached, err := cache.Get(c, cacheKey).Result()
        if err == nil {
            // TODO: десериализовать из JSON
        }
        
        if sub == nil {
            // Из БД
            sub, err = repo.GetActiveByUserIDStr(userID)
            if err != nil {
                audit.LogSecurity("SUBSCRIPTION_CHECK_FAILED", userID, c.ClientIP(), nil)
                c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{
                    "error": "no active subscription",
                })
                return
            }
            
            // Кэшируем на 1 час
            // TODO: сериализовать в JSON
            cache.Set(c, cacheKey, sub, time.Hour)
        }
        
        c.Set("subscription_obj", sub)
        
        // Проверка конкретных лимитов в зависимости от эндпоинта
        // Будет вызвано в конкретных хендлерах
        
        c.Next()
    }
}

// CheckPatientLimit - проверка лимита пациентов
func CheckPatientLimit(repo *repository.PatientRepository, audit *security.AuditLogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetString("user_id")
        sub := c.MustGet("subscription_obj").(*models.Subscription)
        
        if sub.MaxPatients < 0 { // -1 означает безлимит
            c.Next()
            return
        }
        
        current, err := repo.CountByDoctorID(userID)
        if err != nil {
            audit.LogError("PATIENT_COUNT_FAILED", userID, err)
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                "error": "failed to check patient limit",
            })
            return
        }
        
        if current >= sub.MaxPatients {
            audit.LogLimit("PATIENT_LIMIT_REACHED", userID, map[string]interface{}{
                "current": current,
                "limit":   sub.MaxPatients,
            })
            c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{
                "error":      "patient limit reached",
                "current":    current,
                "limit":      sub.MaxPatients,
                "upgrade_url": "/subscription/upgrade",
            })
            return
        }
        
        c.Set("patient_count", current)
        c.Next()
    }
}

// RateLimitMiddleware - ограничение частоты запросов
func RateLimitMiddleware(redis *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        key := "rate_limit:" + ip
        
        count, err := redis.Incr(c, key).Result()
        if err != nil {
            c.Next()
            return
        }
        
        if count == 1 {
            redis.Expire(c, key, time.Minute)
        }
        
        if count > 60 { // 60 запросов в минуту
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error": "rate limit exceeded",
                "retry_after": 60,
            })
            return
        }
        
        c.Next()
    }
}

// AuditMiddleware - логирование всех действий
func AuditMiddleware(audit *security.AuditLogger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // До обработки
        c.Next()
        // После обработки
        audit.LogAccess(
            c.Request.Method+" "+c.FullPath(),
            c.GetString("user_id"),
            c.ClientIP(),
            map[string]interface{}{
                "status": c.Writer.Status(),
                "path":   c.Request.URL.Path,
                "method": c.Request.Method,
                "user_agent": c.Request.UserAgent(),
            },
        )
    }
}
