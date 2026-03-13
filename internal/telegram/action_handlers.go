package telegram

import (
    "fmt"
    "math/rand"
    "strings"
    "time"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Добавление лекарства
func (b *TelegramBot) handleMedAdd(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    chatID := update.CallbackQuery.Message.Chat.ID
    
    if user.UserID == "" {
        b.requireAuth(chatID)
        return
    }
    
    session.State = StateAwaitingMedication
    b.db.Save(session)
    
    msg := "💊 Введите название лекарства:"
    b.editMessage(chatID, update.CallbackQuery.Message.MessageID, msg, nil)
}

func (b *TelegramBot) handleMedicationInput(update tgbotapi.Update, user *TelegramUser, session *TelegramSession, text string) {
    chatID := update.Message.Chat.ID
    
    session.State = StateAwaitingDosage
    session.TempData = text
    b.db.Save(session)
    
    msg := fmt.Sprintf("💊 Лекарство: <b>%s</b>\n\n", text)
    msg += "Введите дозировку (например: 500 мг):"
    b.sendMessage(chatID, msg, nil)
}

func (b *TelegramBot) handleDosageInput(update tgbotapi.Update, user *TelegramUser, session *TelegramSession, text string) {
    chatID := update.Message.Chat.ID
    medication := session.TempData
    
    session.State = StateAwaitingFrequency
    session.TempData = medication + "|" + text
    b.db.Save(session)
    
    msg := fmt.Sprintf("💊 <b>%s</b> %s\n\n", medication, text)
    msg += "Как часто принимать? (например: 3 раза в день)"
    b.sendMessage(chatID, msg, nil)
}

func (b *TelegramBot) handleFrequencyInput(update tgbotapi.Update, user *TelegramUser, session *TelegramSession, text string) {
    chatID := update.Message.Chat.ID
    data := strings.Split(session.TempData, "|")
    medication := data[0]
    dosage := data[1]
    
    session.State = StateAwaitingDuration
    session.TempData = medication + "|" + dosage + "|" + text
    b.db.Save(session)
    
    msg := fmt.Sprintf("💊 <b>%s</b> %s\n", medication, dosage)
    msg += fmt.Sprintf("📅 Частота: %s\n\n", text)
    msg += "Сколько дней принимать? (например: 7 дней)"
    b.sendMessage(chatID, msg, nil)
}

func (b *TelegramBot) handleDurationInput(update tgbotapi.Update, user *TelegramUser, session *TelegramSession, text string) {
    chatID := update.Message.Chat.ID
    data := strings.Split(session.TempData, "|")
    medication := data[0]
    dosage := data[1]
    frequency := data[2]
    
    // Заглушка сохранения
    fmt.Printf("Сохранение: %s %s %s %s\n", medication, dosage, frequency, text)
    
    session.State = StateNone
    session.TempData = ""
    b.db.Save(session)
    
    msg := fmt.Sprintf("✅ Лекарство <b>%s</b> добавлено!\n\n", medication)
    b.sendMessage(chatID, msg, MedicationsMenu(true))
}

// Выбор симптома
func (b *TelegramBot) handleSymptomSelect(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    chatID := update.CallbackQuery.Message.Chat.ID
    data := update.CallbackQuery.Data
    symptom := strings.TrimPrefix(data, "symptom_")
    
    session.TempData = symptom
    b.db.Save(session)
    
    msg := fmt.Sprintf("📝 Симптом: <b>%s</b>\n\n", symptom)
    msg += "Оцените интенсивность:"
    
    b.editMessage(chatID, update.CallbackQuery.Message.MessageID, msg, IntensityMenu(symptom))
}

func (b *TelegramBot) handleIntensity(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    chatID := update.CallbackQuery.Message.Chat.ID
    data := update.CallbackQuery.Data
    parts := strings.Split(data, "_")
    intensity := parts[1]
    symptom := parts[2]
    
    fmt.Printf("Сохранение симптома: %s, интенсивность %s\n", symptom, intensity)
    
    msg := fmt.Sprintf("✅ Симптом <b>%s</b> (интенсивность %s/10) записан!", symptom, intensity)
    
    session.State = StateNone
    session.TempData = ""
    b.db.Save(session)
    
    b.editMessage(chatID, update.CallbackQuery.Message.MessageID, msg, SymptomsMenu())
}

func (b *TelegramBot) handleSymptomCustom(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    chatID := update.CallbackQuery.Message.Chat.ID
    
    session.State = StateAwaitingSymptoms
    b.db.Save(session)
    
    b.editMessage(chatID, update.CallbackQuery.Message.MessageID, "📝 Опишите ваш симптом:", nil)
}

func (b *TelegramBot) handleSymptomsInput(update tgbotapi.Update, user *TelegramUser, session *TelegramSession, text string) {
    chatID := update.Message.Chat.ID
    
    session.TempData = text
    b.db.Save(session)
    
    msg := fmt.Sprintf("📝 Симптом: <b>%s</b>\n\n", text)
    msg += "Оцените интенсивность:"
    
    b.sendMessage(chatID, msg, IntensityMenu(text))
}
