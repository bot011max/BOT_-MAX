package api

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"

    "github.com/bot011max/medical-bot/internal/models"
    "github.com/bot011max/medical-bot/internal/security"
    "github.com/bot011max/medical-bot/pkg/logger"
)

type AuthHandler struct {
    db     *gorm.DB
    armor  *security.AbsoluteArmor
    logger *logger.Logger
}

func NewAuthHandler(db *gorm.DB, armor *security.AbsoluteArmor, logger *logger.Logger) *AuthHandler {
    return &AuthHandler{
        db:     db,
        armor:  armor,
        logger: logger,
    }
}

type RegisterRequest struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=8,max=50"`
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name" binding:"required"`
    Role      string `json:"role" binding:"required,oneof=patient doctor admin"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
    TwoFACode string `json:"twofa_code,omitempty"`
}

type UserResponse struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}

// Register - регистрация нового пользователя
func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Warn("Invalid register request: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    // Проверка сложности пароля
    if !h.armor.ValidatePassword(req.Password) {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Password must contain uppercase, lowercase, number and special character",
        })
        return
    }

    // Хеширование пароля
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        h.logger.Error("Failed to hash password: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    // Генерация 2FA секрета
    twoFASecret := h.armor.GenerateTOTPSecret()

    // Создание пользователя
    user := models.User{
        ID:           uuid.New(),
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
        FirstName:    req.FirstName,
        LastName:     req.LastName,
        Role:         req.Role,
        TwoFASecret:  twoFASecret,
        TwoFAEnabled: false,
    }

    if err := h.db.Create(&user).Error; err != nil {
        h.logger.Error("Failed to create user: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Email already exists"})
        return
    }

    // Аудит
    h.armor.AuditLog("USER_REGISTERED", user.ID.String(), c.ClientIP())

    h.logger.Info("User registered: %s", user.Email)
    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "message": "Registration successful. Please set up 2FA.",
        "twofa_secret": twoFASecret,
    })
}

// Login - вход пользователя
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Warn("Invalid login request: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    // Поиск пользователя
    var user models.User
    if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
        // Защита от timing attack
        bcrypt.CompareHashAndPassword([]byte("fakehash"), []byte(req.Password))
        h.logger.Warn("Failed login attempt for email: %s", req.Email)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Проверка пароля
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        h.logger.Warn("Failed login attempt for user: %s", user.Email)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Проверка 2FA если включена
    if user.TwoFAEnabled {
        if req.TwoFACode == "" {
            c.JSON(http.StatusOK, gin.H{
                "twofa_required": true,
                "message": "2FA code required",
            })
            return
        }
        if !h.armor.VerifyTOTP(user.TwoFASecret, req.TwoFACode) {
            h.logger.Warn("Invalid 2FA code for user: %s", user.Email)
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid 2FA code"})
            return
        }
    }

    // Создание JWT токена
    token, err := h.armor.CreateJWT(user.ID.String(), user.Role)
    if err != nil {
        h.logger.Error("Failed to create JWT: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    // Аудит
    h.armor.AuditLog("USER_LOGIN", user.ID.String(), c.ClientIP())

    h.logger.Info("User logged in: %s", user.Email)
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "token": token,
        "user": UserResponse{
            ID:        user.ID.String(),
            Email:     user.Email,
            FirstName: user.FirstName,
            LastName:  user.LastName,
            Role:      user.Role,
            CreatedAt: user.CreatedAt,
        },
    })
}

// Profile - получение профиля
func (h *AuthHandler) Profile(c *gin.Context) {
    userID := c.GetString("user_id")
    
    var user models.User
    if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
        h.logger.Error("User not found: %s", userID)
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": UserResponse{
            ID:        user.ID.String(),
            Email:     user.Email,
            FirstName: user.FirstName,
            LastName:  user.LastName,
            Role:      user.Role,
            CreatedAt: user.CreatedAt,
        },
    })
}

// SetupRoutes - настройка маршрутов
func SetupRoutes(r *gin.Engine, armor *security.AbsoluteArmor, logger *logger.Logger) {
    // Инициализация БД
    db := initDB()

    authHandler := NewAuthHandler(db, armor, logger)

    api := r.Group("/api")
    {
        // Публичные маршруты
        api.POST("/register", authHandler.Register)
        api.POST("/login", authHandler.Login)

        // Защищенные маршруты
        protected := api.Group("/")
        protected.Use(middleware.AuthRequired(armor))
        {
            protected.GET("/profile", authHandler.Profile)
        }
    }
}

func initDB() *gorm.DB {
    // Подключение к БД
    dsn := "host=postgres user=postgres password=postgres dbname=medical_bot port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    return db
}
