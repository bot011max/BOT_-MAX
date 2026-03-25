package repository

import (
    "errors"
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "github.com/bot011max/medical-bot/internal/models"
)

type MedicationRepository struct {
    db *gorm.DB
}

func NewMedicationRepository(db *gorm.DB) *MedicationRepository {
    return &MedicationRepository{db: db}
}

func (r *MedicationRepository) Create(medication *models.Medication) error {
    return r.db.Create(medication).Error
}

func (r *MedicationRepository) FindByID(id uuid.UUID) (*models.Medication, error) {
    var medication models.Medication
    err := r.db.Preload("User").First(&medication, "id = ?", id).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    return &medication, err
}

func (r *MedicationRepository) FindByUserID(userID uuid.UUID) ([]models.Medication, error) {
    var medications []models.Medication
    err := r.db.Where("user_id = ? AND is_active = ?", userID, true).
        Order("created_at DESC").
        Find(&medications).Error
    return medications, err
}

func (r *MedicationRepository) FindActiveByUserID(userID uuid.UUID) ([]models.Medication, error) {
    var medications []models.Medication
    err := r.db.Where("user_id = ? AND is_active = ? AND (end_date IS NULL OR end_date > ?)", 
        userID, true, time.Now()).
        Order("created_at DESC").
        Find(&medications).Error
    return medications, err
}

func (r *MedicationRepository) Update(medication *models.Medication) error {
    return r.db.Save(medication).Error
}

func (r *MedicationRepository) Delete(id uuid.UUID) error {
    return r.db.Delete(&models.Medication{}, "id = ?", id).Error
}

func (r *MedicationRepository) CountByUserID(userID uuid.UUID) (int64, error) {
    var count int64
    err := r.db.Model(&models.Medication{}).
        Where("user_id = ?", userID).
        Count(&count).Error
    return count, err
}
