package api

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/bot011max/medical-bot/internal/models"
    "github.com/bot011max/medical-bot/internal/repository"
)

type MedicationHandler struct {
    repo *repository.MedicationRepository
}

func NewMedicationHandler(repo *repository.MedicationRepository) *MedicationHandler {
    return &MedicationHandler{repo: repo}
}

type CreateMedicationRequest struct {
    Name         string     `json:"name" binding:"required"`
    Dosage       string     `json:"dosage"`
    Frequency    string     `json:"frequency"`
    Instructions string     `json:"instructions"`
    StartDate    *time.Time `json:"start_date"`
    EndDate      *time.Time `json:"end_date"`
}

func (h *MedicationHandler) Create(c *gin.Context) {
    var req CreateMedicationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    uid := userID.(string) // Теперь это строка, а не uuid.UUID

    startDate := time.Now()
    if req.StartDate != nil {
        startDate = *req.StartDate
    }

    medication := &models.Medication{
        UserID:       uid,
        Name:         req.Name,
        Dosage:       req.Dosage,
        Frequency:    req.Frequency,
        Instructions: req.Instructions,
        StartDate:    startDate,
        EndDate:      req.EndDate,
        IsActive:     true,
    }

    if err := h.repo.Create(medication); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "data":    medication,
    })
}

func (h *MedicationHandler) List(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    uid := userID.(string)

    medications, err := h.repo.FindByUserID(uid)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    medications,
    })
}

func (h *MedicationHandler) Get(c *gin.Context) {
    id := c.Param("id")
    medication, err := h.repo.FindByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Medication not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    medication,
    })
}

func (h *MedicationHandler) Update(c *gin.Context) {
    id := c.Param("id")

    var req CreateMedicationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    medication, err := h.repo.FindByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Medication not found"})
        return
    }

    if req.Name != "" {
        medication.Name = req.Name
    }
    if req.Dosage != "" {
        medication.Dosage = req.Dosage
    }
    if req.Frequency != "" {
        medication.Frequency = req.Frequency
    }
    if req.Instructions != "" {
        medication.Instructions = req.Instructions
    }
    if req.EndDate != nil {
        medication.EndDate = req.EndDate
    }

    if err := h.repo.Update(medication); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    medication,
    })
}

func (h *MedicationHandler) Delete(c *gin.Context) {
    id := c.Param("id")

    if err := h.repo.Delete(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Medication deleted successfully",
    })
}
