package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/bot011max/medical-bot/internal/ai"
)

type AIHandler struct {
    analyzer *ai.SymptomAnalyzer
}

func NewAIHandler() *AIHandler {
    return &AIHandler{
        analyzer: ai.NewSymptomAnalyzer(),
    }
}

func (h *AIHandler) AnalyzeSymptoms(c *gin.Context) {
    var req struct {
        Symptoms []string `json:"symptoms"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    results := h.analyzer.Analyze(req.Symptoms)
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    results,
    })
}

// Добавляем в main.go маршруты
// router.POST("/api/ai/analyze", aiHandler.AnalyzeSymptoms)
// router.GET("/api/doctor/dashboard", doctorHandler.GetDashboard)
// router.POST("/api/billing/upgrade", billingHandler.UpgradeToPremium)
