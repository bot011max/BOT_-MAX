package repository

import (
    "github.com/bot011max/medical-bot/internal/models"
    "gorm.io/gorm"
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

func (r *MedicationRepository) FindByUserID(userID string) ([]models.Medication, error) {
    var medications []models.Medication
    err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).
        Order("created_at DESC").
        Find(&medications).Error
    return medications, err
}

func (r *MedicationRepository) FindByID(id string) (*models.Medication, error) {
    var medication models.Medication
    err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&medication).Error
    if err != nil {
        return nil, err
    }
    return &medication, nil
}

func (r *MedicationRepository) Update(medication *models.Medication) error {
    return r.db.Save(medication).Error
}

func (r *MedicationRepository) Delete(id string) error {
    return r.db.Where("id = ?", id).Delete(&models.Medication{}).Error
}
