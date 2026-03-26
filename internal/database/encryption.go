package database

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "io"
    "os"
)

var encryptionKey []byte

func init() {
    // Загрузка ключа из переменной окружения
    keyStr := os.Getenv("DB_ENCRYPTION_KEY")
    if keyStr != "" {
        key, err := base64.StdEncoding.DecodeString(keyStr)
        if err == nil {
            encryptionKey = key
        }
    }
    
    if encryptionKey == nil {
        // Генерация ключа из мастер-ключа
        masterKey := os.Getenv("MASTER_KEY")
        if masterKey != "" {
            hash := sha256.Sum256([]byte(masterKey))
            encryptionKey = hash[:]
        }
    }
}

type EncryptedField struct {
    Data  string `json:"data"`
    Nonce string `json:"nonce"`
}

func EncryptData(plaintext []byte) (string, error) {
    if len(encryptionKey) == 0 {
        return string(plaintext), nil
    }
    
    block, err := aes.NewCipher(encryptionKey)
    if err != nil {
        return "", err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptData(encrypted string) ([]byte, error) {
    if len(encryptionKey) == 0 {
        return []byte(encrypted), nil
    }
    
    ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
    if err != nil {
        return nil, err
    }
    
    block, err := aes.NewCipher(encryptionKey)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }
    
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}
