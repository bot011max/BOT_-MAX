package service

import (
    "github.com/google/uuid"
    "github.com/bot011max/medical-bot/internal/models"
    "github.com/bot011max/medical-bot/internal/repository"
)

type UserService struct {
    repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
    userID, err := uuid.Parse(id)
    if err != nil {
        return nil, err
    }
    return s.repo.FindByID(userID)
}
