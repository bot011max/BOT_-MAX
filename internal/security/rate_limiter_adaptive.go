package security

import (
    "sync"
    "time"
)

type AdaptiveRateLimiter struct {
    limits     map[string]*UserLimit
    mu         sync.RWMutex
    thresholds map[string]int
}

type UserLimit struct {
    Count      int
    LastReset  time.Time
    Violations int
}

func NewAdaptiveRateLimiter() *AdaptiveRateLimiter {
    return &AdaptiveRateLimiter{
        limits: make(map[string]*UserLimit),
        thresholds: map[string]int{
            "normal":    100,
            "suspicious": 50,
            "attack":     10,
        },
    }
}

func (rl *AdaptiveRateLimiter) Allow(key string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    now := time.Now()
    limit, exists := rl.limits[key]
    
    if !exists {
        rl.limits[key] = &UserLimit{
            Count:     1,
            LastReset: now,
        }
        return true
    }
    
    // Сброс каждую минуту
    if now.Sub(limit.LastReset) > time.Minute {
        limit.Count = 1
        limit.LastReset = now
        return true
    }
    
    // Динамический лимит в зависимости от нарушений
    maxAllowed := rl.thresholds["normal"]
    if limit.Violations > 3 {
        maxAllowed = rl.thresholds["suspicious"]
    }
    if limit.Violations > 10 {
        maxAllowed = rl.thresholds["attack"]
    }
    
    if limit.Count > maxAllowed {
        limit.Violations++
        return false
    }
    
    limit.Count++
    return true
}
