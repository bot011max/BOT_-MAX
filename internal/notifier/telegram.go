package notifier

import (
    "log"
    "github.com/bot011max/medical-bot/internal/models"
    "github.com/bot011max/medical-bot/internal/repository"
)

type TelegramNotifier struct {
    bot      interface{}
    userRepo *repository.UserRepository
}

func NewTelegramNotifier(bot interface{}, userRepo *repository.UserRepository) *TelegramNotifier {
    return &TelegramNotifier{
        bot:      bot,
        userRepo: userRepo,
    }
}

func (n *TelegramNotifier) SendReminder(reminder *models.Reminder) error {
    log.Printf("📨 Sending reminder to user %s", reminder.UserID)
    return nil
}

func (n *TelegramNotifier) SendMessage(userID string, message string) error {
    log.Printf("📨 Sending message to user %s: %s", userID, message)
    return nil
}
