package handlers

import (
	"booktracker/backend/internal/models"
	"booktracker/backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookService *services.BookService
}

func NewBookHandler(bookService *services.BookService) *BookHandler {
	return &BookHandler{bookService: bookService}
}

func (h *BookHandler) GetAllBooks(c *gin.Context) {
	userID := c.GetString("userID")
	books, err := h.bookService.GetAllBooks(c.Request.Context(), userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": books})
}

func (h *BookHandler) GetBook(c *gin.Context) {
	userID := c.GetString("userID")
	bookID := c.Param("id")

	book, err := h.bookService.GetBook(c.Request.Context(), bookID, userID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	userID := c.GetString("userID")
	var req models.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	book, err := h.bookService.CreateBook(c.Request.Context(), userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": book})
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	userID := c.GetString("userID")
	bookID := c.Param("id")
	var req models.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	book, err := h.bookService.UpdateBook(c.Request.Context(), bookID, userID, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	userID := c.GetString("userID")
	bookID := c.Param("id")

	if err := h.bookService.DeleteBook(c.Request.Context(), bookID, userID); err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
}
