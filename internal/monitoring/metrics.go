// internal/monitoring/metrics.go
package monitoring

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // HTTP метрики
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

    // Бизнес метрики
    ActiveUsers = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_users_total",
            Help: "Total number of active users",
        },
    )

    PrescriptionsCreated = promauto.NewCounter(
        prometheus.CounterOpts{
            Name: "prescriptions_created_total",
            Help: "Total number of prescriptions created",
        },
    )

    // Метрики безопасности
    FailedLogins = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "failed_logins_total",
            Help: "Total number of failed login attempts",
        },
        []string{"reason"},
    )

    WAFBlocks = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "waf_blocks_total",
            Help: "Total number of requests blocked by WAF",
        },
        []string{"rule"},
    )

    RateLimiterBlocks = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "rate_limiter_blocks_total",
            Help: "Total number of requests blocked by rate limiter",
        },
        []string{"ip"},
    )

    // Метрики базы данных
    DatabaseQueriesTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "database_queries_total",
            Help: "Total number of database queries",
        },
        []string{"operation", "table"},
    )

    DatabaseQueryDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "database_query_duration_seconds",
            Help:    "Database query duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"operation", "table"},
    )

    // Метрики очередей
    QueueSize = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "queue_size",
            Help: "Current size of the queue",
        },
        []string{"queue_name"},
    )

    // Метрики Telegram бота
    TelegramMessagesTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "telegram_messages_total",
            Help: "Total number of Telegram messages",
        },
        []string{"type", "status"},
    )

    TelegramActiveChats = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "telegram_active_chats_total",
            Help: "Total number of active Telegram chats",
        },
    )
)

// Middleware для сбора метрик
func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start).Seconds()
        status := c.Writer.Status()
        
        HTTPRequestsTotal.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            http.StatusText(status),
        ).Inc()
        
        HTTPRequestDuration.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
        ).Observe(duration)
    }
}
