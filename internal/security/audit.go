// internal/security/audit.go
package security

import (
    "bytes"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/google/uuid"
    "github.com/segmentio/kafka-go"
)

type AuditLogger struct {
    kafkaWriter *kafka.Writer
    httpClient  *http.Client
    siemURL     string
    enabled     bool
    environment string
}

type AuditEvent struct {
    ID          string                 `json:"id"`
    Timestamp   time.Time              `json:"timestamp"`
    Type        string                 `json:"type"`
    UserID      string                 `json:"user_id,omitempty"`
    SessionID   string                 `json:"session_id,omitempty"`
    Action      string                 `json:"action"`
    Resource    string                 `json:"resource"`
    Status      string                 `json:"status"`
    IP          string                 `json:"ip,omitempty"`
    UserAgent   string                 `json:"user_agent,omitempty"`
    Details     map[string]interface{} `json:"details,omitempty"`
    Severity    string                 `json:"severity"`
    Hash        string                 `json:"hash"` // Для целостности логов
}

func NewAuditLogger() (*AuditLogger, error) {
    // Для production - отправка в Kafka/SIEM
    kafkaWriter := &kafka.Writer{
        Addr:     kafka.TCP(os.Getenv("KAFKA_BROKERS")),
        Topic:    "audit-logs",
        Balancer: &kafka.LeastBytes{},
    }

    return &AuditLogger{
        kafkaWriter: kafkaWriter,
        httpClient:  &http.Client{Timeout: 5 * time.Second},
        siemURL:     os.Getenv("SIEM_URL"),
        enabled:     true,
        environment: os.Getenv("ENVIRONMENT"),
    }, nil
}

func (al *AuditLogger) Log(eventType, action, resource, status, userID string, details map[string]interface{}) {
    if !al.enabled {
        return
    }

    event := &AuditEvent{
        ID:        uuid.New().String(),
        Timestamp: time.Now().UTC(),
        Type:      eventType,
        UserID:    userID,
        Action:    action,
        Resource:  resource,
        Status:    status,
        Details:   details,
        Severity:  al.determineSeverity(eventType, status),
    }

    // Добавляем хеш для целостности
    event.Hash = al.calculateHash(event)

    // Отправляем в несколько каналов
    go al.sendToKafka(event)
    go al.sendToSIEM(event)
    go al.logToFile(event)
}

func (al *AuditLogger) determineSeverity(eventType, status string) string {
    if status == "ERROR" || status == "FAILED" {
        return "HIGH"
    }
    if strings.Contains(eventType, "SECURITY_") {
        return "CRITICAL"
    }
    if status == "WARNING" {
        return "MEDIUM"
    }
    return "LOW"
}

func (al *AuditLogger) calculateHash(event *AuditEvent) string {
    data := fmt.Sprintf("%s:%s:%s:%s:%s", 
        event.ID, event.Timestamp, event.UserID, event.Action, event.Status)
    
    h := hmac.New(sha256.New, []byte(os.Getenv("AUDIT_HMAC_KEY")))
    h.Write([]byte(data))
    return hex.EncodeToString(h.Sum(nil))
}

func (al *AuditLogger) sendToKafka(event *AuditEvent) {
    if al.kafkaWriter == nil {
        return
    }

    data, _ := json.Marshal(event)
    msg := kafka.Message{
        Key:   []byte(event.ID),
        Value: data,
        Headers: []kafka.Header{
            {Key: "type", Value: []byte(event.Type)},
            {Key: "severity", Value: []byte(event.Severity)},
        },
    }

    if err := al.kafkaWriter.WriteMessages(context.Background(), msg); err != nil {
        fmt.Printf("Failed to send to Kafka: %v\n", err)
    }
}

func (al *AuditLogger) sendToSIEM(event *AuditEvent) {
    if al.siemURL == "" {
        return
    }

    data, _ := json.Marshal(event)
    req, _ := http.NewRequest("POST", al.siemURL+"/api/events", bytes.NewBuffer(data))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-API-Key", os.Getenv("SIEM_API_KEY"))

    resp, err := al.httpClient.Do(req)
    if err != nil {
        fmt.Printf("Failed to send to SIEM: %v\n", err)
        return
    }
    defer resp.Body.Close()
}

func (al *AuditLogger) logToFile(event *AuditEvent) {
    filename := fmt.Sprintf("audit-%s.log", time.Now().Format("2006-01-02"))
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return
    }
    defer f.Close()

    data, _ := json.Marshal(event)
    f.Write(append(data, '\n'))
}

// Глобальный экземпляр
var GlobalAuditLogger *AuditLogger

func InitAuditLogger() error {
    logger, err := NewAuditLogger()
    if err != nil {
        return err
    }
    GlobalAuditLogger = logger
    return nil
}

func AuditLog(action, userID string, details map[string]interface{}) {
    if GlobalAuditLogger != nil {
        GlobalAuditLogger.Log("AUDIT", action, "system", "SUCCESS", userID, details)
    }
}

func SecurityAlert(alertType string, details map[string]interface{}) {
    if GlobalAuditLogger != nil {
        GlobalAuditLogger.Log("SECURITY_ALERT", alertType, "security", "WARNING", "system", details)
    }
}
