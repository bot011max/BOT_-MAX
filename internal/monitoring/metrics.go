package monitoring

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    HTTPRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )

    RateLimiterBlocks = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "rate_limiter_blocks_total",
            Help: "Total number of requests blocked by rate limiter",
        },
        []string{"ip"},
    )

    LoginAttempts = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "login_attempts_total",
            Help: "Total number of login attempts",
        },
        []string{"status", "ip"},
    )
)

// MetricsMiddleware собирает метрики для HTTP запросов
func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        duration := time.Since(start).Seconds()
        
        status := c.Writer.Status()
        method := c.Request.Method
        endpoint := c.FullPath()
        
        HTTPRequestsTotal.WithLabelValues(method, endpoint, string(rune(status))).Inc()
        HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
    }
}

// MetricsHandler возвращает обработчик для метрик Prometheus
// Note: promhttp.Handler должен быть добавлен в main.go
