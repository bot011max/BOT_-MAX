package middleware

import (
    "sync"
    "time"
    "github.com/gin-gonic/gin"
)

type RateLimiter struct {
    attempts map[string][]time.Time
    blocked  map[string]time.Time
    mu       sync.RWMutex
    limit    int
    window   time.Duration
    block    time.Duration
}

func NewRateLimiter() *RateLimiter {
    return &RateLimiter{
        attempts: make(map[string][]time.Time),
        blocked:  make(map[string]time.Time),
        limit:    5,
        window:   time.Minute,
        block:    15 * time.Minute,
    }
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Пропускаем только для login и register
        if c.Request.URL.Path != "/api/login" && c.Request.URL.Path != "/api/register" {
            c.Next()
            return
        }
        
        ip := c.ClientIP()
        
        // Проверка блокировки
        rl.mu.RLock()
        if blockUntil, exists := rl.blocked[ip]; exists {
            if time.Now().Before(blockUntil) {
                rl.mu.RUnlock()
                c.JSON(429, gin.H{
                    "error":       "Too many attempts. Please try again later.",
                    "retry_after": int(blockUntil.Sub(time.Now()).Seconds()),
                })
                c.Abort()
                return
            }
            rl.mu.RUnlock()
            rl.mu.Lock()
            delete(rl.blocked, ip)
            rl.mu.Unlock()
        } else {
            rl.mu.RUnlock()
        }
        
        // Проверка лимита
        rl.mu.Lock()
        now := time.Now()
        attempts := rl.attempts[ip]
        
        // Очистка старых попыток
        valid := make([]time.Time, 0)
        for _, t := range attempts {
            if now.Sub(t) <= rl.window {
                valid = append(valid, t)
            }
        }
        
        if len(valid) >= rl.limit {
            rl.blocked[ip] = now.Add(rl.block)
            rl.mu.Unlock()
            c.JSON(429, gin.H{
                "error": "Rate limit exceeded. IP blocked for 15 minutes.",
            })
            c.Abort()
            return
        }
        
        rl.attempts[ip] = append(valid, now)
        rl.mu.Unlock()
        
        c.Next()
    }
}

var globalRateLimiter = NewRateLimiter()

func RateLimiterMiddleware() gin.HandlerFunc {
    return globalRateLimiter.Middleware()
}
