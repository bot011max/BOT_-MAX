package telegram

import (
    "time"
)

// TelegramUser связывает пользователя системы с Telegram
type TelegramUser struct {
    ID           uint      `json:"id" gorm:"primaryKey"`
    UserID       string    `json:"user_id" gorm:"index;not null"`        // ID в системе
    TelegramID   int64     `json:"telegram_id" gorm:"uniqueIndex;not null"` // ID в Telegram
    ChatID       int64     `json:"chat_id" gorm:"not null"`              // ID чата
    Username     string    `json:"username"`                              // @username
    FirstName    string    `json:"first_name"`
    LastName     string    `json:"last_name"`
    LanguageCode string    `json:"language_code"`
    IsActive     bool      `json:"is_active" gorm:"default:true"`
    AuthToken    string    `json:"auth_token"`                           // Токен для привязки аккаунта
    TokenExpires *time.Time `json:"token_expires"`                        // Срок действия токена
    
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

// TelegramSession сессия диалога с пользователем
type TelegramSession struct {
    ID           uint      `json:"id" gorm:"primaryKey"`
    TelegramID   int64     `json:"telegram_id" gorm:"index;not null"`
    State        string    `json:"state"`                                // Текущее состояние диалога
    TempData     string    `json:"temp_data"`                            // Временные данные (JSON)
    LastMessageID int      `json:"last_message_id"`                      // ID последнего сообщения
    LastCommand  string    `json:"last_command"`                         // Последняя команда
    UpdatedAt    time.Time `json:"updated_at"`
}

// Reminder напоминание для отправки в Telegram
type Reminder struct {
    ID           uint      `json:"id" gorm:"primaryKey"`
    UserID       string    `json:"user_id" gorm:"index;not null"`
    TelegramID   int64     `json:"telegram_id" gorm:"index"`
    Message      string    `json:"message"`
    ScheduledFor time.Time `json:"scheduled_for" gorm:"index"`
    SentAt       *time.Time `json:"sent_at"`
    Status       string    `json:"status" gorm:"default:'pending'"` // pending, sent, failed
    RetryCount   int       `json:"retry_count" gorm:"default:0"`
    CreatedAt    time.Time `json:"created_at"`
}

// Коды состояний диалога
const (
    StateNone               = "none"
    StateAwaitingAuth       = "awaiting_auth"        // Ожидание кода авторизации
    StateAwaitingSymptoms   = "awaiting_symptoms"     // Ожидание описания симптомов
    StateAwaitingMedication = "awaiting_medication"   // Ожидание названия лекарства
    StateAwaitingDosage     = "awaiting_dosage"       // Ожидание дозировки
    StateAwaitingFrequency  = "awaiting_frequency"    // Ожидание частоты приема
    StateAwaitingDuration   = "awaiting_duration"     // Ожидание длительности
)
