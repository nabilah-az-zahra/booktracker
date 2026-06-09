package services

import (
	"booktracker/backend/internal/apperrors"
	"booktracker/backend/internal/models"
	"booktracker/backend/internal/repositories"
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

func (s *SessionService) StartSession(bookID, userID string) (*models.ReadingSession, error) {
	existing, err := s.sessionRepo.GetActiveSession(userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: you already have an active reading session", apperrors.ErrConflict)
	}

	_, err = s.bookRepo.GetByID(bookID, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: book not found", apperrors.ErrNotFound)
	}

	return s.sessionRepo.StartSession(bookID, userID)
}

func (s *SessionService) PauseSession(sessionID, userID string, duration int) (*models.ReadingSession, error) {
	return s.sessionRepo.PauseSession(sessionID, userID, duration)
}

func (s *SessionService) ResumeSession(sessionID, userID string) (*models.ReadingSession, error) {
	return s.sessionRepo.ResumeSession(sessionID, userID)
}

func (s *SessionService) StopSession(sessionID, userID string, duration, pagesRead int) (*models.ReadingSession, error) {
	if pagesRead < 0 {
		return nil, fmt.Errorf("%w: pages read cannot be negative", apperrors.ErrValidation)
	}
	if duration < 0 {
		return nil, fmt.Errorf("%w: duration cannot be negative", apperrors.ErrValidation)
	}

	session, err := s.sessionRepo.StopSession(sessionID, userID, duration, pagesRead)
	if err != nil {
		return nil, err
	}

	progress, _ := s.sessionRepo.GetProgress(session.BookID, userID)
	newPage := pagesRead
	if progress != nil {
		newPage = progress.CurrentPage + pagesRead
	}

	if err := s.sessionRepo.UpdateProgress(session.BookID, userID, newPage); err != nil {
		log.Printf("Warning: failed to update progress for book %s: %v", session.BookID, err)
	}

	return session, nil
}

func (s *SessionService) CancelSession(sessionID, userID string) error {
	return s.sessionRepo.CancelSession(sessionID, userID)
}

func (s *SessionService) GetActiveSession(userID string) (*models.ReadingSession, error) {
	return s.sessionRepo.GetActiveSession(userID)
}

func (s *SessionService) GetSessionsByBook(bookID, userID string) ([]models.ReadingSession, error) {
	return s.sessionRepo.GetSessionsByBookID(bookID, userID)
}

func (s *SessionService) GetProgress(bookID, userID string) (*models.ReadingProgress, error) {
	return s.sessionRepo.GetProgress(bookID, userID)
}
