package repositories

import (
	"booktracker/backend/internal/models"
	"database/sql"
	"time"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) StartSession(bookID, userID string) (*models.ReadingSession, error) {
	session := &models.ReadingSession{}
	err := r.db.QueryRow(
		`INSERT INTO reading_sessions (book_id, user_id, started_at, status)
		VALUES ($1, $2, $3, 'active')
		RETURNING id, book_id, user_id, started_at, ended_at, duration_seconds, pages_read, status, created_at`,
		bookID, userID, time.Now(),
	).Scan(
		&session.ID, &session.BookID, &session.UserID, &session.StartedAt,
		&session.EndedAt, &session.DurationSeconds, &session.PagesRead,
		&session.Status, &session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) PauseSession(sessionID, userID string, duration int) (*models.ReadingSession, error) {
	session := &models.ReadingSession{}
	err := r.db.QueryRow(
		`UPDATE reading_sessions SET status='paused', duration_seconds=$1
		WHERE id=$2 AND user_id=$3
		RETURNING id, book_id, user_id, started_at, ended_at, duration_seconds, pages_read, status, created_at`,
		duration, sessionID, userID,
	).Scan(
		&session.ID, &session.BookID, &session.UserID, &session.StartedAt,
		&session.EndedAt, &session.DurationSeconds, &session.PagesRead,
		&session.Status, &session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) ResumeSession(sessionID, userID string) (*models.ReadingSession, error) {
	session := &models.ReadingSession{}
	err := r.db.QueryRow(
		`UPDATE reading_sessions SET status='active'
		WHERE id=$1 AND user_id=$2
		RETURNING id, book_id, user_id, started_at, ended_at, duration_seconds, pages_read, status, created_at`,
		sessionID, userID,
	).Scan(
		&session.ID, &session.BookID, &session.UserID, &session.StartedAt,
		&session.EndedAt, &session.DurationSeconds, &session.PagesRead,
		&session.Status, &session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) StopSession(sessionID, userID string, duration, pagesRead int) (*models.ReadingSession, error) {
	now := time.Now()
	session := &models.ReadingSession{}
	err := r.db.QueryRow(
		`UPDATE reading_sessions 
		SET status='completed', ended_at=$1, duration_seconds=$2, pages_read=$3
		WHERE id=$4 AND user_id=$5
		RETURNING id, book_id, user_id, started_at, ended_at, duration_seconds, pages_read, status, created_at`,
		now, duration, pagesRead, sessionID, userID,
	).Scan(
		&session.ID, &session.BookID, &session.UserID, &session.StartedAt,
		&session.EndedAt, &session.DurationSeconds, &session.PagesRead,
		&session.Status, &session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) CancelSession(sessionID, userID string) error {
	_, err := r.db.Exec(
		`DELETE FROM reading_sessions 
        WHERE id=$1 AND user_id=$2 AND status IN ('active', 'paused')`,
		sessionID, userID,
	)
	return err
}

func (r *SessionRepository) GetActiveSession(userID string) (*models.ReadingSession, error) {
	session := &models.ReadingSession{}
	err := r.db.QueryRow(
		`SELECT id, book_id, user_id, started_at, ended_at, duration_seconds, pages_read, status, created_at
		FROM reading_sessions
		WHERE user_id=$1 AND status IN ('active', 'paused')
		ORDER BY created_at DESC LIMIT 1`,
		userID,
	).Scan(
		&session.ID, &session.BookID, &session.UserID, &session.StartedAt,
		&session.EndedAt, &session.DurationSeconds, &session.PagesRead,
		&session.Status, &session.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) GetSessionsByBookID(bookID, userID string) ([]models.ReadingSession, error) {
	sessions := make([]models.ReadingSession, 0)

	rows, err := r.db.Query(
		`SELECT id, book_id, user_id, started_at, ended_at, duration_seconds, pages_read, status, created_at
        FROM reading_sessions
        WHERE book_id=$1 AND user_id=$2 AND status='completed'
        ORDER BY created_at DESC`,
		bookID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.ReadingSession
		err := rows.Scan(
			&s.ID, &s.BookID, &s.UserID, &s.StartedAt,
			&s.EndedAt, &s.DurationSeconds, &s.PagesRead,
			&s.Status, &s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func (r *SessionRepository) UpdateProgress(bookID, userID string, currentPage int) error {
	_, err := r.db.Exec(
		`INSERT INTO reading_progress (book_id, user_id, current_page, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (book_id, user_id)
		DO UPDATE SET current_page=$3, updated_at=NOW()`,
		bookID, userID, currentPage,
	)
	return err
}

func (r *SessionRepository) GetProgress(bookID, userID string) (*models.ReadingProgress, error) {
	p := &models.ReadingProgress{}
	err := r.db.QueryRow(
		`SELECT id, book_id, user_id, current_page, updated_at
        FROM reading_progress WHERE book_id=$1 AND user_id=$2`,
		bookID, userID,
	).Scan(&p.ID, &p.BookID, &p.UserID, &p.CurrentPage, &p.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

func (r *SessionRepository) HasActiveSession(bookID, userID string) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM reading_sessions 
        WHERE book_id=$1 AND user_id=$2 AND status IN ('active', 'paused')`,
		bookID, userID,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
