package recovery

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "io"
    "log"
    "os"
    "time"
)

type AutoBackup struct {
    backupDir string
    interval  time.Duration
    masterKey []byte
}

func NewAutoBackup(masterKey []byte) *AutoBackup {
    return &AutoBackup{
        backupDir: "/var/backups/medical-bot",
        interval:  1 * time.Hour,
        masterKey: masterKey,
    }
}

func (b *AutoBackup) Start() {
    // Создаем директорию для бэкапов
    os.MkdirAll(b.backupDir, 0700)
    
    ticker := time.NewTicker(b.interval)
    go func() {
        for range ticker.C {
            b.createBackup()
        }
    }()
    log.Println("📦 Auto backup started, interval:", b.interval)
}

func (b *AutoBackup) createBackup() {
    timestamp := time.Now().Format("20060102_150405")
    backupFile := b.backupDir + "/backup_" + timestamp + ".enc"
    
    // Симуляция бэкапа БД
    data := []byte("Database dump simulation")
    
    // Шифрование
    encrypted, err := b.encrypt(data)
    if err != nil {
        log.Printf("❌ Backup encryption failed: %v", err)
        return
    }
    
    // Сохранение
    err = os.WriteFile(backupFile, encrypted, 0600)
    if err != nil {
        log.Printf("❌ Backup save failed: %v", err)
        return
    }
    
    log.Printf("✅ Backup created: %s (%d bytes)", backupFile, len(encrypted))
    
    // Очистка старых бэкапов (хранить 7 дней)
    b.cleanupOldBackups()
}

func (b *AutoBackup) encrypt(data []byte) ([]byte, error) {
    block, err := aes.NewCipher(b.masterKey)
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

func (b *AutoBackup) cleanupOldBackups() {
    files, err := os.ReadDir(b.backupDir)
    if err != nil {
        return
    }
    
    cutoff := time.Now().AddDate(0, 0, -7) // 7 дней
    for _, file := range files {
        info, _ := file.Info()
        if info.ModTime().Before(cutoff) {
            os.Remove(b.backupDir + "/" + file.Name())
            log.Printf("🗑️ Deleted old backup: %s", file.Name())
        }
    }
}
