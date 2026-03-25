package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Email        string         `json:"email" gorm:"uniqueIndex;not null;size:255"`
    PasswordHash string         `json:"-" gorm:"not null;size:255"`
    FirstName    string         `json:"first_name" gorm:"size:100"`
    LastName     string         `json:"last_name" gorm:"size:100"`
    Role         string         `json:"role" gorm:"default:'patient';size:50"`
    Phone        string         `json:"phone" gorm:"size:20"`
    IsActive     bool           `json:"is_active" gorm:"default:true"`
    TelegramID   *int64         `json:"telegram_id,omitempty" gorm:"uniqueIndex"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return nil
}

func (u *User) GetFullName() string {
    return u.FirstName + " " + u.LastName
}
