package services

import (
	"booktracker/backend/internal/apperrors"
	"booktracker/backend/internal/models"
	"booktracker/backend/internal/repositories"
	"context"
	"fmt"
	"log"
)

type SessionService struct {
	sessionRepo *repositories.SessionRepository
	bookRepo    *repositories.BookRepository
}

func NewSessionService(
	sessionRepo *repositories.SessionRepository,
	bookRepo *repositories.BookRepository,
) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
		bookRepo:    bookRepo,
	}
}

func (s *SessionService) StartSession(ctx context.Context, bookID, userID string) (*models.ReadingSession, error) {
	existing, err := s.sessionRepo.GetActiveSession(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: you already have an active reading session", apperrors.ErrConflict)
	}

	_, err = s.bookRepo.GetByID(ctx, bookID, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: book not found", apperrors.ErrNotFound)
	}

	return s.sessionRepo.StartSession(ctx, bookID, userID)
}

func (s *SessionService) PauseSession(ctx context.Context, sessionID, userID string, duration int) (*models.ReadingSession, error) {
	return s.sessionRepo.PauseSession(ctx, sessionID, userID, duration)
}

func (s *SessionService) ResumeSession(ctx context.Context, sessionID, userID string) (*models.ReadingSession, error) {
	return s.sessionRepo.ResumeSession(ctx, sessionID, userID)
}

func (s *SessionService) StopSession(ctx context.Context, sessionID, userID string, duration, pagesRead int) (*models.ReadingSession, error) {
	if pagesRead < 0 {
		return nil, fmt.Errorf("%w: pages read cannot be negative", apperrors.ErrValidation)
	}
	if duration < 0 {
		return nil, fmt.Errorf("%w: duration cannot be negative", apperrors.ErrValidation)
	}

	session, err := s.sessionRepo.StopSession(ctx, sessionID, userID, duration, pagesRead)
	if err != nil {
		return nil, err
	}

	progress, _ := s.sessionRepo.GetProgress(ctx, session.BookID, userID)
	newPage := pagesRead
	if progress != nil {
		newPage = progress.CurrentPage + pagesRead
	}
	book, err := s.bookRepo.GetByID(ctx, session.BookID, userID)
	if err == nil && book != nil && book.TotalPages > 0 && newPage > book.TotalPages {
		newPage = book.TotalPages
	}

	if err := s.sessionRepo.UpdateProgress(ctx, session.BookID, userID, newPage); err != nil {
		log.Printf("Warning: failed to update progress for book %s: %v", session.BookID, err)
	}

	return session, nil
}

func (s *SessionService) CancelSession(ctx context.Context, sessionID, userID string) error {
	return s.sessionRepo.CancelSession(ctx, sessionID, userID)
}

func (s *SessionService) GetActiveSession(ctx context.Context, userID string) (*models.ReadingSession, error) {
	return s.sessionRepo.GetActiveSession(ctx, userID)
}

func (s *SessionService) GetSessionsByBook(ctx context.Context, bookID, userID string) ([]models.ReadingSession, error) {
	return s.sessionRepo.GetSessionsByBookID(ctx, bookID, userID)
}

func (s *SessionService) GetProgress(ctx context.Context, bookID, userID string) (*models.ReadingProgress, error) {
	return s.sessionRepo.GetProgress(ctx, bookID, userID)
}
