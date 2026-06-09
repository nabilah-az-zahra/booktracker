package handlers

import (
	"booktracker/backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsService *services.StatsService
}

func NewStatsHandler(statsService *services.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	userID := c.GetString("userID")
	stats, err := h.statsService.GetStats(userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}
