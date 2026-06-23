package repositories_test

import (
	"booktracker/backend/internal/models"
	"booktracker/backend/internal/repositories"
	"booktracker/backend/internal/testutil"
	"context"
	"testing"

	"github.com/google/uuid"
)

func setupStatsTestData(t *testing.T, db interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (interface{}, error)
}) {
	t.Helper()
}

func TestGetStats_NoBooks(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	authRepo := repositories.NewAuthRepository(db)
	statsRepo := repositories.NewStatsRepository(db)
	ctx := context.Background()

	user, err := authRepo.CreateUser(ctx, "Test User", "stats@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	stats, err := statsRepo.GetStats(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if stats.TotalBooks != 0 {
		t.Errorf("expected 0 total books, got: %d", stats.TotalBooks)
	}
	if stats.FinishedBooks != 0 {
		t.Errorf("expected 0 finished books, got: %d", stats.FinishedBooks)
	}
	if stats.TotalReadingTime != 0 {
		t.Errorf("expected 0 reading time, got: %d", stats.TotalReadingTime)
	}
	if stats.CurrentStreak != 0 {
		t.Errorf("expected 0 streak, got: %d", stats.CurrentStreak)
	}
}

func TestGetStats_WithBooks(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	authRepo := repositories.NewAuthRepository(db)
	bookRepo := repositories.NewBookRepository(db)
	statsRepo := repositories.NewStatsRepository(db)
	ctx := context.Background()

	user, err := authRepo.CreateUser(ctx, "Test User", "stats@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	book1, err := bookRepo.Create(ctx, user.ID, models.CreateBookRequest{
		Title:  "Finished Book",
		Status: "finished",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	_, err = db.ExecContext(ctx,
		`UPDATE books SET finished_at = NOW() WHERE id = $1`,
		book1.ID,
	)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	_, err = bookRepo.Create(ctx, user.ID, models.CreateBookRequest{
		Title:  "Reading Book",
		Status: "reading",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	stats, err := statsRepo.GetStats(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if stats.TotalBooks != 2 {
		t.Errorf("expected 2 total books, got: %d", stats.TotalBooks)
	}
	if stats.FinishedBooks != 1 {
		t.Errorf("expected 1 finished book, got: %d", stats.FinishedBooks)
	}
}

func TestGetStats_WithSessions(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	authRepo := repositories.NewAuthRepository(db)
	bookRepo := repositories.NewBookRepository(db)
	statsRepo := repositories.NewStatsRepository(db)
	ctx := context.Background()

	user, err := authRepo.CreateUser(ctx, "Test User", "stats@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	book, err := bookRepo.Create(ctx, user.ID, models.CreateBookRequest{
		Title:  "Test Book",
		Status: "reading",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	_, err = db.ExecContext(ctx,
		`INSERT INTO reading_sessions (book_id, user_id, started_at, ended_at, duration_seconds, pages_read, status)
		VALUES ($1, $2, NOW(), NOW() + INTERVAL '30 minutes', 1800, 20, 'completed')`,
		book.ID, user.ID,
	)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	stats, err := statsRepo.GetStats(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if stats.TotalReadingTime != 1800 {
		t.Errorf("expected 1800 seconds reading time, got: %d", stats.TotalReadingTime)
	}
	if stats.TotalPagesRead != 20 {
		t.Errorf("expected 20 pages read, got: %d", stats.TotalPagesRead)
	}
}

func TestCalculateStreak_NoSessions(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	authRepo := repositories.NewAuthRepository(db)
	statsRepo := repositories.NewStatsRepository(db)
	ctx := context.Background()

	user, err := authRepo.CreateUser(ctx, "Test User", "stats@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	streak, err := statsRepo.CalculateStreak(ctx, user.ID, "UTC")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if streak != 0 {
		t.Errorf("expected 0 streak, got: %d", streak)
	}
}

func TestCalculateStreak_Today(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	authRepo := repositories.NewAuthRepository(db)
	bookRepo := repositories.NewBookRepository(db)
	statsRepo := repositories.NewStatsRepository(db)
	ctx := context.Background()

	user, err := authRepo.CreateUser(ctx, "Test User", "stats@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	book, err := bookRepo.Create(ctx, user.ID, models.CreateBookRequest{
		Title:  "Test Book",
		Status: "reading",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	_, err = db.ExecContext(ctx,
		`INSERT INTO reading_sessions (book_id, user_id, started_at, ended_at, duration_seconds, pages_read, status)
		VALUES ($1, $2, NOW(), NOW() + INTERVAL '30 minutes', 1800, 20, 'completed')`,
		book.ID, user.ID,
	)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	streak, err := statsRepo.CalculateStreak(ctx, user.ID, "UTC")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if streak != 1 {
		t.Errorf("expected streak of 1, got: %d", streak)
	}
}

func TestCalculateStreak_ConsecutiveDays(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	authRepo := repositories.NewAuthRepository(db)
	bookRepo := repositories.NewBookRepository(db)
	statsRepo := repositories.NewStatsRepository(db)
	ctx := context.Background()

	user, err := authRepo.CreateUser(ctx, "Test User", "stats@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	book, err := bookRepo.Create(ctx, user.ID, models.CreateBookRequest{
		Title:  "Test Book",
		Status: "reading",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	for i := 0; i < 3; i++ {
		_, err = db.ExecContext(ctx,
			`INSERT INTO reading_sessions (book_id, user_id, started_at, ended_at, duration_seconds, pages_read, status)
			VALUES ($1, $2, NOW() - $3 * INTERVAL '1 day', NOW() - $3 * INTERVAL '1 day' + INTERVAL '30 minutes', 1800, 20, 'completed')`,
			book.ID, user.ID, i,
		)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
	}

	streak, err := statsRepo.CalculateStreak(ctx, user.ID, "UTC")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if streak != 3 {
		t.Errorf("expected streak of 3, got: %d", streak)
	}
}

func TestCalculateStreak_BrokenStreak(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	authRepo := repositories.NewAuthRepository(db)
	bookRepo := repositories.NewBookRepository(db)
	statsRepo := repositories.NewStatsRepository(db)
	ctx := context.Background()

	user, err := authRepo.CreateUser(ctx, "Test User", "stats@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	book, err := bookRepo.Create(ctx, user.ID, models.CreateBookRequest{
		Title:  "Test Book",
		Status: "reading",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	for _, daysAgo := range []int{0, 2} {
		_, err = db.ExecContext(ctx,
			`INSERT INTO reading_sessions (book_id, user_id, started_at, ended_at, duration_seconds, pages_read, status)
			VALUES ($1, $2, NOW() - $3 * INTERVAL '1 day', NOW() - $3 * INTERVAL '1 day' + INTERVAL '30 minutes', 1800, 20, 'completed')`,
			book.ID, user.ID, daysAgo,
		)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}
	}

	streak, err := statsRepo.CalculateStreak(ctx, user.ID, "UTC")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if streak != 1 {
		t.Errorf("expected streak of 1, got: %d", streak)
	}
}

func TestGetReadingHistory(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	authRepo := repositories.NewAuthRepository(db)
	bookRepo := repositories.NewBookRepository(db)
	statsRepo := repositories.NewStatsRepository(db)
	ctx := context.Background()

	user, err := authRepo.CreateUser(ctx, "Test User", "stats@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	book, err := bookRepo.Create(ctx, user.ID, models.CreateBookRequest{
		Title:  "Test Book",
		Status: "reading",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	_, err = db.ExecContext(ctx,
		`INSERT INTO reading_sessions (book_id, user_id, started_at, ended_at, duration_seconds, pages_read, status)
		VALUES ($1, $2, NOW(), NOW() + INTERVAL '30 minutes', 1800, 20, 'completed')`,
		book.ID, user.ID,
	)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	history, err := statsRepo.GetReadingHistory(ctx, user.ID, 7)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(history) != 7 {
		t.Errorf("expected 7 days of history, got: %d", len(history))
	}

	found := false
	for _, d := range history {
		if d.Pages == 20 && d.Seconds == 1800 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find today's session in history")
	}

	_ = uuid.New().String()
}

func TestGetReadingHistory_EmptyDays(t *testing.T) {
	db := testutil.NewTestDB(t)
	defer db.Close()
	testutil.CleanDB(t, db)

	authRepo := repositories.NewAuthRepository(db)
	statsRepo := repositories.NewStatsRepository(db)
	ctx := context.Background()

	user, err := authRepo.CreateUser(ctx, "Test User", "stats@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	history, err := statsRepo.GetReadingHistory(ctx, user.ID, 7)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(history) != 7 {
		t.Errorf("expected 7 days even with no sessions, got: %d", len(history))
	}
	for _, d := range history {
		if d.Pages != 0 || d.Seconds != 0 {
			t.Errorf("expected all zeros for empty history, got pages=%d seconds=%d", d.Pages, d.Seconds)
		}
	}
}