package api

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/bot011max/medical-bot/internal/service"
    "github.com/bot011max/medical-bot/internal/security"
)

// Handler - структура хендлера
type Handler struct {
    authService *service.AuthService
    audit       *security.AuditLogger
}

// NewHandler - создание хендлера
func NewHandler(authService *service.AuthService, audit *security.AuditLogger) *Handler {
    return &Handler{
        authService: authService,
        audit:       audit,
    }
}

// RegisterRoutes - регистрация всех маршрутов
func (h *Handler) RegisterRoutes(r *gin.Engine) {
    // Публичные маршруты
    r.POST("/api/register", h.Register)
    r.POST("/api/login", h.Login)
    r.POST("/api/refresh", h.RefreshToken)
    
    // Защищенные маршруты
    authorized := r.Group("/api")
    authorized.Use(AuthMiddleware("your-secret-key", h.audit))
    {
        authorized.GET("/profile", h.GetProfile)
        authorized.POST("/logout", h.Logout)
        
        // Маршруты с проверкой роли
        doctorRoutes := authorized.Group("/doctor")
        doctorRoutes.Use(RequireRole("doctor", "clinic"))
        {
            doctorRoutes.GET("/patients", h.GetMyPatients)
            doctorRoutes.POST("/prescriptions", h.CreatePrescription)
        }
        
        patientRoutes := authorized.Group("/patient")
        patientRoutes.Use(RequireRole("patient"))
        {
            patientRoutes.GET("/medications", h.GetMedications)
            patientRoutes.POST("/symptoms", h.AddSymptom)
        }
        
        clinicRoutes := authorized.Group("/clinic")
        clinicRoutes.Use(RequireRole("clinic"))
        {
            clinicRoutes.GET("/all-patients", h.GetAllPatients)
            clinicRoutes.GET("/analytics", h.GetAnalytics)
        }
    }
}

// Register - регистрация
func (h *Handler) Register(c *gin.Context) {
    var req service.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    resp, err := h.authService.Register(&req, c.ClientIP())
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, resp)
}

// Login - вход
func (h *Handler) Login(c *gin.Context) {
    var req service.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    resp, err := h.authService.Login(&req, c.ClientIP(), c.Request.UserAgent())
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, resp)
}

// RefreshToken - обновление токена
func (h *Handler) RefreshToken(c *gin.Context) {
    var req struct {
        RefreshToken string `json:"refresh_token" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    resp, err := h.authService.RefreshToken(req.RefreshToken, c.ClientIP())
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, resp)
}

// GetProfile - получение профиля
func (h *Handler) GetProfile(c *gin.Context) {
    userID := c.GetString("user_id")
    // TODO: получить профиль из БД
    
    c.JSON(http.StatusOK, gin.H{
        "user_id": userID,
        "profile": "user profile data",
    })
}

// Logout - выход
func (h *Handler) Logout(c *gin.Context) {
    userID := c.GetString("user_id")
    
    h.audit.LogAccess("USER_LOGOUT", userID, c.ClientIP(), nil)
    
    c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// GetMyPatients - получение пациентов врача
func (h *Handler) GetMyPatients(c *gin.Context) {
    // TODO: получить пациентов
    c.JSON(http.StatusOK, gin.H{"patients": []string{}})
}

// CreatePrescription - создание назначения
func (h *Handler) CreatePrescription(c *gin.Context) {
    // TODO: создать назначение
    c.JSON(http.StatusCreated, gin.H{"message": "prescription created"})
}

// GetMedications - получение лекарств пациента
func (h *Handler) GetMedications(c *gin.Context) {
    // TODO: получить лекарства
    c.JSON(http.StatusOK, gin.H{"medications": []string{}})
}

// AddSymptom - добавление симптома
func (h *Handler) AddSymptom(c *gin.Context) {
    // TODO: добавить симптом
    c.JSON(http.StatusCreated, gin.H{"message": "symptom added"})
}

// GetAllPatients - получение всех пациентов (для клиники)
func (h *Handler) GetAllPatients(c *gin.Context) {
    // TODO: получить всех пациентов
    c.JSON(http.StatusOK, gin.H{"patients": []string{}})
}

// GetAnalytics - получение аналитики
func (h *Handler) GetAnalytics(c *gin.Context) {
    // TODO: получить аналитику
    c.JSON(http.StatusOK, gin.H{
        "total_patients": 1234,
        "active_doctors": 56,
        "prescriptions":  7890,
    })
}
