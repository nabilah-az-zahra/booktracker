package repositories

import (
	"booktracker/backend/internal/models"
	"context"
	"database/sql"
	"time"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(ctx context.Context, name, email, passwordHash string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, yearly_goal, created_at`,
		name, email, passwordHash,
	).Scan(&user.ID, &user.Name, &user.Email, &user.YearlyGoal, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, email, password_hash, yearly_goal, created_at
		FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.YearlyGoal, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *AuthRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, email, yearly_goal, created_at
		FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.YearlyGoal, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *AuthRepository) UpdateYearlyGoal(ctx context.Context, userID string, goal int) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE users SET yearly_goal = $1 WHERE id = $2`,
		goal, userID,
	)
	return err
}

func (r *AuthRepository) CreateRefreshToken(ctx context.Context, userID, token string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO refresh_tokens (user_id, token, expires_at) 
		VALUES ($1, $2, $3)`,
		userID, token, expiresAt,
	)
	return err
}

func (r *AuthRepository) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	rt := &models.RefreshToken{}
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens WHERE token = $1 AND expires_at > NOW()`,
		token,
	).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rt, nil
}

func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM refresh_tokens WHERE token = $1`,
		token,
	)
	return err

}

func (r *AuthRepository) DeleteAllRefreshToken(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM refresh_tokens WHERE user_id = $1`,
		userID,
	)
	return err
}