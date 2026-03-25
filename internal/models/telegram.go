package models

import (
    "time"
    "github.com/google/uuid"
)

type TelegramUser struct {
    ID           uint      `json:"id" gorm:"primaryKey"`
    UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
    TelegramID   int64     `json:"telegram_id" gorm:"uniqueIndex;not null"`
    ChatID       int64     `json:"chat_id"`
    Username     string    `json:"username"`
    FirstName    string    `json:"first_name"`
    LastName     string    `json:"last_name"`
    IsActive     bool      `json:"is_active" gorm:"default:true"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
