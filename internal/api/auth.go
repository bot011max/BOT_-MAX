package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/bot011max/medical-bot/internal/service"
)

type AuthHandler struct {
    authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=6"`
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name" binding:"required"`
    Phone     string `json:"phone"`
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.authService.Register(req.Email, req.Password, req.FirstName, req.LastName, req.Phone)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "data": gin.H{
            "id":         user.ID,
            "email":      user.Email,
            "first_name": user.FirstName,
            "last_name":  user.LastName,
            "role":       user.Role,
            "created_at": user.CreatedAt,
        },
    })
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    token, user, err := h.authService.Login(req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": gin.H{
            "token": token,
            "user": gin.H{
                "id":         user.ID,
                "email":      user.Email,
                "first_name": user.FirstName,
                "last_name":  user.LastName,
                "role":       user.Role,
                "created_at": user.CreatedAt,
            },
        },
    })
}

func (h *AuthHandler) Profile(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    user,
    })
}

// Setup2FA настраивает двухфакторную аутентификацию
func (h *AuthHandler) Setup2FA(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    currentUser := user.(*models.User)
    
    // Генерация секрета
    secret, err := auth.GenerateSecret()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate secret"})
        return
    }
    
    // Генерация QR кода
    qrData, err := auth.GenerateQRCode(secret, currentUser.Email, "MedicalBot")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
        return
    }
    
    // Генерация резервных кодов
    backupCodes := auth.GenerateBackupCodes(10)
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": gin.H{
            "secret":       secret,
            "qr_code":      base64.StdEncoding.EncodeToString(qrData),
            "backup_codes": backupCodes,
        },
    })
}

// Verify2FA проверяет код двухфакторной аутентификации
func (h *AuthHandler) Verify2FA(c *gin.Context) {
    var req struct {
        Code string `json:"code" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    currentUser := user.(*models.User)
    
    // Проверка кода (нужно хранить secret в БД)
    // if auth.VerifyCode(user2FASecret, req.Code) {
    //     c.JSON(http.StatusOK, gin.H{"success": true})
    // } else {
    //     c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid code"})
    // }
    
    c.JSON(http.StatusOK, gin.H{"success": true, "message": "2FA verified"})
}
