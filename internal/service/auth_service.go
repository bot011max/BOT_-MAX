package service

import (
    "errors"
    "time"
    "os"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "github.com/bot011max/medical-bot/internal/models"
    "github.com/bot011max/medical-bot/internal/repository"
)

type AuthService struct {
    userRepo  *repository.UserRepository
    jwtSecret []byte
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "default-secret-key-change-in-production"
    }
    return &AuthService{
        userRepo:  userRepo,
        jwtSecret: []byte(secret),
    }
}

type RegisterRequest struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=8"`
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name" binding:"required"`
    Phone     string `json:"phone"`
}

type RegisterResponse struct {
    ID        uuid.UUID `json:"id"`
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}

func (s *AuthService) Register(req *RegisterRequest) (*RegisterResponse, error) {
    existing, _ := s.userRepo.FindByEmail(req.Email)
    if existing != nil {
        return nil, errors.New("email already registered")
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, errors.New("failed to hash password")
    }

    user := &models.User{
        ID:           uuid.New(),
        Email:        req.Email,
        PasswordHash: string(hash),
        FirstName:    req.FirstName,
        LastName:     req.LastName,
        Phone:        req.Phone,
        Role:         "patient",
        IsActive:     true,
    }

    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }

    return &RegisterResponse{
        ID:        user.ID,
        Email:     user.Email,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        Role:      user.Role,
        CreatedAt: user.CreatedAt,
    }, nil
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Token string           `json:"token"`
    User  *RegisterResponse `json:"user"`
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
    user, err := s.userRepo.FindByEmail(req.Email)
    if err != nil || user == nil {
        return nil, errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return nil, errors.New("invalid credentials")
    }

    if !user.IsActive {
        return nil, errors.New("account is disabled")
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID.String(),
        "role":    user.Role,
        "email":   user.Email,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
        "iat":     time.Now().Unix(),
    })

    tokenString, err := token.SignedString(s.jwtSecret)
    if err != nil {
        return nil, errors.New("failed to generate token")
    }

    user.UpdatedAt = time.Now()
    s.userRepo.Update(user)

    return &LoginResponse{
        Token: tokenString,
        User: &RegisterResponse{
            ID:        user.ID,
            Email:     user.Email,
            FirstName: user.FirstName,
            LastName:  user.LastName,
            Role:      user.Role,
            CreatedAt: user.CreatedAt,
        },
    }, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return s.jwtSecret, nil
    })
}
