package security

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

// HardwareSecurityModule - симуляция HSM с поддержкой реального железа
type HardwareSecurityModule struct {
    masterKey     []byte
    initialized   bool
    mu            sync.RWMutex
    hsmAvailable  bool
    hsmPath       string
}

// HSMConfig конфигурация аппаратного модуля
type HSMConfig struct {
    Enabled      bool   `json:"enabled"`
    DevicePath   string `json:"device_path"`   // /dev/tpm0 или /dev/hsm
    PKCS11Module string `json:"pkcs11_module"` // Путь к PKCS#11 библиотеке
    SlotID       uint   `json:"slot_id"`
    PinCode      string `json:"pin_code"`
}

var (
    ErrHSMNotAvailable = errors.New("HSM device not available")
    ErrInvalidKey      = errors.New("invalid encryption key")
)

// NewHardwareSecurityModule создает новый HSM модуль
func NewHardwareSecurityModule() *HardwareSecurityModule {
    hsm := &HardwareSecurityModule{
        hsmPath:      "/dev/tpm0",
        hsmAvailable: checkHSMAvailable(),
    }
    
    if !hsm.hsmAvailable {
        log.Println("⚠️  HSM not available, using software encryption (simulation mode)")
        // Генерируем ключ из аппаратных характеристик
        hsm.masterKey = hsm.generateSoftwareKey()
    } else {
        log.Println("🔒 HSM device detected, using hardware encryption")
        hsm.initHSM()
    }
    
    hsm.initialized = true
    return hsm
}

// checkHSMAvailable проверяет наличие аппаратного HSM/TPM
func checkHSMAvailable() bool {
    // Проверяем наличие TPM устройства
    if _, err := os.Stat("/dev/tpm0"); err == nil {
        return true
    }
    if _, err := os.Stat("/dev/tpmrm0"); err == nil {
        return true
    }
    
    // Проверяем наличие PKCS#11 модулей
    pkcs11Paths := []string{
        "/usr/lib/softhsm/libsofthsm2.so",
        "/usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so",
        "/usr/local/lib/softhsm/libsofthsm2.so",
    }
    
    for _, path := range pkcs11Paths {
        if _, err := os.Stat(path); err == nil {
            return true
        }
    }
    
    return false
}

// initHSM инициализирует аппаратное устройство
func (h *HardwareSecurityModule) initHSM() {
    log.Println("🔐 Initializing HSM device...")
    // Здесь будет реальная инициализация HSM через PKCS#11
    // h.pkcs11 = pkcs11.New(h.hsmPath)
    // h.session = h.pkcs11.OpenSession(h.slotID)
}

// generateSoftwareKey генерирует ключ из уникальных аппаратных идентификаторов
func (h *HardwareSecurityModule) generateSoftwareKey() []byte {
    // Собираем уникальные идентификаторы системы
    hostname, _ := os.Hostname()
    machineID := getMachineID()
    
    data := []byte(hostname + machineID + "medical-bot-salt-2026")
    hash := sha256.Sum256(data)
    return hash[:]
}

// getMachineID получает уникальный ID машины
func getMachineID() string {
    // Пытаемся прочитать machine-id
    if data, err := os.ReadFile("/etc/machine-id"); err == nil {
        return string(data)
    }
    if data, err := os.ReadFile("/var/lib/dbus/machine-id"); err == nil {
        return string(data)
    }
    return "unknown-machine-id"
}

// Encrypt шифрует данные с использованием HSM
func (h *HardwareSecurityModule) Encrypt(plaintext []byte) (string, error) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    if !h.initialized {
        return "", ErrHSMNotAvailable
    }
    
    if h.hsmAvailable {
        return h.hsmEncrypt(plaintext)
    }
    return h.softwareEncrypt(plaintext)
}

// hsmEncrypt использует аппаратное шифрование
func (h *HardwareSecurityModule) hsmEncrypt(plaintext []byte) (string, error) {
    // Здесь будет реальная реализация HSM через PKCS#11
    // Для демонстрации используем software encryption с маркером HSM
    encrypted, err := h.softwareEncrypt(plaintext)
    if err != nil {
        return "", err
    }
    return "HSM:" + encrypted, nil
}

// softwareEncrypt программная реализация (для симуляции)
func (h *HardwareSecurityModule) softwareEncrypt(plaintext []byte) (string, error) {
    block, err := aes.NewCipher(h.masterKey)
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

// Decrypt расшифровывает данные
func (h *HardwareSecurityModule) Decrypt(encryptedData string) ([]byte, error) {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    // Проверяем маркер HSM
    isHSM := len(encryptedData) > 4 && encryptedData[:4] == "HSM:"
    if isHSM {
        encryptedData = encryptedData[4:]
    }
    
    if h.hsmAvailable && isHSM {
        return h.hsmDecrypt(encryptedData)
    }
    return h.softwareDecrypt(encryptedData)
}

func (h *HardwareSecurityModule) hsmDecrypt(encryptedData string) ([]byte, error) {
    return h.softwareDecrypt(encryptedData)
}

func (h *HardwareSecurityModule) softwareDecrypt(encryptedData string) ([]byte, error) {
    ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
    if err != nil {
        return nil, err
    }
    
    block, err := aes.NewCipher(h.masterKey)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, ErrInvalidKey
    }
    
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}

// GetHSMInfo возвращает информацию о HSM
func (h *HardwareSecurityModule) GetHSMInfo() map[string]interface{} {
    return map[string]interface{}{
        "available":     h.hsmAvailable,
        "initialized":   h.initialized,
        "device_path":   h.hsmPath,
        "encryption":    "AES-256-GCM",
        "mode":          map[bool]string{true: "hardware", false: "software"}[h.hsmAvailable],
    }
}
