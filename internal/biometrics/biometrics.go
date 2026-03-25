package biometrics

import (
    "crypto/sha256"
    "encoding/hex"
    "log"
)

type VoiceBiometrics struct {
    profiles  map[string]string
    threshold float64
}

func NewVoiceBiometrics() *VoiceBiometrics {
    return &VoiceBiometrics{
        profiles:  make(map[string]string),
        threshold: 0.95,
    }
}

func (v *VoiceBiometrics) RegisterVoice(userID string, audioHash string) {
    v.profiles[userID] = audioHash
    log.Printf("🎤 Voice profile registered for user %s", userID)
}

func (v *VoiceBiometrics) VerifyVoice(userID string, audioHash string) bool {
    stored, exists := v.profiles[userID]
    if !exists {
        return false
    }
    return stored == audioHash
}

type FaceBiometrics struct {
    profiles  map[string]string
    threshold float64
}

func NewFaceBiometrics() *FaceBiometrics {
    return &FaceBiometrics{
        profiles:  make(map[string]string),
        threshold: 0.85,
    }
}

func (f *FaceBiometrics) RegisterFace(userID string, faceHash string) {
    f.profiles[userID] = faceHash
    log.Printf("👤 Face profile registered for user %s", userID)
}

func (f *FaceBiometrics) VerifyFace(userID string, faceHash string) bool {
    stored, exists := f.profiles[userID]
    if !exists {
        return false
    }
    return stored == faceHash
}

func HashData(data []byte) string {
    hash := sha256.Sum256(data)
    return hex.EncodeToString(hash[:])
}
