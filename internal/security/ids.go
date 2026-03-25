package security

import (
    "sync"
    "time"
)

type IntrusionDetectionSystem struct {
    requestHistory map[string][]time.Time
    mu             sync.RWMutex
}

func NewIntrusionDetectionSystem() *IntrusionDetectionSystem {
    return &IntrusionDetectionSystem{
        requestHistory: make(map[string][]time.Time),
    }
}

// DetectPortScan - обнаружение сканирования портов
func (ids *IntrusionDetectionSystem) DetectPortScan(ip string) bool {
    ids.mu.Lock()
    defer ids.mu.Unlock()
    
    // Симуляция: если IP делает много запросов - блокируем
    if len(ids.requestHistory[ip]) > 10 {
        return true
    }
    return false
}

// DetectBruteForce - обнаружение брутфорса
func (ids *IntrusionDetectionSystem) DetectBruteForce(email string) bool {
    ids.mu.Lock()
    defer ids.mu.Unlock()
    
    key := "login:" + email
    if len(ids.requestHistory[key]) > 5 {
        return true
    }
    return false
}

// RecordRequest - запись запроса для анализа
func (ids *IntrusionDetectionSystem) RecordRequest(ip string, endpoint string) {
    ids.mu.Lock()
    defer ids.mu.Unlock()
    
    ids.requestHistory[ip] = append(ids.requestHistory[ip], time.Now())
    
    // Очистка старых записей (старше 5 минут)
    cutoff := time.Now().Add(-5 * time.Minute)
    var recent []time.Time
    for _, t := range ids.requestHistory[ip] {
        if t.After(cutoff) {
            recent = append(recent, t)
        }
    }
    ids.requestHistory[ip] = recent
}

type IPS struct {
    blockedIPs   map[string]time.Time
    blockedUsers map[string]time.Time
    mu           sync.RWMutex
}

func NewIPS() *IPS {
    return &IPS{
        blockedIPs:   make(map[string]time.Time),
        blockedUsers: make(map[string]time.Time),
    }
}

func (ips *IPS) BlockIP(ip string, duration time.Duration) {
    ips.mu.Lock()
    defer ips.mu.Unlock()
    ips.blockedIPs[ip] = time.Now().Add(duration)
    go func() {
        time.Sleep(duration)
        ips.mu.Lock()
        delete(ips.blockedIPs, ip)
        ips.mu.Unlock()
    }()
}

func (ips *IPS) IsIPBlocked(ip string) bool {
    ips.mu.RLock()
    defer ips.mu.RUnlock()
    blockedUntil, exists := ips.blockedIPs[ip]
    return exists && time.Now().Before(blockedUntil)
}
