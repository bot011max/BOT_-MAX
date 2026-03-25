package ocr

import (
    "log"
    "regexp"
    "strconv"
)

type OCRService struct {
    language string
    active   bool
}

type Medication struct {
    Name    string
    Dosage  string
    Form    string
    Count   int
}

func NewOCRService() *OCRService {
    log.Println("📸 OCR service initialized (simulation mode)")
    return &OCRService{
        language: "rus",
        active:   true,
    }
}

func (o *OCRService) ExtractText(image []byte) (string, error) {
    // Симуляция OCR распознавания
    text := "Амоксициллин 500 мг №20\nПринимать 3 раза в день"
    log.Printf("📄 Extracted text: %s", text)
    return text, nil
}

func (o *OCRService) ExtractMedications(text string) []Medication {
    medications := []Medication{}
    
    // Паттерны для распознавания лекарств
    patterns := []struct {
        name   string
        regex  string
        dosage string
        form   string
    }{
        {"Амоксициллин", `Амоксициллин\s+(\d+)\s*мг`, "мг", "таблетки"},
        {"Парацетамол", `Парацетамол\s+(\d+)\s*мг`, "мг", "таблетки"},
        {"Ибупрофен", `Ибупрофен\s+(\d+)\s*мг`, "мг", "таблетки"},
        {"Цитрамон", `Цитрамон`, "500 мг", "таблетки"},
        {"Нурофен", `Нурофен\s+(\d+)\s*мг`, "мг", "капсулы"},
    }
    
    for _, pattern := range patterns {
        re := regexp.MustCompile(pattern.regex)
        matches := re.FindAllStringSubmatch(text, -1)
        
        for _, match := range matches {
            dosage := pattern.dosage
            if len(match) > 1 {
                dosage = match[1] + " " + pattern.dosage
            }
            
            // Поиск количества
            countRe := regexp.MustCompile(`№\s*(\d+)`)
            countMatch := countRe.FindStringSubmatch(text)
            count := 20 // по умолчанию
            if len(countMatch) > 1 {
                count, _ = strconv.Atoi(countMatch[1])
            }
            
            medications = append(medications, Medication{
                Name:   pattern.name,
                Dosage: dosage,
                Form:   pattern.form,
                Count:  count,
            })
        }
    }
    
    log.Printf("💊 Extracted %d medications", len(medications))
    return medications
}
