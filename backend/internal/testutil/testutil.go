package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	redisclient "booktracker/backend/internal/redis"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func NewTestDB(t *testing.T) *sql.DB {
	t.Helper()

	godotenv.Load("../../.env.test")

	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test_secret_key_for_testing_only")
	}

	if redisclient.Client == nil {
        redisclient.Init()
    }

	host := os.Getenv("TEST_POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("TEST_POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("TEST_POSTGRES_USER")
	if user == "" {
		user = "booktracker"
	}
	password := os.Getenv("TEST_POSTGRES_PASSWORD")
	if password == "" {
		password = ""
	}
	dbname := os.Getenv("TEST_POSTGRES_DB")
	if dbname == "" {
		dbname = "booktracker_test"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	return db
}

func CleanDB(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.ExecContext(context.Background(),
		`TRUNCATE TABLE session_notes, used_refresh_tokens, refresh_tokens, 
		reading_progress, reading_sessions, books, users RESTART IDENTITY CASCADE`,
	)
	if err != nil {
		t.Fatalf("failed to clean test database: %v", err)
	}
}