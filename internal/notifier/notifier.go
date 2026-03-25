package notifier

import "github.com/bot011max/medical-bot/internal/models"

type MessageNotifier interface {
    SendReminder(reminder *models.Reminder) error
    SendMessage(userID string, message string) error
}
