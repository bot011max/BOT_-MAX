package telegram

import (
    "fmt"
    "math/rand"
    "time"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// /start
func (b *TelegramBot) handleStart(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    chatID := update.Message.Chat.ID
    
    welcome := fmt.Sprintf("👋 Здравствуйте, %s!\n\n", update.Message.From.FirstName)
    welcome += "Я ваш персональный медицинский помощник.\n\n"
    welcome += "Используйте /help для списка команд"
    
    b.sendMessage(chatID, welcome, MainMenu())
}

// /help
func (b *TelegramBot) handleHelp(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    chatID := update.Message.Chat.ID
    
    help := "📋 Доступные команды:\n\n"
    help += "/start - Начать работу\n"
    help += "/help - Показать эту справку\n"
    help += "/login - Войти в аккаунт\n"
    help += "/medications - Мои лекарства\n"
    help += "/appointments - Мои визиты\n"
    help += "/symptoms - Записать симптомы\n"
    help += "/today - Расписание на сегодня\n"
    help += "/settings - Настройки\n\n"
    
    if user.UserID != "" {
        help += "✅ Вы авторизованы"
    } else {
        help += "❌ Вы не авторизованы. Используйте /login"
    }
    
    b.sendMessage(chatID, help, nil)
}

// /login
func (b *TelegramBot) handleLogin(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    chatID := update.Message.Chat.ID
    
    if user.UserID != "" {
        b.sendMessage(chatID, "✅ Вы уже авторизованы!", nil)
        return
    }
    
    // Генерируем одноразовый код для авторизации
    code := fmt.Sprintf("%06d", rand.Intn(1000000))
    expires := time.Now().Add(10 * time.Minute)
    
    user.AuthToken = code
    user.TokenExpires = &expires
    b.db.Save(user)
    
    session.State = StateAwaitingAuth
    b.db.Save(session)
    
    msg := "🔐 Для авторизации введите код подтверждения.\n\n"
    msg += fmt.Sprintf("Код: <b>%s</b>\n\n", code)
    msg += "Код действителен 10 минут"
    
    b.sendMessage(chatID, msg, nil)
}

// /medications
func (b *TelegramBot) handleMedications(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    chatID := update.Message.Chat.ID
    
    if user.UserID == "" {
        b.requireAuth(chatID)
        return
    }
    
    msg := "💊 Ваши лекарства:\n\n"
    msg += "Функция в разработке"
    
    b.sendMessage(chatID, msg, MedicationsMenu(false))
}

// /appointments
func (b *TelegramBot) handleAppointments(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    chatID := update.Message.Chat.ID
    
    if user.UserID == "" {
        b.requireAuth(chatID)
        return
    }
    
    msg := "📅 Ваши визиты:\n\n"
    msg += "Функция в разработке"
    
    b.sendMessage(chatID, msg, nil)
}

// /symptoms
func (b *TelegramBot) handleSymptoms(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    chatID := update.Message.Chat.ID
    
    if user.UserID == "" {
        b.requireAuth(chatID)
        return
    }
    
    msg := "📝 Что вас беспокоит?\n\nВыберите симптом:"
    b.sendMessage(chatID, msg, SymptomsMenu())
}

// /today
func (b *TelegramBot) handleToday(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    chatID := update.Message.Chat.ID
    
    if user.UserID == "" {
        b.requireAuth(chatID)
        return
    }
    
    now := time.Now()
    today := now.Format("02.01.2006")
    
    msg := fmt.Sprintf("📅 <b>Расписание на %s</b>\n\n", today)
    msg += "Функция в разработке"
    
    b.sendMessage(chatID, msg, nil)
}

// /settings
func (b *TelegramBot) handleSettings(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    chatID := update.Message.Chat.ID
    
    msg := "⚙️ Настройки\n\n"
    if user.UserID != "" {
        msg += "✅ Аккаунт привязан\n"
    } else {
        msg += "❌ Аккаунт не привязан\n"
    }
    
    b.sendMessage(chatID, msg, SettingsMenu())
}
