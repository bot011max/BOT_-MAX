package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func MainMenu() tgbotapi.InlineKeyboardMarkup {
    buttons := [][]tgbotapi.InlineKeyboardButton{
        {
            tgbotapi.NewInlineKeyboardButtonData("💊 Лекарства", "medications"),
            tgbotapi.NewInlineKeyboardButtonData("📅 Приемы", "appointments"),
        },
        {
            tgbotapi.NewInlineKeyboardButtonData("📝 Симптомы", "symptoms"),
            tgbotapi.NewInlineKeyboardButtonData("⚙️ Настройки", "settings"),
        },
    }
    return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}
