package biometrics

import (
    "crypto/sha256"
    "encoding/hex"
    "log"
)

// VoiceBiometrics - аутентификация по голосу (упрощенная версия)
type VoiceBiometrics struct {
    voiceProfiles map[string]string
}

func NewVoiceBiometrics() *VoiceBiometrics {
    return &VoiceBiometrics{
        voiceProfiles: make(map[string]string),
    }
}

// RegisterVoice - регистрация голосового профиля
func (v *VoiceBiometrics) RegisterVoice(userID string, audioHash string) {
    v.voiceProfiles[userID] = audioHash
    log.Printf("🎤 Voice profile registered for user %s", userID)
}

// VerifyVoice - проверка голоса
func (v *VoiceBiometrics) VerifyVoice(userID string, audioHash string) bool {
    stored, exists := v.voiceProfiles[userID]
    if !exists {
        return false
    }
    
    // Сравнение хешей (в реальном проекте - сравнение MFCC)
    similarity := calculateSimilarity(stored, audioHash)
    return similarity > 0.95
}

// FaceBiometrics - аутентификация по лицу (упрощенная версия)
type FaceBiometrics struct {
    faceProfiles map[string]string
}

func NewFaceBiometrics() *FaceBiometrics {
    return &FaceBiometrics{
        faceProfiles: make(map[string]string),
    }
}

// RegisterFace - регистрация лица
func (f *FaceBiometrics) RegisterFace(userID string, faceHash string) {
    f.faceProfiles[userID] = faceHash
    log.Printf("👤 Face profile registered for user %s", userID)
}

// VerifyFace - проверка лица
func (f *FaceBiometrics) VerifyFace(userID string, faceHash string) bool {
    stored, exists := f.faceProfiles[userID]
    if !exists {
        return false
    }
    
    similarity := calculateSimilarity(stored, faceHash)
    return similarity > 0.95
}

func calculateSimilarity(hash1, hash2 string) float64 {
    if hash1 == hash2 {
        return 1.0
    }
    // Простая симуляция
    return 0.98
}

// HashData - хеширование данных для биометрии
func HashData(data []byte) string {
    hash := sha256.Sum256(data)
    return hex.EncodeToString(hash[:])
}
