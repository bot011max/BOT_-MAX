package auth

import (
    "crypto/rand"
    "encoding/base32"
    "image/png"
    "os"
    
    "github.com/pquerna/otp/totp"
    "github.com/skip2/go-qrcode"
)

type TwoFactorAuth struct {
    Secret   string   `json:"secret"`
    Verified bool     `json:"verified"`
    Codes    []string `json:"backup_codes,omitempty"`
}

// GenerateSecret генерирует секрет для TOTP
func GenerateSecret() (string, error) {
    secret := make([]byte, 20)
    _, err := rand.Read(secret)
    if err != nil {
        return "", err
    }
    return base32.StdEncoding.EncodeToString(secret), nil
}

// GenerateQRCode генерирует QR код для Google Authenticator
func GenerateQRCode(secret, email, issuer string) ([]byte, error) {
    url := totp.URL(secret, issuer, email)
    return qrcode.Encode(url, qrcode.Medium, 256)
}

// VerifyCode проверяет TOTP код
func VerifyCode(secret, code string) bool {
    return totp.Validate(code, secret)
}

// GenerateBackupCodes генерирует резервные коды
func GenerateBackupCodes(count int) []string {
    codes := make([]string, count)
    for i := 0; i < count; i++ {
        code := make([]byte, 5)
        rand.Read(code)
        codes[i] = base32.StdEncoding.EncodeToString(code)[:8]
    }
    return codes
}

// SaveQRCode сохраняет QR код в файл
func SaveQRCode(secret, email, issuer, filename string) error {
    qrData, err := GenerateQRCode(secret, email, issuer)
    if err != nil {
        return err
    }
    return os.WriteFile(filename, qrData, 0644)
}
