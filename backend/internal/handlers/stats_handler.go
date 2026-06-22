package handlers

import (
	"booktracker/backend/internal/services"
	"net/http"
	"strconv"

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
	stats, err := h.statsService.GetStats(c.Request.Context(), userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *StatsHandler) GetReadingHistory(c *gin.Context) {
	userID := c.GetString("userID")
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || (days != 7 && days != 30 && days != 90) {
        c.JSON(http.StatusBadRequest, gin.H{"message": "days must be 7, 30, or 90"})
        return
    }
	history, err := h.statsService.GetReadingHistory(c.Request.Context(), userID, days)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": history})
}