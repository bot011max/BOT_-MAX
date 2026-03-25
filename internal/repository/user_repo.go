package repository

import (
    "github.com/bot011max/medical-bot/internal/models"
    "gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    result := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
    var user models.User
    result := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
    return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id string) error {
    return r.db.Where("id = ?", id).Delete(&models.User{}).Error
}
