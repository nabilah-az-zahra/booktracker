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

func (r *AuthRepository) CreateUser(ctx context.Context, name, email, passwordHash, timezone string) (*models.User, error) {
	if timezone == "" {
        timezone = "UTC"
    }
	user := &models.User{}
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO users (name, email, password_hash, timezone)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, yearly_goal, timezone, created_at`,
		name, email, passwordHash, timezone,
	).Scan(&user.ID, &user.Name, &user.Email, &user.YearlyGoal, &user.Timezone, &user.CreatedAt)
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

func (r *AuthRepository) UpdateTimezone(ctx context.Context, userID, timezone string) error {
    _, err := r.db.ExecContext(
        ctx,
        `UPDATE users SET timezone = $1 WHERE id = $2`,
        timezone, userID,
    )
    return err
}

func (r *AuthRepository) CreateRefreshToken(ctx context.Context, userID, token, familyID string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO refresh_tokens (user_id, token, family_id, expires_at) 
		VALUES ($1, $2, $3, $4)`,
		userID, token, familyID, expiresAt,
	)
	return err
}

func (r *AuthRepository) CreateRefreshTokenInFamily(ctx context.Context, userID, token, familyID, parentToken string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO refresh_tokens (user_id, token, family_id, parent_token, expires_at)
		VALUES ($1, $2, $3, $4, $5)`,
		userID, token, familyID, parentToken, expiresAt,	
	)
	return err
}

func (r *AuthRepository) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	rt := &models.RefreshToken{}
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, token, family_id, parent_token, expires_at, created_at
		FROM refresh_tokens WHERE token = $1 AND expires_at > NOW()`,
		token,
	).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.FamilyID, &rt.ParentToken, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rt, nil
}

func (r *AuthRepository) MarkTokenAsUsed(ctx context.Context, token, familyID, userID string) error {
    _, err := r.db.ExecContext(
		ctx,
        `INSERT INTO used_refresh_tokens (token, family_id, user_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (token) DO NOTHING`,
        token, familyID, userID,
    )
    return err
}

func (r *AuthRepository) GetUsedToken(ctx context.Context, token string) (*models.UsedRefreshToken, error) {
    ut := &models.UsedRefreshToken{}
    err := r.db.QueryRowContext(
		ctx,
        `SELECT token, family_id, user_id, used_at
        FROM used_refresh_tokens WHERE token = $1`,
        token,
    ).Scan(&ut.Token, &ut.FamilyID, &ut.UserID, &ut.UsedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return ut, nil
}

func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM refresh_tokens WHERE token = $1`,
		token,
	)
	return err
}

func (r *AuthRepository) DeleteRefreshTokenFamily(ctx context.Context, familyID string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM refresh_tokens WHERE family_id = $1`,
		familyID,
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