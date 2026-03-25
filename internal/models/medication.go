package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Medication struct {
    ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID       uuid.UUID      `json:"user_id" gorm:"type:uuid;index;not null"`
    Name         string         `json:"name" gorm:"not null"`
    Dosage       string         `json:"dosage"`
    Frequency    string         `json:"frequency"`
    Instructions string         `json:"instructions" gorm:"type:text"`
    StartDate    *time.Time     `json:"start_date"`
    EndDate      *time.Time     `json:"end_date"`
    IsActive     bool           `json:"is_active" gorm:"default:true"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (m *Medication) BeforeCreate(tx *gorm.DB) error {
    if m.ID == uuid.Nil {
        m.ID = uuid.New()
    }
    return nil
}
