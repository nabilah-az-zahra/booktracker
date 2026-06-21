package handlers

import (
	"booktracker/backend/internal/models"
	"booktracker/backend/internal/services"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	response, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	response, refreshToken, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	setRefreshCookie(c, refreshToken)
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "refresh token required"})
		return
	}

	response, newRefreshToken, err := h.authService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		clearRefreshCookie(c)
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	setRefreshCookie(c, newRefreshToken)
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	accessToken := ""
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		accessToken = authHeader[7:]
	}

	refreshToken, _ := c.Cookie("refresh_token")
	h.authService.Logout(c.Request.Context(), accessToken, refreshToken)
	clearRefreshCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("userID")
	user, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) UpdateGoal(c *gin.Context) {
	userID := c.GetString("userID")
	var req models.UpdateGoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.authService.UpdateGoal(c.Request.Context(), userID, req.YearlyGoal); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "goal updated successfully"})
}

func (h *AuthHandler) UpdateTimezone(ctx *gin.Context) {
    userID := ctx.GetString("userID")
    var req models.UpdateTimezoneRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.authService.UpdateTimezone(ctx, userID, req.Timezone); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid timezone"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "timezone updated"})
}

func setRefreshCookie(c *gin.Context, token string) {
	secure := os.Getenv("GO_ENV") == "production"
	c.SetCookie(
		"refresh_token",
		token,
		int(30*24*time.Hour/time.Second),
		"/",
		"",
		secure,
		true,
	)
}

func clearRefreshCookie(c *gin.Context) {
	secure := os.Getenv("GO_ENV") == "production"
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		secure,
		true,
	)
}