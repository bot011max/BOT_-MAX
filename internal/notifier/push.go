package notifier

import (
    "log"
    "time"
)

type PushNotifier struct {
    enabled bool
}

func NewPushNotifier() *PushNotifier {
    log.Println("📱 Push notifier initialized (simulation mode)")
    return &PushNotifier{enabled: true}
}

func (n *PushNotifier) SendReminder(userID, medicationName, message string) error {
    log.Printf("📨 [PUSH] To: %s, Medication: %s, Message: %s", userID, medicationName, message)
    return nil
}

func (n *PushNotifier) SendHealthReport(userID string, report *HealthReport) error {
    log.Printf("📊 [REPORT] To: %s, Adherence: %.1f%%, Missed: %d", 
        userID, report.AdherenceRate, report.MissedDoses)
    return nil
}

type HealthReport struct {
    UserID        string    `json:"user_id"`
    Date          time.Time `json:"date"`
    AdherenceRate float64   `json:"adherence_rate"`
    MissedDoses   int       `json:"missed_doses"`
    Symptoms      []string  `json:"symptoms"`
}
