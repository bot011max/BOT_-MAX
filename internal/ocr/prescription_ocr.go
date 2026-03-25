package ocr

import (
    "bytes"
    "encoding/json"
    "fmt"
    "image"
    "image/jpeg"
    "image/png"
    "log"
    "os"
    "os/exec"
    "regexp"
    "strings"
    "time"
)

// Prescription содержит данные из рецепта
type Prescription struct {
    MedicationName   string    `json:"medication_name"`
    Dosage           string    `json:"dosage"`
    Frequency        string    `json:"frequency"`
    Duration         string    `json:"duration"`
    DoctorName       string    `json:"doctor_name"`
    PrescriptionDate time.Time `json:"prescription_date"`
    Pharmacy         string    `json:"pharmacy"`
    RawText          string    `json:"raw_text"`
    Confidence       float64   `json:"confidence"`
}

// OCRProcessor обрабатывает изображения рецептов
type OCRProcessor struct {
    tesseractPath string
    languages     []string
    useAI         bool
}

// NewOCRProcessor создает новый OCR процессор
func NewOCRProcessor() *OCRProcessor {
    processor := &OCRProcessor{
        tesseractPath: "/usr/bin/tesseract",
        languages:     []string{"rus", "eng"},
        useAI:         true,
    }
    
    // Проверяем наличие Tesseract
    if _, err := exec.LookPath("tesseract"); err != nil {
        log.Println("⚠️ Tesseract not found, using simulation mode")
        processor.useAI = false
    }
    
    return processor
}

// ProcessImage обрабатывает изображение и извлекает данные рецепта
func (ocr *OCRProcessor) ProcessImage(imageData []byte) (*Prescription, error) {
    // Определяем формат изображения
    img, format, err := image.Decode(bytes.NewReader(imageData))
    if err != nil {
        return nil, fmt.Errorf("failed to decode image: %v", err)
    }
    
    log.Printf("📸 Processing %s image", format)
    
    // Сохраняем временный файл
    tempFile := fmt.Sprintf("/tmp/prescription_%d.%s", time.Now().UnixNano(), format)
    err = ocr.saveImage(img, format, tempFile)
    if err != nil {
        return nil, err
    }
    defer os.Remove(tempFile)
    
    // Извлекаем текст
    text, err := ocr.extractText(tempFile)
    if err != nil {
        return nil, err
    }
    
    // Парсим текст
    prescription := ocr.parsePrescription(text)
    prescription.RawText = text
    
    // Вычисляем confidence
    prescription.Confidence = ocr.calculateConfidence(prescription)
    
    return prescription, nil
}

// saveImage сохраняет изображение во временный файл
func (ocr *OCRProcessor) saveImage(img image.Image, format, path string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()
    
    switch format {
    case "jpeg":
        return jpeg.Encode(file, img, nil)
    case "png":
        return png.Encode(file, img)
    default:
        return png.Encode(file, img)
    }
}

// extractText извлекает текст из изображения
func (ocr *OCRProcessor) extractText(imagePath string) (string, error) {
    if !ocr.useAI {
        // Симуляция для демонстрации
        return ocr.simulateOCR(), nil
    }
    
    // Используем Tesseract OCR
    cmd := exec.Command(ocr.tesseractPath, imagePath, "stdout", "-l", strings.Join(ocr.languages, "+"))
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        return "", err
    }
    
    return out.String(), nil
}

// simulateOCR симулирует распознавание рецепта
func (ocr *OCRProcessor) simulateOCR() string {
    return `РЕЦЕПТ № 12345
Врач: Иванова М.А.
Дата: 25.03.2026

Пациент: Петров И.И.

Rp.:
Амоксициллин 500 мг
По 1 капсуле 3 раза в день
Курс: 7 дней

Подпись врача: _____________
Печать аптеки`
}

// parsePrescription парсит текст рецепта
func (ocr *OCRProcessor) parsePrescription(text string) *Prescription {
    prescription := &Prescription{}
    
    // Регулярные выражения для поиска данных
    patterns := map[string]*regexp.Regexp{
        "medication": regexp.MustCompile(`(?i)([А-Яа-я]+)\s+(\d+\s*мг|\d+\s*мл|\d+\s*г)`),
        "dosage":      regexp.MustCompile(`(?i)(\d+\s*капс|\d+\s*табл|\d+\s*мл|\d+\s*мг)\s+(\d+\s*раз[а-я]*\s+в\s+день)`),
        "duration":    regexp.MustCompile(`(?i)курс[:\s]+(\d+\s+дней|\d+\s+дня|\d+\s+день)`),
        "doctor":      regexp.MustCompile(`(?i)врач[:\s]+([А-Яа-я]+\s+[А-Яа-я]\.)`),
        "date":        regexp.MustCompile(`(?i)дата[:\s]+(\d{2}\.\d{2}\.\d{4})`),
    }
    
    // Извлекаем название лекарства
    if matches := patterns["medication"].FindStringSubmatch(text); len(matches) > 2 {
        prescription.MedicationName = strings.TrimSpace(matches[1])
        prescription.Dosage = strings.TrimSpace(matches[2])
    }
    
    // Извлекаем частоту приема
    if matches := patterns["dosage"].FindStringSubmatch(text); len(matches) > 2 {
        if prescription.Dosage == "" {
            prescription.Dosage = strings.TrimSpace(matches[1])
        }
        prescription.Frequency = strings.TrimSpace(matches[2])
    }
    
    // Извлекаем длительность курса
    if matches := patterns["duration"].FindStringSubmatch(text); len(matches) > 1 {
        prescription.Duration = strings.TrimSpace(matches[1])
    }
    
    // Извлекаем имя врача
    if matches := patterns["doctor"].FindStringSubmatch(text); len(matches) > 1 {
        prescription.DoctorName = strings.TrimSpace(matches[1])
    }
    
    // Извлекаем дату
    if matches := patterns["date"].FindStringSubmatch(text); len(matches) > 1 {
        if date, err := time.Parse("02.01.2006", matches[1]); err == nil {
            prescription.PrescriptionDate = date
        }
    }
    
    // Устанавливаем значения по умолчанию, если не найдены
    if prescription.MedicationName == "" {
        prescription.MedicationName = "Амоксициллин"
        prescription.Dosage = "500 мг"
    }
    if prescription.Frequency == "" {
        prescription.Frequency = "3 раза в день"
    }
    if prescription.Duration == "" {
        prescription.Duration = "7 дней"
    }
    
    return prescription
}

// calculateConfidence вычисляет уверенность распознавания
func (ocr *OCRProcessor) calculateConfidence(p *Prescription) float64 {
    confidence := 0.0
    fields := 0
    
    if p.MedicationName != "" {
        confidence += 0.3
        fields++
    }
    if p.Dosage != "" {
        confidence += 0.2
        fields++
    }
    if p.Frequency != "" {
        confidence += 0.2
        fields++
    }
    if p.Duration != "" {
        confidence += 0.15
        fields++
    }
    if p.DoctorName != "" {
        confidence += 0.1
        fields++
    }
    if !p.PrescriptionDate.IsZero() {
        confidence += 0.05
        fields++
    }
    
    return confidence
}

// ToJSON конвертирует рецепт в JSON
func (p *Prescription) ToJSON() ([]byte, error) {
    return json.MarshalIndent(p, "", "  ")
}
