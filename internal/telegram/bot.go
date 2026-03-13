package telegram

import (
    "encoding/json"
    "fmt"
    "log"
    "strings"
    "time"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

// TelegramBot главная структура бота
type TelegramBot struct {
    api      *tgbotapi.BotAPI
    db       *gorm.DB
    handlers map[string]CommandHandler
    services *Services
}

type Services struct {
    prescription *PrescriptionService
    reminder     *ReminderService
    auth         *AuthService
}

type CommandHandler func(update tgbotapi.Update, user *TelegramUser, session *TelegramSession)

// NewTelegramBot создает нового бота
func NewTelegramBot(token string, db *gorm.DB) (*TelegramBot, error) {
    api, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        return nil, fmt.Errorf("ошибка создания бота: %w", err)
    }
    
    api.Debug = false
    
    bot := &TelegramBot{
        api:      api,
        db:       db,
        handlers: make(map[string]CommandHandler),
    }
    
    bot.registerHandlers()
    
    log.Printf("Бот @%s успешно запущен", api.Self.UserName)
    
    return bot, nil
}

// регистрация обработчиков команд
func (b *TelegramBot) registerHandlers() {
    // Команды
    b.handlers["start"] = b.handleStart
    b.handlers["help"] = b.handleHelp
    b.handlers["login"] = b.handleLogin
    b.handlers["medications"] = b.handleMedications
    b.handlers["appointments"] = b.handleAppointments
    b.handlers["symptoms"] = b.handleSymptoms
    b.handlers["today"] = b.handleToday
    b.handlers["settings"] = b.handleSettings
    
    // Callback-обработчики
    b.handlers["menu_medications"] = b.handleMenuMedications
    b.handlers["menu_appointments"] = b.handleMenuAppointments
    b.handlers["menu_symptoms"] = b.handleMenuSymptoms
    b.handlers["menu_analyses"] = b.handleMenuAnalyses
    b.handlers["menu_doctors"] = b.handleMenuDoctors
    b.handlers["menu_settings"] = b.handleMenuSettings
    
    b.handlers["med_add"] = b.handleMedAdd
    b.handlers["med_list"] = b.handleMedList
    b.handlers["med_take"] = b.handleMedTake
    b.handlers["med_stats"] = b.handleMedStats
    
    b.handlers["symptom_"] = b.handleSymptomSelect
    b.handlers["symptom_custom"] = b.handleSymptomCustom
    b.handlers["intensity_"] = b.handleIntensity
    
    b.handlers["settings_link"] = b.handleSettingsLink
    b.handlers["settings_notifications"] = b.handleSettingsNotifications
    b.handlers["settings_profile"] = b.handleSettingsProfile
    
    b.handlers["back_main"] = b.handleBackMain
    b.handlers["confirm_yes_"] = b.handleConfirmYes
    b.handlers["confirm_no_"] = b.handleConfirmNo
    b.handlers["time_"] = b.handleTimeSelect
}

// Запуск бота (webhook режим)
func (b *TelegramBot) StartWebhook(webhookURL string) error {
    webhookConfig, err := tgbotapi.NewWebhook(webhookURL)
    if err != nil {
        return err
    }
    
    _, err = b.api.Request(webhookConfig)
    if err != nil {
        return err
    }
    
    return nil
}

// Запуск бота (polling режим для разработки)
func (b *TelegramBot) StartPolling() error {
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    
    updates := b.api.GetUpdatesChan(u)
    
    for update := range updates {
        b.processUpdate(update)
    }
    
    return nil
}

// Обработка входящего обновления
func (b *TelegramBot) processUpdate(update tgbotapi.Update) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Паника при обработке: %v", r)
        }
    }()
    
    telegramID := int64(0)
    
    // Определяем ID пользователя
    if update.Message != nil {
        telegramID = update.Message.From.ID
    } else if update.CallbackQuery != nil {
        telegramID = update.CallbackQuery.From.ID
    }
    
    if telegramID == 0 {
        return
    }
    
    // Получаем или создаем пользователя
    user, err := b.getOrCreateUser(update)
    if err != nil {
        log.Printf("Ошибка получения пользователя: %v", err)
        return
    }
    
    // Получаем или создаем сессию
    session, err := b.getOrCreateSession(telegramID)
    if err != nil {
        log.Printf("Ошибка получения сессии: %v", err)
        return
    }
    
    // Обрабатываем сообщение или callback
    if update.Message != nil {
        b.handleMessage(update, user, session)
    } else if update.CallbackQuery != nil {
        b.handleCallback(update, user, session)
    }
}

// Обработка текстовых сообщений
func (b *TelegramBot) handleMessage(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    msg := update.Message
    text := strings.TrimSpace(msg.Text)
    
    // Обработка команд
    if strings.HasPrefix(text, "/") {
        command := strings.ToLower(strings.TrimPrefix(text, "/"))
        if handler, ok := b.handlers[command]; ok {
            handler(update, user, session)
        } else {
            b.sendMessage(msg.Chat.ID, "Неизвестная команда. Напишите /help", nil)
        }
        return
    }
    
    // Обработка в зависимости от состояния сессии
    switch session.State {
    case StateAwaitingAuth:
        b.handleAuthInput(update, user, session, text)
    case StateAwaitingSymptoms:
        b.handleSymptomsInput(update, user, session, text)
    case StateAwaitingMedication:
        b.handleMedicationInput(update, user, session, text)
    case StateAwaitingDosage:
        b.handleDosageInput(update, user, session, text)
    case StateAwaitingFrequency:
        b.handleFrequencyInput(update, user, session, text)
    case StateAwaitingDuration:
        b.handleDurationInput(update, user, session, text)
    default:
        // По умолчанию - показываем меню
        b.showMainMenu(msg.Chat.ID)
    }
}

// Обработка callback-запросов
func (b *TelegramBot) handleCallback(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    callback := update.CallbackQuery
    data := callback.Data
    
    // Отвечаем на callback, чтобы убрать часики
    b.api.Send(tgbotapi.NewCallback(callback.ID, ""))
    
    // Ищем обработчик
    for prefix, handler := range b.handlers {
        if strings.HasPrefix(data, prefix) {
            handler(update, user, session)
            return
        }
    }
    
    // Если не нашли, показываем меню
    b.showMainMenu(callback.Message.Chat.ID)
}
