package services

import (
	"booktracker/backend/internal/apperrors"
	"booktracker/backend/internal/models"
	"booktracker/backend/internal/repositories"
	"errors"
	"fmt"
)

type BookService struct {
	bookRepo    *repositories.BookRepository
	sessionRepo *repositories.SessionRepository
}

func NewBookService(bookRepo *repositories.BookRepository, sessionRepo *repositories.SessionRepository) *BookService {
	return &BookService{
		bookRepo:    bookRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *BookService) GetAllBooks(userID string) ([]models.Book, error) {
	return s.bookRepo.GetAllByUserID(userID)
}

func (s *BookService) GetBook(bookID, userID string) (*models.Book, error) {
	book, err := s.bookRepo.GetByID(bookID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, fmt.Errorf("%w: book not found", apperrors.ErrNotFound)
		}
		return nil, err
	}
	return book, nil
}

func (s *BookService) CreateBook(userID string, req models.CreateBookRequest) (*models.Book, error) {
	if req.TotalPages < 0 {
		return nil, fmt.Errorf("%w: total pages must be greater than 0", apperrors.ErrValidation)
	}

	return s.bookRepo.Create(userID, req)
}

func (s *BookService) UpdateBook(bookID, userID string, req models.UpdateBookRequest) (*models.Book, error) {
	return s.bookRepo.Update(bookID, userID, req)
}

func (s *BookService) DeleteBook(bookID, userID string) error {
	hasActiveSession, err := s.sessionRepo.HasActiveSession(bookID, userID)
	if err != nil {
		return err
	}

	if hasActiveSession {
		return fmt.Errorf("%w: cannot delete a book with an active reading session", apperrors.ErrConflict)
	}

	return s.bookRepo.Delete(bookID, userID)
}
