package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/bot011max/medical-bot/internal/biometrics"
)

type BiometricsHandler struct {
    voiceBiometrics *biometrics.VoiceBiometricsReal
}

func NewBiometricsHandler() *BiometricsHandler {
    return &BiometricsHandler{
        voiceBiometrics: biometrics.NewVoiceBiometricsReal(),
    }
}

func (h *BiometricsHandler) RegisterVoice(c *gin.Context) {
    userID := c.PostForm("user_id")
    audioFile, err := c.FormFile("audio")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Audio file required"})
        return
    }
    
    // Читаем файл
    file, err := audioFile.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read audio"})
        return
    }
    defer file.Close()
    
    audioData := make([]byte, audioFile.Size)
    file.Read(audioData)
    
    // Извлекаем признаки
    features := h.voiceBiometrics.ExtractFeatures(audioData)
    
    // Регистрируем
    h.voiceBiometrics.Register(userID, features)
    
    c.JSON(http.StatusOK, gin.H{"success": true, "message": "Voice registered"})
}

func (h *BiometricsHandler) VerifyVoice(c *gin.Context) {
    userID := c.PostForm("user_id")
    audioFile, err := c.FormFile("audio")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Audio file required"})
        return
    }
    
    file, err := audioFile.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read audio"})
        return
    }
    defer file.Close()
    
    audioData := make([]byte, audioFile.Size)
    file.Read(audioData)
    
    features := h.voiceBiometrics.ExtractFeatures(audioData)
    verified := h.voiceBiometrics.Verify(userID, features)
    
    c.JSON(http.StatusOK, gin.H{"success": verified, "message": "Voice verified"})
}
