package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Reminder struct {
    ID             uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID         uuid.UUID      `json:"user_id" gorm:"type:uuid;index;not null"`
    MedicationID   *uuid.UUID     `json:"medication_id,omitempty" gorm:"type:uuid;index"`
    ScheduledAt    time.Time      `json:"scheduled_at" gorm:"index"`
    Message        string         `json:"message" gorm:"type:text"`
    Status         string         `json:"status" gorm:"default:'pending';size:20"`
    SentAt         *time.Time     `json:"sent_at"`
    AcknowledgedAt *time.Time     `json:"acknowledged_at"`
    RetryCount     int            `json:"retry_count" gorm:"default:0"`
    ErrorMessage   string         `json:"error_message" gorm:"type:text"`
    CreatedAt      time.Time      `json:"created_at"`
    UpdatedAt      time.Time      `json:"updated_at"`
    DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
    
    User           User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
    Medication     *Medication    `json:"medication,omitempty" gorm:"foreignKey:MedicationID"`
}

func (r *Reminder) BeforeCreate(tx *gorm.DB) error {
    if r.ID == uuid.Nil {
        r.ID = uuid.New()
    }
    return nil
}

func (r *Reminder) CanRetry() bool {
    return r.Status == "failed" && r.RetryCount < 3
}
