package handlers

import (
	"booktracker/backend/internal/models"
	"booktracker/backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionService *services.SessionService
}

func NewSessionHandler(sessionService *services.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

func (h *SessionHandler) StartSession(c *gin.Context) {
	userID := c.GetString("userID")
	var req models.StartSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	session, err := h.sessionService.StartSession(c.Request.Context(), req.BookID, userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": session})
}

func (h *SessionHandler) PauseSession(c *gin.Context) {
	userID := c.GetString("userID")
	sessionID := c.Param("id")

	var body struct {
		Duration int `json:"duration_seconds"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if body.Duration <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "duration must be greater than 0"})
		return
	}

	session, err := h.sessionService.PauseSession(c.Request.Context(), sessionID, userID, body.Duration)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": session})
}

func (h *SessionHandler) ResumeSession(c *gin.Context) {
	userID := c.GetString("userID")
	sessionID := c.Param("id")

	session, err := h.sessionService.ResumeSession(c.Request.Context(), sessionID, userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": session})
}

func (h *SessionHandler) StopSession(c *gin.Context) {
	userID := c.GetString("userID")
	sessionID := c.Param("id")

	var body struct {
		Duration  int `json:"duration_seconds"`
		PagesRead int `json:"pages_read" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	session, err := h.sessionService.StopSession(c.Request.Context(), sessionID, userID, body.Duration, body.PagesRead)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": session})
}

func (h *SessionHandler) CancelSession(c *gin.Context) {
	userID := c.GetString("userID")
	sessionID := c.Param("id")

	if err := h.sessionService.CancelSession(c.Request.Context(), sessionID, userID); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "session cancelled"})
}

func (h *SessionHandler) GetActiveSession(c *gin.Context) {
	userID := c.GetString("userID")
	session, err := h.sessionService.GetActiveSession(c.Request.Context(), userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": session})
}

func (h *SessionHandler) GetSessionsByBook(c *gin.Context) {
	userID := c.GetString("userID")
	bookID := c.Param("bookId")

	sessions, err := h.sessionService.GetSessionsByBook(c.Request.Context(), bookID, userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sessions})
}

func (h *SessionHandler) GetProgress(c *gin.Context) {
	userID := c.GetString("userID")
	bookID := c.Param("bookId")

	progress, err := h.sessionService.GetProgress(c.Request.Context(), bookID, userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": progress})
}

func (h *SessionHandler) AddNote(c *gin.Context) {
	userID := c.GetString("userID")
	sessionID := c.Param("id")

	var req models.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	note, err := h.sessionService.AddNote(c.Request.Context(), sessionID, userID, req)
	if err != nil {
		handleServiceError(c, err)
		return 
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": note})
}

func (h *SessionHandler) GetNotes(c *gin.Context) {
	userID := c.GetString("userID")
	sessionID := c.Param("id")

	notes, err := h.sessionService.GetNotes(c.Request.Context(), sessionID, userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": notes})
}

func (h *SessionHandler) DeleteNote(c *gin.Context) {
	userID := c.GetString("userID")
	noteID := c.Param("noteId")

	if err := h.sessionService.DeleteNote(c.Request.Context(), noteID, userID); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "note deleted"})
}