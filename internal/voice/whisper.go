package voice

import (
    "log"
    "regexp"
    "strings"
)

type WhisperService struct {
    model string
    active bool
}

type Symptom struct {
    Name     string
    Severity int
    Timestamp int64
}

func NewWhisperService() *WhisperService {
    log.Println("🎤 Whisper voice recognition service initialized (simulation mode)")
    return &WhisperService{
        model:  "base",
        active: true,
    }
}

func (w *WhisperService) Transcribe(audio []byte) (string, error) {
    // Симуляция распознавания речи
    text := "У меня болит голова и температура 38.5"
    log.Printf("🎤 Transcribed: %s", text)
    return text, nil
}

func (w *WhisperService) ExtractSymptoms(text string) []Symptom {
    symptoms := []Symptom{}
    text = strings.ToLower(text)
    
    symptomPatterns := map[string]int{
        "головная боль":   7,
        "голова":          6,
        "температура":     8,
        "кашель":          5,
        "насморк":         4,
        "боль в горле":    6,
        "тошнота":         5,
        "слабость":        3,
        "ломота":          4,
    }
    
    for symptom, severity := range symptomPatterns {
        if strings.Contains(text, symptom) {
            symptoms = append(symptoms, Symptom{
                Name:     symptom,
                Severity: severity,
            })
        }
    }
    
    // Поиск температуры
    re := regexp.MustCompile(`(\d+\.?\d*)`)
    matches := re.FindAllString(text, -1)
    for _, match := range matches {
        symptoms = append(symptoms, Symptom{
            Name:     "температура",
            Severity: int(match[0] - '0'),
        })
    }
    
    log.Printf("📊 Extracted %d symptoms", len(symptoms))
    return symptoms
}
