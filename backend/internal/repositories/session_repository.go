package repositories

import (
	"booktracker/backend/internal/models"
	"context"
	"database/sql"
	"time"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) StartSession(ctx context.Context, bookID, userID string) (*models.ReadingSession, error) {
	session := &models.ReadingSession{}
	err := r.db.QueryRowContext(
		ctx,
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

func (r *SessionRepository) PauseSession(ctx context.Context, sessionID, userID string, duration int) (*models.ReadingSession, error) {
	session := &models.ReadingSession{}
	err := r.db.QueryRowContext(
		ctx,
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

func (r *SessionRepository) ResumeSession(ctx context.Context, sessionID, userID string) (*models.ReadingSession, error) {
	session := &models.ReadingSession{}
	err := r.db.QueryRowContext(
		ctx,
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

func (r *SessionRepository) StopSession(ctx context.Context, sessionID, userID string, duration, pagesRead int) (*models.ReadingSession, error) {
	now := time.Now()
	session := &models.ReadingSession{}
	err := r.db.QueryRowContext(
		ctx,
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

func (r *SessionRepository) CancelSession(ctx context.Context, sessionID, userID string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM reading_sessions 
        WHERE id=$1 AND user_id=$2 AND status IN ('active', 'paused')`,
		sessionID, userID,
	)
	return err
}

func (r *SessionRepository) GetActiveSession(ctx context.Context, userID string) (*models.ReadingSession, error) {
	session := &models.ReadingSession{}
	err := r.db.QueryRowContext(
		ctx,
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

func (r *SessionRepository) GetSessionsByBookID(ctx context.Context, bookID, userID string) ([]models.ReadingSession, error) {
	sessions := make([]models.ReadingSession, 0)

	rows, err := r.db.QueryContext(
		ctx,
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

func (r *SessionRepository) UpdateProgress(ctx context.Context, bookID, userID string, currentPage int) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO reading_progress (book_id, user_id, current_page, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (book_id, user_id)
		DO UPDATE SET current_page=$3, updated_at=NOW()`,
		bookID, userID, currentPage,
	)
	return err
}

func (r *SessionRepository) GetProgress(ctx context.Context, bookID, userID string) (*models.ReadingProgress, error) {
	p := &models.ReadingProgress{}
	err := r.db.QueryRowContext(
		ctx,
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

func (r *SessionRepository) HasActiveSession(ctx context.Context, bookID, userID string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM reading_sessions 
        WHERE book_id=$1 AND user_id=$2 AND status IN ('active', 'paused')`,
		bookID, userID,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *SessionRepository) CreateNote(ctx context.Context, sessionID, userID string, req models.CreateNoteRequest) (*models.SessionNote, error) {
	note := &models.SessionNote{}
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO session_notes (session_id, user_id, chapter, pages, note)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, session_id, user_id, chapter, pages, note, created_at`,
		sessionID, userID, req.Chapter, req.Pages, req.Note,
	).Scan(
		&note.ID, &note.SessionID, &note.UserID,
		&note.Chapter, &note.Pages, &note.Note, &note.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (r *SessionRepository) GetNotesBySessionID(ctx context.Context, sessionID, userID string) ([]models.SessionNote, error) {
	notes := make([]models.SessionNote, 0)
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, session_id, user_id, chapter, pages, note, created_at
		FROM session_notes
		WHERE session_id = $1 AND user_id = $2
		ORDER BY created_at ASC`,
		sessionID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var n models.SessionNote
		err := rows.Scan(
			&n.ID, &n.SessionID, &n.UserID,
			&n.Chapter, &n.Pages, &n.Note, &n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func (r *SessionRepository) DeleteNote(ctx context.Context, noteID, userID string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM session_notes WHERE id = $1 AND user_id = $2`,
		noteID, userID,
	)
	return err
}