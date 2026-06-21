package services

import (
	"booktracker/backend/internal/apperrors"
	"booktracker/backend/internal/models"
	redisclient "booktracker/backend/internal/redis"
	"booktracker/backend/internal/repositories"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo *repositories.AuthRepository
}

func NewAuthService(authRepo *repositories.AuthRepository) *AuthService {
	return &AuthService{authRepo: authRepo}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.authRepo.CreateUser(ctx, req.Name, req.Email, string(hashedPassword), req.Timezone)
	if err != nil {
		log.Printf("registration error: %v", err)
		return nil, errors.New("registration failed, try again.")
	}

	token, err := generateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, string, error) {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.authRepo.GetUserByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	accessToken, err := generateAccessToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, "", err
	}

	familyID := uuid.New().String()
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	if err := s.authRepo.CreateRefreshToken(ctx, user.ID, refreshToken, familyID, expiresAt); err != nil {
		return nil, "", err
	}

	return &models.AuthResponse{Token: accessToken, User: *user}, refreshToken, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, string, error) {
	rt, err := s.authRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, "", errors.New("invalid or expired refresh token")
	}
	if rt == nil {
		usedToken, err := s.authRepo.GetUsedToken(ctx, refreshToken)
		if err != nil {
			return nil, "", errors.New("invalid or expired refresh token")
		}
		if usedToken != nil {
			log.Printf("refresh token reuse detected for family %s and invalidating entire family", usedToken.FamilyID)
			s.authRepo.DeleteRefreshTokenFamily(ctx, usedToken.FamilyID)
			return nil, "", errors.New("invalid or expired refresh token")
		}
		return  nil, "", errors.New("invalid or expired refresh token")
	}

	user, err := s.authRepo.GetUserByID(ctx, rt.UserID)
	if err != nil || user == nil {
		return nil, "", errors.New("invalid or expired refresh token")
	}

	if err := s.authRepo.MarkTokenAsUsed(ctx, refreshToken, rt.FamilyID, rt.UserID); err != nil {
		log.Printf("warning: failed to mark token as used: %v", err)
	}

	if err := s.authRepo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return nil, "", err
	}

	accessToken, err := generateAccessToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, "", err
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	if err := s.authRepo.CreateRefreshTokenInFamily(ctx, user.ID, newRefreshToken, rt.FamilyID, refreshToken, expiresAt); err != nil {
		return nil, "", err
	}

	return &models.AuthResponse{Token: accessToken, User: *user}, newRefreshToken, nil
}

func (s * AuthService) Logout(ctx context.Context, accessToken, refreshToken string) error {
	if accessToken != "" {
		remaining := time.Until(time.Now().Add(15 * time.Minute))
		if remaining > 0 {
			if err := redisclient.Client.Set(
				ctx,
				"blacklist:"+accessToken,
				1,
				remaining,
			).Err(); err != nil {
				log.Printf("warning: failed to blacklist token in Redis: %v", err)
			}
		}
	}

	if refreshToken != "" {
		return s.authRepo.DeleteRefreshToken(ctx, refreshToken)
	}
	
	return nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	return s.authRepo.GetUserByID(ctx, userID)
}

func (s *AuthService) UpdateGoal(ctx context.Context, userID string, goal int) error {
	return s.authRepo.UpdateYearlyGoal(ctx, userID, goal)
}

func (s *AuthService) UpdateTimezone(ctx context.Context, userID, timezone string) error {
    if _, err := time.LoadLocation(timezone); err != nil {
        return apperrors.ErrInvalidTimezone 
    }
    return s.authRepo.UpdateTimezone(ctx, userID, timezone)
}

func generateAccessToken(userID string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not configured")
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}