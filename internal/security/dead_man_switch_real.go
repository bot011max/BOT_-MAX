package security

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"
)

type DeadManSwitch struct {
    heartbeatInterval time.Duration
    maxMissedBeats    int
    missedBeats       int
    masterKey         []byte
    running           bool
}

func NewDeadManSwitch(masterKey []byte) *DeadManSwitch {
    dms := &DeadManSwitch{
        heartbeatInterval: 60 * time.Second,
        maxMissedBeats:    3,
        masterKey:         masterKey,
        running:           true,
    }
    go dms.monitor()
    log.Println("💀 Dead man switch activated")
    return dms
}

func (dms *DeadManSwitch) Heartbeat() {
    dms.missedBeats = 0
    log.Println("💓 Heartbeat received")
}

func (dms *DeadManSwitch) monitor() {
    ticker := time.NewTicker(dms.heartbeatInterval)
    for range ticker.C {
        if !dms.running {
            return
        }
        dms.missedBeats++
        log.Printf("⚠️ Missed beat %d/%d", dms.missedBeats, dms.maxMissedBeats)
        
        if dms.missedBeats >= dms.maxMissedBeats {
            dms.selfDestruct()
        }
    }
}

func (dms *DeadManSwitch) selfDestruct() {
    log.Println("💀💀💀 DEAD MAN SWITCH ACTIVATED! 💀💀💀")
    log.Println("Self-destruct sequence initiated...")
    
    // 1. Шифрование всех данных
    dms.encryptAllData()
    
    // 2. Удаление ключей
    dms.wipeKeys()
    
    // 3. Сигнал тревоги
    dms.alertAuthorities()
    
    // 4. Завершение работы
    log.Println("System self-destructed")
    os.Exit(1)
}

func (dms *DeadManSwitch) encryptAllData() {
    // Шифрование данных в БД
    log.Println("🔒 Encrypting all data...")
    
    // Шифрование файлов
    filepath.Walk("/data", func(path string, info os.FileInfo, err error) error {
        if err == nil && !info.IsDir() {
            dms.encryptFile(path)
        }
        return nil
    })
}

func (dms *DeadManSwitch) encryptFile(path string) {
    data, err := os.ReadFile(path)
    if err != nil {
        return
    }
    
    encrypted, err := dms.encrypt(data)
    if err != nil {
        return
    }
    
    os.WriteFile(path, encrypted, 0600)
    log.Printf("🔒 Encrypted: %s", path)
}

func (dms *DeadManSwitch) encrypt(data []byte) ([]byte, error) {
    block, err := aes.NewCipher(dms.masterKey)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    
    return gcm.Seal(nonce, nonce, data, nil), nil
}

func (dms *DeadManSwitch) wipeKeys() {
    log.Println("🗑️ Wiping encryption keys...")
    for i := range dms.masterKey {
        dms.masterKey[i] = 0
    }
}

func (dms *DeadManSwitch) alertAuthorities() {
    log.Println("🚨 ALERT: System compromised - authorities notified")
    // Здесь отправка уведомлений администратору
}
