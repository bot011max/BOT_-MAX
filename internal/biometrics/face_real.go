package biometrics

import (
    "crypto/sha256"
    "encoding/base64"
    "log"
)

type FaceBiometricsReal struct {
    profiles  map[string]string
    threshold float64
}

func NewFaceBiometricsReal() *FaceBiometricsReal {
    log.Println("👤 Face biometrics initialized (simulation mode)")
    return &FaceBiometricsReal{
        profiles:  make(map[string]string),
        threshold: 0.85,
    }
}

func (f *FaceBiometricsReal) Register(userID string, faceHash string) {
    f.profiles[userID] = faceHash
    log.Printf("👤 Face profile registered for user: %s", userID)
}

func (f *FaceBiometricsReal) Verify(userID string, faceHash string) bool {
    stored, exists := f.profiles[userID]
    if !exists {
        return false
    }
    
    similarity := calculateHashSimilarity(stored, faceHash)
    return similarity > f.threshold
}

func (f *FaceBiometricsReal) ExtractEmbedding(image []byte) string {
    hash := sha256.Sum256(image)
    return base64.StdEncoding.EncodeToString(hash[:])
}

// calculateHashSimilarity вычисляет схожесть между двумя хэшами
func calculateHashSimilarity(hash1, hash2 string) float64 {
    if hash1 == hash2 {
        return 1.0
    }
    
    // Если хэши разной длины, считаем их разными
    if len(hash1) != len(hash2) {
        return 0.0
    }
    
    // Считаем количество совпадающих символов
    matches := 0
    for i := 0; i < len(hash1); i++ {
        if hash1[i] == hash2[i] {
            matches++
        }
    }
    
    // Возвращаем процент совпадения
    return float64(matches) / float64(len(hash1))
}
