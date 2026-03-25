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

func (h *AuthHandler) Register(c *gin.Context) {
    var req service.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "invalid request",
            "details": err.Error(),
        })
        return
    }

    resp, err := h.authService.Register(&req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "data":    resp,
    })
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req service.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "invalid request",
            "details": err.Error(),
        })
        return
    }

    resp, err := h.authService.Login(&req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    resp,
    })
}

func (h *AuthHandler) Profile(c *gin.Context) {
    userID := c.GetString("user_id")

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": gin.H{
            "user_id": userID,
        },
    })
}
