package hsm

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "io"
    "log"
    "os"
    "sync"
)

type SoftwareHSM struct {
    masterKey []byte
    mu        sync.RWMutex
}

var ErrKeyNotFound = errors.New("key not found")

func NewSoftwareHSM() *SoftwareHSM {
    // Получаем ключ из переменной окружения или генерируем
    keyEnv := os.Getenv("HSM_MASTER_KEY")
    var masterKey []byte
    
    if keyEnv != "" {
        masterKey = []byte(keyEnv)
    } else {
        // Генерируем ключ из уникальных идентификаторов системы
        hostname, _ := os.Hostname()
        machineID := getMachineID()
        hash := sha256.Sum256([]byte(hostname + machineID + "medical-bot-salt-2026"))
        masterKey = hash[:]
    }
    
    log.Println("🔐 Software HSM initialized (simulation mode)")
    return &SoftwareHSM{
        masterKey: masterKey,
    }
}

func getMachineID() string {
    if data, err := os.ReadFile("/etc/machine-id"); err == nil {
        return string(data)
    }
    if data, err := os.ReadFile("/var/lib/dbus/machine-id"); err == nil {
        return string(data)
    }
    return "unknown-machine-id"
}

func (h *SoftwareHSM) Encrypt(data []byte) (string, error) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    block, err := aes.NewCipher(h.masterKey[:32])
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
    
    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (h *SoftwareHSM) Decrypt(encrypted string) ([]byte, error) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    data, err := base64.StdEncoding.DecodeString(encrypted)
    if err != nil {
        return nil, err
    }
    
    block, err := aes.NewCipher(h.masterKey[:32])
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }
    
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}

func (h *SoftwareHSM) GenerateKey(label string) (string, error) {
    key := make([]byte, 32)
    if _, err := rand.Read(key); err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(key), nil
}
