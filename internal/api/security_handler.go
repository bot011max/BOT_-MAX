package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/bot011max/medical-bot/internal/security"
    "github.com/bot011max/medical-bot/internal/ocr"
)

type SecurityHandler struct {
    hsm      *security.HardwareSecurityModule
    selfHeal *security.SelfHealingManager
    ocr      *ocr.OCRProcessor
}

func NewSecurityHandler() *SecurityHandler {
    return &SecurityHandler{
        hsm:      security.NewHardwareSecurityModule(),
        selfHeal: security.NewSelfHealingManager("/tmp/medical-bot-backups"),
        ocr:      ocr.NewOCRProcessor(),
    }
}

// GetHSMInfo возвращает информацию о HSM
func (h *SecurityHandler) GetHSMInfo(c *gin.Context) {
    info := h.hsm.GetHSMInfo()
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    info,
    })
}

// CreateBackup создает новый бэкап
func (h *SecurityHandler) CreateBackup(c *gin.Context) {
    var req struct {
        Description string `json:"description"`
    }
    if err := c.BindJSON(&req); err != nil {
        req.Description = "Manual backup"
    }
    
    backup, err := h.selfHeal.CreateBackup(req.Description)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    backup,
    })
}

// ListBackups возвращает список бэкапов
func (h *SecurityHandler) ListBackups(c *gin.Context) {
    backups := h.selfHeal.ListBackups()
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    backups,
    })
}

// Rollback восстанавливает из бэкапа
func (h *SecurityHandler) Rollback(c *gin.Context) {
    backupID := c.Param("id")
    
    err := h.selfHeal.Rollback(backupID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "System restored successfully",
    })
}

// ProcessPrescription обрабатывает фото рецепта
func (h *SecurityHandler) ProcessPrescription(c *gin.Context) {
    // Получаем файл из запроса
    file, err := c.FormFile("prescription")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Prescription image required",
        })
        return
    }
    
    // Читаем файл
    src, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "Failed to read file",
        })
        return
    }
    defer src.Close()
    
    imageData := make([]byte, file.Size)
    _, err = src.Read(imageData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "Failed to read image data",
        })
        return
    }
    
    // Обрабатываем рецепт
    prescription, err := h.ocr.ProcessImage(imageData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    prescription,
    })
}
