package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Medication struct {
    ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID       uuid.UUID      `json:"user_id" gorm:"type:uuid;index;not null"`
    Name         string         `json:"name" gorm:"not null;size:255"`
    Dosage       string         `json:"dosage" gorm:"size:100"`
    Frequency    string         `json:"frequency" gorm:"size:100"`
    Form         string         `json:"form" gorm:"size:50"`
    Instructions string         `json:"instructions" gorm:"type:text"`
    StartDate    *time.Time     `json:"start_date"`
    EndDate      *time.Time     `json:"end_date"`
    IsActive     bool           `json:"is_active" gorm:"default:true"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
    
    User         User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (m *Medication) BeforeCreate(tx *gorm.DB) error {
    if m.ID == uuid.Nil {
        m.ID = uuid.New()
    }
    return nil
}

func (m *Medication) IsExpired() bool {
    if m.EndDate == nil {
        return false
    }
    return time.Now().After(*m.EndDate)
}
