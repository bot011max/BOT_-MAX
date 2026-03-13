package telegram

import (
    "encoding/json"
    "fmt"
    "math/rand"
    "time"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/google/uuid"
)

// /start
func (b *TelegramBot) handleStart(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    chatID := update.Message.Chat.ID
    
    welcome := fmt.Sprintf("👋 Здравствуйте, %s!\n\n", update.Message.From.FirstName)
    welcome += "Я ваш персональный медицинский помощник. Вот что я умею:\n\n"
    welcome += "💊 Отслеживать приём лекарств\n"
    welcome += "📝 Записывать симптомы\n"
    welcome += "📅 Напоминать о визитах к врачу\n"
    welcome += "📊 Анализировать результаты анализов\n"
    welcome += "👨‍⚕️ Связывать с вашим врачом\n\n"
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
        help += "✅ Вы авторизованы как " + user.FirstName + " " + user.LastName
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
    session.TempData = code
    b.db.Save(session)
    
    msg := "🔐 Для авторизации введите код подтверждения.\n\n"
    msg += "Код отправлен на вашу почту (в демо-режиме показываем):\n\n"
    msg += fmt.Sprintf("📱 Код: <b>%s</b>\n\n", code)
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
    
    // Получаем лекарства из БД
    medications := b.getUserMedications(user.UserID)
    
    if len(medications) == 0 {
        msg := "💊 У вас пока нет назначенных лекарств.\n\n"
        msg += "Нажмите кнопку «➕ Добавить», чтобы добавить новое лекарство."
        b.sendMessage(chatID, msg, MedicationsMenu(false))
        return
    }
    
    msg := "💊 Ваши лекарства:\n\n"
    for i, med := range medications {
        msg += fmt.Sprintf("%d. <b>%s</b> %s\n", i+1, med.Name, med.Dosage)
        msg += fmt.Sprintf("   🕐 %s\n", med.Frequency)
        msg += fmt.Sprintf("   📅 до %s\n\n", med.EndDate.Format("02.01.2006"))
    }
    
    b.sendMessage(chatID, msg, MedicationsMenu(true))
}

// /appointments
func (b *TelegramBot) handleAppointments(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    chatID := update.Message.Chat.ID
    
    if user.UserID == "" {
        b.requireAuth(chatID)
        return
    }
    
    appointments := b.getUserAppointments(user.UserID)
    
    if len(appointments) == 0 {
        msg := "📅 У вас нет предстоящих визитов к врачу."
        b.sendMessage(chatID, msg, nil)
        return
    }
    
    msg := "📅 Предстоящие визиты:\n\n"
    for i, apt := range appointments {
        msg += fmt.Sprintf("%d. <b>%s</b>\n", i+1, apt.Date.Format("02.01.2006 15:04"))
        msg += fmt.Sprintf("   👨‍⚕️ %s\n", apt.DoctorName)
        msg += fmt.Sprintf("   🏥 %s\n\n", apt.Location)
    }
    
    b.sendMessage(chatID, msg, nil)
}

// /symptoms
func (b *TelegramBot) handleSymptoms(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    chatID := update.Message.Chat.ID
    
    if user.UserID == "" {
        b.requireAuth(chatID)
        return
    }
    
    msg := "📝 Что вас беспокоит?\n\nВыберите симптом или напишите свой вариант:"
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
    
    // Лекарства на сегодня
    medications := b.getTodayMedications(user.UserID)
    msg += "💊 Лекарства:\n"
    if len(medications) == 0 {
        msg += "   Нет назначений\n"
    } else {
        for _, med := range medications {
            status := "⏳ ожидает"
            if med.Taken {
                status = "✅ принято"
            }
            msg += fmt.Sprintf("   • %s %s - %s %s\n", med.Name, med.Dosage, med.Time, status)
        }
    }
    
    msg += "\n"
    
    // Приемы на сегодня
    appointments := b.getTodayAppointments(user.UserID)
    msg += "👨‍⚕️ Приемы:\n"
    if len(appointments) == 0 {
        msg += "   Нет приемов\n"
    } else {
        for _, apt := range appointments {
            msg += fmt.Sprintf("   • %s - %s (%s)\n", apt.Time, apt.DoctorName, apt.Specialty)
        }
    }
    
    b.sendMessage(chatID, msg, nil)
}

// /settings
func (b *TelegramBot) handleSettings(update tgbotapi.Update, user *TelegramUser, session *TelegramSession) {
    chatID := update.Message.Chat.ID
    
    msg := "⚙️ Настройки\n\n"
    if user.UserID != "" {
        msg += fmt.Sprintf("👤 Аккаунт: %s %s\n", user.FirstName, user.LastName)
        msg += fmt.Sprintf("📧 Email: %s\n", user.Email)
        msg += fmt.Sprintf("🔔 Уведомления: включены\n")
    } else {
        msg += "❌ Аккаунт не привязан\n"
    }
    
    b.sendMessage(chatID, msg, SettingsMenu())
}
