package repositories_test

import (
	"booktracker/backend/internal/models"
	"booktracker/backend/internal/repositories"
	"booktracker/backend/internal/testutil"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateUser(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	repo := repositories.NewAuthRepository(db)
	ctx := context.Background()

	user, err := repo.CreateUser(ctx, "Test User", "test@example.com", "hashedpassword", "Asia/Singapore")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if user.ID == "" {
		t.Error("expected user ID to be set")
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got: %s", user.Email)
	}
	if user.Timezone != "Asia/Singapore" {
		t.Errorf("expected timezone Asia/Singapore, got: %s", user.Timezone)
	}
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	repo := repositories.NewAuthRepository(db)
	ctx := context.Background()

	_, err := repo.CreateUser(ctx, "Test User", "test@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("first insert failed: %v", err)
	}

	_, err = repo.CreateUser(ctx, "Test User 2", "test@example.com", "hashedpassword2", "UTC")
	if err == nil {
		t.Error("expected error for duplicate email, got nil")
	}
}

func TestGetUserByEmail(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	repo := repositories.NewAuthRepository(db)
	ctx := context.Background()

	_, err := repo.CreateUser(ctx, "Test User", "test@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	user, err := repo.GetUserByEmail(ctx, "test@example.com")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got: %s", user.Email)
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	repo := repositories.NewAuthRepository(db)
	ctx := context.Background()

	user, err := repo.GetUserByEmail(ctx, "notexist@example.com")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if user != nil {
		t.Error("expected nil user for non-existent email")
	}
}

func TestCreateAndGetRefreshToken(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	repo := repositories.NewAuthRepository(db)
	ctx := context.Background()

	user, err := repo.CreateUser(ctx, "Test User", "test@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	familyID := uuid.New().String()
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	err = repo.CreateRefreshToken(ctx, user.ID, "testtoken123", familyID, expiresAt)
	if err != nil {
		t.Fatalf("expected no error creating token, got: %v", err)
	}

	rt, err := repo.GetRefreshToken(ctx, "testtoken123")
	if err != nil {
		t.Fatalf("expected no error getting token, got: %v", err)
	}
	if rt == nil {
		t.Fatal("expected token, got nil")
	}
	if rt.UserID != user.ID {
		t.Errorf("expected userID %s, got: %s", user.ID, rt.UserID)
	}
}

func TestGetRefreshToken_Expired(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	repo := repositories.NewAuthRepository(db)
	ctx := context.Background()

	user, err := repo.CreateUser(ctx, "Test User", "test@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	familyID := uuid.New().String()
	expiresAt := time.Now().Add(-1 * time.Hour)
	err = repo.CreateRefreshToken(ctx, user.ID, "expiredtoken", familyID, expiresAt)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	rt, err := repo.GetRefreshToken(ctx, "expiredtoken")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if rt != nil {
		t.Error("expected nil for expired token, got token")
	}
}

func TestDeleteRefreshTokenFamily(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	repo := repositories.NewAuthRepository(db)
	ctx := context.Background()

	user, err := repo.CreateUser(ctx, "Test User", "test@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	familyID := uuid.New().String()
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	err = repo.CreateRefreshToken(ctx, user.ID, "token1", familyID, expiresAt)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	err = repo.CreateRefreshTokenInFamily(ctx, user.ID, "token2", familyID, "token1", expiresAt)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	err = repo.DeleteRefreshTokenFamily(ctx, familyID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	rt1, _ := repo.GetRefreshToken(ctx, "token1")
	rt2, _ := repo.GetRefreshToken(ctx, "token2")
	if rt1 != nil || rt2 != nil {
		t.Error("expected all tokens in family to be deleted")
	}
}

func TestMarkTokenAsUsed(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	repo := repositories.NewAuthRepository(db)
	ctx := context.Background()

	user, err := repo.CreateUser(ctx, "Test User", "test@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	familyID := uuid.New().String()
	err = repo.MarkTokenAsUsed(ctx, "usedtoken123", familyID, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	ut, err := repo.GetUsedToken(ctx, "usedtoken123")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if ut == nil {
		t.Fatal("expected used token, got nil")
	}
	if ut.UserID != user.ID {
		t.Errorf("expected userID %s, got: %s", user.ID, ut.UserID)
	}
}

func TestUpdateTimezone(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	repo := repositories.NewAuthRepository(db)
	ctx := context.Background()

	user, err := repo.CreateUser(ctx, "Test User", "test@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	err = repo.UpdateTimezone(ctx, user.ID, "Asia/Singapore")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	updated, err := repo.GetUserByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if updated.Timezone != "Asia/Singapore" {
		t.Errorf("expected timezone Asia/Singapore, got: %s", updated.Timezone)
	}
}

func setupTestUser(t *testing.T, repo *repositories.AuthRepository, ctx context.Context) *models.User {
	t.Helper()
	user, err := repo.CreateUser(ctx, "Test User", "test@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return user
}