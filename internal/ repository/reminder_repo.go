package repository

import (
    "time"
    "errors"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "github.com/bot011max/medical-bot/internal/models"
)

type ReminderRepository struct {
    db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) *ReminderRepository {
    return &ReminderRepository{db: db}
}

func (r *ReminderRepository) Create(reminder *models.Reminder) error {
    return r.db.Create(reminder).Error
}

func (r *ReminderRepository) FindByID(id uuid.UUID) (*models.Reminder, error) {
    var reminder models.Reminder
    err := r.db.Preload("User").Preload("Medication").
        First(&reminder, "id = ?", id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &reminder, err
}

func (r *ReminderRepository) FindByUserID(userID uuid.UUID) ([]models.Reminder, error) {
    var reminders []models.Reminder
    err := r.db.Where("user_id = ?", userID).
        Order("scheduled_at ASC").
        Find(&reminders).Error
    return reminders, err
}

func (r *ReminderRepository) FindPending(scheduledBefore time.Time) ([]models.Reminder, error) {
    var reminders []models.Reminder
    err := r.db.Where("status = ? AND scheduled_at <= ?", "pending", scheduledBefore).
        Order("scheduled_at ASC").
        Limit(100).
        Find(&reminders).Error
    return reminders, err
}

func (r *ReminderRepository) Update(reminder *models.Reminder) error {
    return r.db.Save(reminder).Error
}

func (r *ReminderRepository) MarkAsSent(id uuid.UUID) error {
    now := time.Now()
    return r.db.Model(&models.Reminder{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "status": "sent",
            "sent_at": now,
        }).Error
}

func (r *ReminderRepository) MarkAsAcknowledged(id uuid.UUID) error {
    now := time.Now()
    return r.db.Model(&models.Reminder{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "status": "acknowledged",
            "acknowledged_at": now,
        }).Error
}

func (r *ReminderRepository) MarkAsFailed(id uuid.UUID, errorMsg string) error {
    return r.db.Model(&models.Reminder{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "status": "failed",
            "error_message": errorMsg,
            "retry_count": gorm.Expr("retry_count + 1"),
        }).Error
}

func (r *ReminderRepository) DeleteExpired() error {
    return r.db.Where("status IN (?) AND scheduled_at < ?", 
        []string{"sent", "acknowledged", "failed"}, 
        time.Now().AddDate(0, 0, -30)).
        Delete(&models.Reminder{}).Error
}
