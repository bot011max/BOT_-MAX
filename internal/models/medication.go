package models

import (
    "time"
    "gorm.io/gorm"
)

type Medication struct {
    ID           string         `json:"id" gorm:"primaryKey;type:text"`
    UserID       string         `json:"user_id" gorm:"type:text;not null;index"`
    Name         string         `json:"name" gorm:"not null"`
    Dosage       string         `json:"dosage"`
    Frequency    string         `json:"frequency"`
    Instructions string         `json:"instructions"`
    StartDate    time.Time      `json:"start_date"`
    EndDate      *time.Time     `json:"end_date,omitempty"`
    IsActive     bool           `json:"is_active" gorm:"default:true"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type Reminder struct {
    ID             string         `json:"id" gorm:"primaryKey;type:text"`
    UserID         string         `json:"user_id" gorm:"type:text;not null;index"`
    MedicationID   *string        `json:"medication_id,omitempty" gorm:"type:text"`
    ScheduledAt    time.Time      `json:"scheduled_at" gorm:"index"`
    Message        string         `json:"message"`
    Status         string         `json:"status" gorm:"default:'pending'"`
    SentAt         *time.Time     `json:"sent_at,omitempty"`
    AcknowledgedAt *time.Time     `json:"acknowledged_at,omitempty"`
    RetryCount     int            `json:"retry_count" gorm:"default:0"`
    ErrorMessage   string         `json:"error_message,omitempty"`
    CreatedAt      time.Time      `json:"created_at"`
    UpdatedAt      time.Time      `json:"updated_at"`
    DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}
