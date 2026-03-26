package database

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/hex"
    "io"
    "os"
)

type EncryptedDB struct {
    key      []byte
    dbPath   string
}

func NewEncryptedDB(dbPath string) (*EncryptedDB, error) {
    keyHex := os.Getenv("DB_ENCRYPTION_KEY")
    if keyHex == "" {
        key := make([]byte, 32)
        if _, err := rand.Read(key); err != nil {
            return nil, err
        }
        keyHex = hex.EncodeToString(key)
    }
    
    key, err := hex.DecodeString(keyHex)
    if err != nil {
        return nil, err
    }
    
    return &EncryptedDB{
        key:    key,
        dbPath: dbPath,
    }, nil
}

func (e *EncryptedDB) encrypt(data []byte) ([]byte, error) {
    block, err := aes.NewCipher(e.key)
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

func (e *EncryptedDB) decrypt(data []byte) ([]byte, error) {
    block, err := aes.NewCipher(e.key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
        return nil, err
    }
    
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}

func (e *EncryptedDB) EncryptFile() error {
    data, err := os.ReadFile(e.dbPath)
    if err != nil {
        return err
    }
    
    encrypted, err := e.encrypt(data)
    if err != nil {
        return err
    }
    
    return os.WriteFile(e.dbPath+".enc", encrypted, 0600)
}

func (e *EncryptedDB) DecryptFile() error {
    encrypted, err := os.ReadFile(e.dbPath + ".enc")
    if err != nil {
        return err
    }
    
    decrypted, err := e.decrypt(encrypted)
    if err != nil {
        return err
    }
    
    return os.WriteFile(e.dbPath, decrypted, 0600)
}
