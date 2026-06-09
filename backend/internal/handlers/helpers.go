package handlers

import (
	"booktracker/backend/internal/apperrors"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrors.ErrValidation):
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	case errors.Is(err, apperrors.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	case errors.Is(err, apperrors.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "an unexpected error occurred"})
	}
}
