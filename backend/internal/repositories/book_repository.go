package repositories

import (
	"booktracker/backend/internal/apperrors"
	"booktracker/backend/internal/models"
	"database/sql"
	"errors"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetAllByUserID(userID string) ([]models.Book, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, title, author, cover_url, total_pages, status, rating, finished_at, created_at
		FROM books WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := make([]models.Book, 0)
	for rows.Next() {
		var b models.Book
		err := rows.Scan(
			&b.ID, &b.UserID, &b.Title, &b.Author, &b.CoverURL,
			&b.TotalPages, &b.Status, &b.Rating, &b.FinishedAt, &b.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (r *BookRepository) GetByID(bookID, userID string) (*models.Book, error) {
	b := &models.Book{}
	err := r.db.QueryRow(
		`SELECT id, user_id, title, author, cover_url, total_pages, status, rating, finished_at, created_at
		FROM books WHERE id = $1 AND user_id = $2`,
		bookID, userID,
	).Scan(
		&b.ID, &b.UserID, &b.Title, &b.Author, &b.CoverURL,
		&b.TotalPages, &b.Status, &b.Rating, &b.FinishedAt, &b.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return b, nil
}

func (r *BookRepository) Create(userID string, req models.CreateBookRequest) (*models.Book, error) {
	status := req.Status
	if status == "" {
		status = "want_to_read"
	}
	b := &models.Book{}
	err := r.db.QueryRow(
		`INSERT INTO books (user_id, title, author, cover_url, total_pages, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, title, author, cover_url, total_pages, status, rating, finished_at, created_at`,
		userID, req.Title, req.Author, req.CoverURL, req.TotalPages, status,
	).Scan(
		&b.ID, &b.UserID, &b.Title, &b.Author, &b.CoverURL,
		&b.TotalPages, &b.Status, &b.Rating, &b.FinishedAt, &b.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *BookRepository) Update(bookID, userID string, req models.UpdateBookRequest) (*models.Book, error) {
	b := &models.Book{}
	err := r.db.QueryRow(
		`UPDATE books SET 
			title = COALESCE($1, title), 
			author = COALESCE($2, author), 
			cover_url = COALESCE($3, cover_url), 
			total_pages = COALESCE($4, total_pages), 
			status = COALESCE($5, status), 
			rating = COALESCE($6, rating), 
			finished_at = CASE
				WHEN COALESCE($5, status) = 'finished' AND finished_at IS NULL THEN NOW()
				WHEN COALESCE($5, status) != 'finished' THEN NULL
				ELSE finished_at
			END
		WHERE id=$7 AND user_id=$8
		RETURNING id, user_id, title, author, cover_url, total_pages, status, rating, finished_at, created_at`,
		req.Title, req.Author, req.CoverURL, req.TotalPages,
		req.Status, req.Rating, bookID, userID,
	).Scan(
		&b.ID, &b.UserID, &b.Title, &b.Author, &b.CoverURL,
		&b.TotalPages, &b.Status, &b.Rating, &b.FinishedAt, &b.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return b, nil
}

func (r *BookRepository) Delete(bookID, userID string) error {
	result, err := r.db.Exec(
		`DELETE FROM books WHERE id = $1 AND user_id = $2`,
		bookID, userID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
