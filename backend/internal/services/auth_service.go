package services

import (
	"booktracker/backend/internal/models"
	"booktracker/backend/internal/repositories"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo *repositories.AuthRepository
}

func NewAuthService(authRepo *repositories.AuthRepository) *AuthService {
	return &AuthService{authRepo: authRepo}
}

func (s *AuthService) Register(req models.RegisterRequest) (*models.AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.authRepo.CreateUser(req.Name, req.Email, string(hashedPassword))
	if err != nil {
		log.Printf("registration error: %v", err)
		return nil, errors.New("registration failed, try again.")
	}

	token, err := generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.authRepo.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) GetProfile(userID string) (*models.User, error) {
	return s.authRepo.GetUserByID(userID)
}

func (s *AuthService) UpdateGoal(userID string, goal int) error {
	return s.authRepo.UpdateYearlyGoal(userID, goal)
}

func generateToken(userID string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not configured")
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
