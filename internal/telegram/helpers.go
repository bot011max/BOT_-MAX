package telegram

import (
    "encoding/json"
    "fmt"
    "log"
    "time"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/google/uuid"
)

// Получить или создать пользователя
func (b *TelegramBot) getOrCreateUser(update tgbotapi.Update) (*TelegramUser, error) {
    var from *tgbotapi.User
    if update.Message != nil {
        from = update.Message.From
    } else if update.CallbackQuery != nil {
        from = update.CallbackQuery.From
    } else {
        return nil, fmt.Errorf("не удалось определить пользователя")
    }
    
    var user TelegramUser
    err := b.db.Where("telegram_id = ?", from.ID).First(&user).Error
    
    if err == nil {
        return &user, nil
    }
    
    // Создаем нового пользователя
    user = TelegramUser{
        TelegramID:   from.ID,
        ChatID:       update.Message.Chat.ID,
        Username:     from.UserName,
        FirstName:    from.FirstName,
        LastName:     from.LastName,
        LanguageCode: from.LanguageCode,
        IsActive:     true,
        CreatedAt:    time.Now(),
    }
    
    if err := b.db.Create(&user).Error; err != nil {
        return nil, err
    }
    
    return &user, nil
}

// Получить или создать сессию
func (b *TelegramBot) getOrCreateSession(telegramID int64) (*TelegramSession, error) {
    var session TelegramSession
    err := b.db.Where("telegram_id = ?", telegramID).First(&session).Error
    
    if err == nil {
        return &session, nil
    }
    
    session = TelegramSession{
        TelegramID: telegramID,
        State:      StateNone,
        UpdatedAt:  time.Now(),
    }
    
    if err := b.db.Create(&session).Error; err != nil {
        return nil, err
    }
    
    return &session, nil
}

// Отправить сообщение
func (b *TelegramBot) sendMessage(chatID int64, text string, keyboard interface{}) {
    msg := tgbotapi.NewMessage(chatID, text)
    msg.ParseMode = "HTML"
    
    if keyboard != nil {
        msg.ReplyMarkup = keyboard
    }
    
    if _, err := b.api.Send(msg); err != nil {
        log.Printf("Ошибка отправки сообщения: %v", err)
    }
}

// Отредактировать сообщение
func (b *TelegramBot) editMessage(chatID int64, messageID int, text string, keyboard interface{}) {
    msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
    msg.ParseMode = "HTML"
    
    if keyboard != nil {
        msg.ReplyMarkup = keyboard.(*tgbotapi.InlineKeyboardMarkup)
    }
    
    if _, err := b.api.Send(msg); err != nil {
        log.Printf("Ошибка редактирования сообщения: %v", err)
    }
}

// Требовать авторизацию
func (b *TelegramBot) requireAuth(chatID int64) {
    msg := "🔐 Для использования этой функции необходимо авторизоваться.\n\nИспользуйте /login"
    b.sendMessage(chatID, msg, nil)
}

// Получить лекарства пользователя
func (b *TelegramBot) getUserMedications(userID string) []Medication {
    // TODO: реализовать запрос к БД
    return []Medication{}
}

// Получить лекарства на сегодня
func (b *TelegramBot) getTodayMedications(userID string) []TodayMedication {
    // TODO: реализовать запрос к БД
    return []TodayMedication{}
}

// Получить приемы пользователя
func (b *TelegramBot) getUserAppointments(userID string) []Appointment {
    // TODO: реализовать запрос к БД
    return []Appointment{}
}

// Получить приемы на сегодня
func (b *TelegramBot) getTodayAppointments(userID string) []TodayAppointment {
    // TODO: реализовать запрос к БД
    return []TodayAppointment{}
}

// Получить врачей пользователя
func (b *TelegramBot) getUserDoctors(userID string) []Doctor {
    // TODO: реализовать запрос к БД
    return []Doctor{}
}

// Сохранить лекарство
func (b *TelegramBot) saveMedication(userID, name, dosage, frequency, duration string) {
    // TODO: реализовать сохранение в БД
}

// Сохранить симптом
func (b *TelegramBot) saveSymptom(userID, symptom, intensity string) {
    // TODO: реализовать сохранение в БД
}

// Структуры данных
type Medication struct {
    ID        int
    Name      string
    Dosage    string
    Frequency string
    EndDate   time.Time
}

type TodayMedication struct {
    ID     int
    Name   string
    Dosage string
    Time   string
    Taken  bool
}

type Appointment struct {
    Date       time.Time
    DoctorName string
    Specialty  string
    Location   string
}

type TodayAppointment struct {
    Time      string
    DoctorName string
    Specialty string
}

type Doctor struct {
    FullName  string
    Specialty string
    Phone     string
}
