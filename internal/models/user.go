package models

import (
    "time"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "io"
    "errors"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

// User - модель пользователя с шифрованием ПДн (152-ФЗ)
type User struct {
    ID                uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Email             string         `json:"email" gorm:"uniqueIndex;not null"`
    PasswordHash      string         `json:"-" gorm:"not null"` // bcrypt hash
    FirstName         string         `json:"first_name"`
    LastName          string         `json:"last_name"`
    PhoneEncrypted    []byte         `json:"-"` // AES-256 зашифрованный телефон
    Role              string         `json:"role" gorm:"index;not null;default:'patient'"`
    SubscriptionTier  string         `json:"subscription_tier" gorm:"default:'free'"`
    TwoFactorSecret   string         `json:"-"` // TOTP secret
    TwoFactorEnabled  bool           `json:"two_factor_enabled" gorm:"default:false"`
    LastLoginAt       *time.Time     `json:"last_login_at"`
    LoginAttempts     int            `json:"-" gorm:"default:0"`
    LockedUntil       *time.Time     `json:"-"`
    CreatedAt         time.Time      `json:"created_at"`
    UpdatedAt         time.Time      `json:"updated_at"`
    DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}

// EncryptPhone - шифрование телефона перед сохранением (AES-256-GCM)
func (u *User) EncryptPhone(phone string, key []byte) error {
    block, err := aes.NewCipher(key)
    if err != nil {
        return err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return err
    }

    u.PhoneEncrypted = gcm.Seal(nonce, nonce, []byte(phone), nil)
    return nil
}

// DecryptPhone - расшифровка телефона при чтении
func (u *User) DecryptPhone(key []byte) (string, error) {
    if len(u.PhoneEncrypted) == 0 {
        return "", errors.New("phone not encrypted")
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonceSize := gcm.NonceSize()
    if len(u.PhoneEncrypted) < nonceSize {
        return "", errors.New("invalid ciphertext")
    }

    nonce, ciphertext := u.PhoneEncrypted[:nonceSize], u.PhoneEncrypted[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}

// BeforeCreate - хук GORM перед созданием
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return nil
}

// Subscription - модель подписки с лимитами
type Subscription struct {
    ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID          uuid.UUID      `json:"user_id" gorm:"index;not null"`
    Tier            string         `json:"tier" gorm:"not null"` // free, patient_pro, doctor_pro, clinic
    MaxPatients     int            `json:"max_patients"`
    MaxReminders    int            `json:"max_reminders"`     // в месяц
    MaxAnalyses     int            `json:"max_analyses"`      // в месяц
    StorageYears    int            `json:"storage_years"`
    PriceMonthly    float64        `json:"price_monthly"`
    PriceYearly     float64        `json:"price_yearly"`
    Features        JSONStringMap  `json:"features" gorm:"type:jsonb"`
    StartedAt       time.Time      `json:"started_at"`
    ExpiresAt       time.Time      `json:"expires_at"`
    AutoRenew       bool           `json:"auto_renew" gorm:"default:true"`
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
}

// JSONStringMap - для хранения JSON в PostgreSQL
type JSONStringMap map[string]interface{}

// SubscriptionLimits - лимиты по тарифам (для проверки)
var SubscriptionLimits = map[string]map[string]interface{}{
    "free": {
        "max_patients":  1,
        "max_reminders": 10,
        "max_analyses":  5,
        "storage_years": 1,
        "price_monthly": 0,
        "price_yearly":  0,
        "features": map[string]bool{
            "telegram_bot":     true,
            "basic_reminders":  true,
            "symptom_diary":    true,
        },
    },
    "patient_pro": {
        "max_patients":  5, // семья
        "max_reminders": -1, // безлимит
        "max_analyses":  -1,
        "storage_years": 3,
        "price_monthly": 299,
        "price_yearly":  2990,
        "features": map[string]bool{
            "telegram_bot":          true,
            "unlimited_reminders":   true,
            "voice_input":           true,
            "photo_recognition":     true,
            "family_access":         true,
            "export_history":        true,
        },
    },
    "doctor_pro": {
        "max_patients":  100,
        "max_reminders": -1,
        "max_analyses":  -1,
        "storage_years": 5,
        "price_monthly": 1490,
        "price_yearly":  14900,
        "features": map[string]bool{
            "telegram_bot":          true,
            "unlimited_reminders":   true,
            "voice_input":           true,
            "photo_recognition":     true,
            "patient_management":    true,
            "prescription_templates": true,
            "analytics":             true,
            "mis_integration":       true,
        },
    },
    "clinic": {
        "max_patients":  1000,
        "max_reminders": -1,
        "max_analyses":  -1,
        "storage_years": 10,
        "price_monthly": 9900,
        "price_yearly":  99000,
        "features": map[string]bool{
            "telegram_bot":          true,
            "unlimited_reminders":   true,
            "voice_input":           true,
            "photo_recognition":     true,
            "patient_management":    true,
            "prescription_templates": true,
            "analytics":             true,
            "mis_integration":       true,
            "staff_accounts":        true,
            "api_access":            true,
            "custom_branding":       true,
            "priority_support":      true,
        },
    },
}
