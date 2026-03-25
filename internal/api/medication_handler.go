package api

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
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
    Name         string `json:"name" binding:"required"`
    Dosage       string `json:"dosage"`
    Frequency    string `json:"frequency" binding:"required"`
    Instructions string `json:"instructions"`
    StartDate    string `json:"start_date"`
    EndDate      string `json:"end_date"`
}

func (h *MedicationHandler) Create(c *gin.Context) {
    userID := c.GetString("user_id")

    var req CreateMedicationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    med := &models.Medication{
        ID:           uuid.New(),
        UserID:       uuid.MustParse(userID),
        Name:         req.Name,
        Dosage:       req.Dosage,
        Frequency:    req.Frequency,
        Instructions: req.Instructions,
        IsActive:     true,
    }

    if req.StartDate != "" {
        startDate, err := time.Parse("2006-01-02", req.StartDate)
        if err == nil {
            med.StartDate = &startDate
        }
    }

    if req.EndDate != "" {
        endDate, err := time.Parse("2006-01-02", req.EndDate)
        if err == nil {
            med.EndDate = &endDate
        }
    }

    if err := h.repo.Create(med); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "success": true,
        "data":    med,
    })
}

func (h *MedicationHandler) List(c *gin.Context) {
    userID := c.GetString("user_id")

    medications, err := h.repo.FindByUserID(uuid.MustParse(userID))
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
    userID := c.GetString("user_id")

    med, err := h.repo.FindByID(uuid.MustParse(id))
    if err != nil || med == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "medication not found"})
        return
    }

    if med.UserID.String() != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    med,
    })
}

func (h *MedicationHandler) Update(c *gin.Context) {
    id := c.Param("id")
    userID := c.GetString("user_id")

    var req struct {
        Name         string `json:"name"`
        Dosage       string `json:"dosage"`
        Frequency    string `json:"frequency"`
        Instructions string `json:"instructions"`
        IsActive     bool   `json:"is_active"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    med, err := h.repo.FindByID(uuid.MustParse(id))
    if err != nil || med == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "medication not found"})
        return
    }

    if med.UserID.String() != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
        return
    }

    if req.Name != "" {
        med.Name = req.Name
    }
    if req.Dosage != "" {
        med.Dosage = req.Dosage
    }
    if req.Frequency != "" {
        med.Frequency = req.Frequency
    }
    if req.Instructions != "" {
        med.Instructions = req.Instructions
    }
    med.IsActive = req.IsActive

    if err := h.repo.Update(med); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    med,
    })
}

func (h *MedicationHandler) Delete(c *gin.Context) {
    id := c.Param("id")
    userID := c.GetString("user_id")

    med, err := h.repo.FindByID(uuid.MustParse(id))
    if err != nil || med == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "medication not found"})
        return
    }

    if med.UserID.String() != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
        return
    }

    if err := h.repo.Delete(med.ID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "medication deleted",
    })
}
