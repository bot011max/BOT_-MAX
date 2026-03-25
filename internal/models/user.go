package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID           string         `json:"id" gorm:"primaryKey;type:text"`
    Email        string         `json:"email" gorm:"uniqueIndex;not null"`
    PasswordHash string         `json:"-" gorm:"not null"`
    FirstName    string         `json:"first_name"`
    LastName     string         `json:"last_name"`
    Role         string         `json:"role" gorm:"default:'patient'"`
    Phone        string         `json:"phone"`
    IsActive     bool           `json:"is_active" gorm:"default:true"`
    TelegramID   *int64         `json:"telegram_id,omitempty" gorm:"uniqueIndex"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == "" {
        u.ID = uuid.New().String()
    }
    return nil
}
