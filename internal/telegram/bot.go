package telegram

import (
    "log"
    "strings"
    "time"

    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/bot011max/medical-bot/internal/repository"
)

type Bot struct {
    api            *tgbotapi.BotAPI
    userRepo       *repository.UserRepository
    medicationRepo *repository.MedicationRepository
}

func NewBot(token string, userRepo *repository.UserRepository, medicationRepo *repository.MedicationRepository) (*Bot, error) {
    api, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        return nil, err
    }

    log.Printf("✅ Telegram bot @%s started", api.Self.UserName)

    return &Bot{
        api:            api,
        userRepo:       userRepo,
        medicationRepo: medicationRepo,
    }, nil
}

func (b *Bot) Start() {
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := b.api.GetUpdatesChan(u)

    for update := range updates {
        if update.Message != nil {
            b.handleMessage(update.Message)
        }
    }
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
    chatID := msg.Chat.ID
    text := strings.TrimSpace(msg.Text)

    if strings.HasPrefix(text, "/") {
        command := strings.ToLower(strings.TrimPrefix(text, "/"))
        b.handleCommand(chatID, command, msg.From.ID)
        return
    }

    b.sendMessage(chatID, "❓ Неизвестная команда. Используйте /help")
}

func (b *Bot) handleCommand(chatID int64, command string, telegramID int64) {
    switch command {
    case "start":
        b.handleStart(chatID)
    case "help":
        b.handleHelp(chatID)
    case "bind":
        b.handleBind(chatID, telegramID)
    case "medications":
        b.handleMedications(chatID, telegramID)
    case "take":
        b.handleTake(chatID, telegramID)
    default:
        b.sendMessage(chatID, "❓ Неизвестная команда. Используйте /help")
    }
}

func (b *Bot) handleStart(chatID int64) {
    msg := "👋 Добро пожаловать в Медицинского бота!\n\n"
    msg += "🤖 Я помогу вам:\n"
    msg += "• 💊 Отслеживать прием лекарств\n"
    msg += "• 📝 Записывать симптомы\n"
    msg += "• 📅 Напоминать о визитах\n\n"
    msg += "🔐 Используйте /bind для привязки аккаунта\n"
    msg += "📋 Используйте /help для списка команд"

    b.sendMessage(chatID, msg)
}

func (b *Bot) handleHelp(chatID int64) {
    msg := "📋 <b>Доступные команды:</b>\n\n"
    msg += "/start - начать работу\n"
    msg += "/help - показать справку\n"
    msg += "/bind - привязать аккаунт\n"
    msg += "/medications - мои лекарства\n"
    msg += "/take - отметить прием лекарства\n\n"
    msg += "💡 Для привязки аккаунта:\n"
    msg += "1. Войдите в веб-версию\n"
    msg += "2. В разделе 'Настройки' выберите 'Привязать Telegram'\n"
    msg += "3. Введите код, который пришлет бот"

    b.sendMessage(chatID, msg)
}

func (b *Bot) handleBind(chatID int64, telegramID int64) {
    code := generateCode(6)

    msg := "🔐 <b>Привязка Telegram к аккаунту</b>\n\n"
    msg += "Ваш код: <code>" + code + "</code>\n\n"
    msg += "Введите этот код в веб-версии в разделе 'Привязать Telegram'\n"
    msg += "⏰ Код действителен 10 минут"

    b.sendMessage(chatID, msg)

    log.Printf("Код привязки для %d: %s", telegramID, code)
}

func (b *Bot) handleMedications(chatID int64, telegramID int64) {
    msg := "💊 <b>Ваши лекарства</b>\n\n"
    msg += "Для просмотра лекарств необходимо:\n"
    msg += "1. Привязать аккаунт через /bind\n"
    msg += "2. Войти в веб-версию\n\n"
    msg += "После привязки здесь появится список ваших лекарств"

    b.sendMessage(chatID, msg)
}

func (b *Bot) handleTake(chatID int64, telegramID int64) {
    b.sendMessage(chatID, "✅ Прием лекарства отмечен! Спасибо, что заботитесь о здоровье.")
}

func (b *Bot) sendMessage(chatID int64, text string) {
    msg := tgbotapi.NewMessage(chatID, text)
    msg.ParseMode = "HTML"
    if _, err := b.api.Send(msg); err != nil {
        log.Printf("Ошибка отправки сообщения: %v", err)
    }
}

func generateCode(length int) string {
    const digits = "0123456789"
    code := make([]byte, length)
    for i := range code {
        code[i] = digits[time.Now().UnixNano()%int64(len(digits))]
        time.Sleep(1 * time.Nanosecond)
    }
    return string(code)
}
