package repositories

import (
	"booktracker/backend/internal/models"
	"database/sql"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(name, email, passwordHash string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
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

func (r *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
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

func (r *AuthRepository) GetUserByID(id string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
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

func (r *AuthRepository) UpdateYearlyGoal(userID string, goal int) error {
	_, err := r.db.Exec(
		`UPDATE users SET yearly_goal = $1 WHERE id = $2`,
		goal, userID,
	)
	return err
}
