package service

import (
    "errors"
    "github.com/bot011max/medical-bot/internal/models"
    "github.com/bot011max/medical-bot/internal/repository"
    "github.com/bot011max/medical-bot/internal/security"
    "golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
    return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(email, password, firstName, lastName, phone string) (*models.User, error) {
    existing, _ := s.userRepo.FindByEmail(email)
    if existing != nil {
        return nil, errors.New("user already exists")
    }
    
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    user := &models.User{
        Email:        email,
        PasswordHash: string(hashedPassword),
        FirstName:    firstName,
        LastName:     lastName,
        Phone:        phone,
        Role:         "patient",
        IsActive:     true,
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    
    return user, nil
}

func (s *AuthService) Login(email, password string) (string, *models.User, error) {
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return "", nil, errors.New("invalid credentials")
    }
    
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return "", nil, errors.New("invalid credentials")
    }
    
    // ID уже строка, не нужно вызывать .String()
    token, err := security.GenerateJWT(user.ID, user.Email, user.Role)
    if err != nil {
        return "", nil, err
    }
    
    return token, user, nil
}
