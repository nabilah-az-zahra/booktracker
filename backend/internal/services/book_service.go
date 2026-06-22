package services

import (
	"booktracker/backend/internal/apperrors"
	"booktracker/backend/internal/models"
	redisclient "booktracker/backend/internal/redis"
	"booktracker/backend/internal/repositories"
	"context"
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

func (s *BookService) GetAllBooks(ctx context.Context, userID string) ([]models.Book, error) {
	return s.bookRepo.GetAllByUserID(ctx, userID)
}

func (s *BookService) GetBook(ctx context.Context, bookID, userID string) (*models.Book, error) {
	book, err := s.bookRepo.GetByID(ctx, bookID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, fmt.Errorf("%w: book not found", apperrors.ErrNotFound)
		}
		return nil, err
	}
	return book, nil
}

func (s *BookService) CreateBook(ctx context.Context, userID string, req models.CreateBookRequest) (*models.Book, error) {
	if req.TotalPages < 0 {
		return nil, fmt.Errorf("%w: total pages must be greater than 0", apperrors.ErrValidation)
	}

	return s.bookRepo.Create(ctx, userID, req)
}

func (s *BookService) UpdateBook(ctx context.Context, bookID, userID string, req models.UpdateBookRequest) (*models.Book, error) {
    book, err := s.bookRepo.Update(ctx, bookID, userID, req)
    if err != nil {
        return nil, err
    }
    redisclient.Client.Del(ctx, "stats:"+userID)
    return book, nil
}

func (s *BookService) DeleteBook(ctx context.Context, bookID, userID string) error {
    hasActiveSession, err := s.sessionRepo.HasActiveSession(ctx, bookID, userID)
    if err != nil {
        return err
    }
    if hasActiveSession {
        return fmt.Errorf("%w: cannot delete a book with an active reading session", apperrors.ErrConflict)
    }
    if err := s.bookRepo.Delete(ctx, bookID, userID); err != nil {
        return err
    }
    redisclient.Client.Del(ctx, "stats:"+userID)
    return nil
}