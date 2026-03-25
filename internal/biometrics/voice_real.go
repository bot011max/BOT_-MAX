package biometrics

import (
    "crypto/sha256"
    "encoding/hex"
    "log"
)

type VoiceBiometricsReal struct {
    profiles  map[string]string
    threshold float64
}

func NewVoiceBiometricsReal() *VoiceBiometricsReal {
    log.Println("🎤 Voice biometrics initialized (simulation mode)")
    return &VoiceBiometricsReal{
        profiles:  make(map[string]string),
        threshold: 0.95,
    }
}

func (v *VoiceBiometricsReal) Register(userID string, audioHash string) {
    v.profiles[userID] = audioHash
    log.Printf("🎤 Voice profile registered for user: %s", userID)
}

func (v *VoiceBiometricsReal) Verify(userID string, audioHash string) bool {
    stored, exists := v.profiles[userID]
    if !exists {
        return false
    }
    return stored == audioHash
}

func (v *VoiceBiometricsReal) ExtractFeatures(audio []byte) string {
    hash := sha256.Sum256(audio)
    return hex.EncodeToString(hash[:])
}
