package repositories_test

import (
	"booktracker/backend/internal/models"
	"booktracker/backend/internal/repositories"
	"booktracker/backend/internal/testutil"
	"context"
	"testing"
)

func setupSessionTestData(t *testing.T) (*repositories.AuthRepository, *repositories.BookRepository, *repositories.SessionRepository, *models.User, *models.Book, context.Context) {
	t.Helper()
	db := testutil.NewTestDB(t)
	testutil.CleanDB(t, db)

	authRepo := repositories.NewAuthRepository(db)
	bookRepo := repositories.NewBookRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	ctx := context.Background()

	user, err := authRepo.CreateUser(ctx, "Test User", "session@example.com", "hashedpassword", "UTC")
	if err != nil {
		t.Fatalf("setup failed creating user: %v", err)
	}

	book, err := bookRepo.Create(ctx, user.ID, models.CreateBookRequest{
		Title:      "Test Book",
		Status:     "reading",
		TotalPages: 300,
	})
	if err != nil {
		t.Fatalf("setup failed creating book: %v", err)
	}

	return authRepo, bookRepo, sessionRepo, user, book, ctx
}

func TestStartSession(t *testing.T) {
	_, _, sessionRepo, user, book, ctx := setupSessionTestData(t)

	session, err := sessionRepo.StartSession(ctx, book.ID, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if session.ID == "" {
		t.Error("expected session ID to be set")
	}
	if session.Status != "active" {
		t.Errorf("expected status active, got: %s", session.Status)
	}
	if session.BookID != book.ID {
		t.Errorf("expected book ID %s, got: %s", book.ID, session.BookID)
	}
}

func TestPauseSession(t *testing.T) {
	_, _, sessionRepo, user, book, ctx := setupSessionTestData(t)

	session, err := sessionRepo.StartSession(ctx, book.ID, user.ID)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	paused, err := sessionRepo.PauseSession(ctx, session.ID, user.ID, 300)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if paused.Status != "paused" {
		t.Errorf("expected status paused, got: %s", paused.Status)
	}
	if paused.DurationSeconds != 300 {
		t.Errorf("expected duration 300, got: %d", paused.DurationSeconds)
	}
}

func TestResumeSession(t *testing.T) {
	_, _, sessionRepo, user, book, ctx := setupSessionTestData(t)

	session, err := sessionRepo.StartSession(ctx, book.ID, user.ID)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	_, err = sessionRepo.PauseSession(ctx, session.ID, user.ID, 300)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	resumed, err := sessionRepo.ResumeSession(ctx, session.ID, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resumed.Status != "active" {
		t.Errorf("expected status active, got: %s", resumed.Status)
	}
}

func TestStopSessionWithProgress(t *testing.T) {
	_, _, sessionRepo, user, book, ctx := setupSessionTestData(t)

	session, err := sessionRepo.StartSession(ctx, book.ID, user.ID)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	stopped, err := sessionRepo.StopSessionWithProgress(ctx, session.ID, user.ID, 1800, 20)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if stopped.Status != "completed" {
		t.Errorf("expected status completed, got: %s", stopped.Status)
	}
	if stopped.DurationSeconds != 1800 {
		t.Errorf("expected duration 1800, got: %d", stopped.DurationSeconds)
	}
	if stopped.PagesRead != 20 {
		t.Errorf("expected 20 pages read, got: %d", stopped.PagesRead)
	}

	progress, err := sessionRepo.GetProgress(ctx, book.ID, user.ID)
	if err != nil {
		t.Fatalf("expected no error getting progress, got: %v", err)
	}
	if progress == nil {
		t.Fatal("expected progress to be set, got nil")
	}
	if progress.CurrentPage != 20 {
		t.Errorf("expected current page 20, got: %d", progress.CurrentPage)
	}
}

func TestStopSessionWithProgress_PagesCapped(t *testing.T) {
	_, _, sessionRepo, user, book, ctx := setupSessionTestData(t)

	session, err := sessionRepo.StartSession(ctx, book.ID, user.ID)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	stopped, err := sessionRepo.StopSessionWithProgress(ctx, session.ID, user.ID, 1800, 400)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if stopped.Status != "completed" {
		t.Errorf("expected status completed, got: %s", stopped.Status)
	}

	progress, err := sessionRepo.GetProgress(ctx, book.ID, user.ID)
	if err != nil {
		t.Fatalf("expected no error getting progress, got: %v", err)
	}
	if progress.CurrentPage != 300 {
		t.Errorf("expected current page capped at 300, got: %d", progress.CurrentPage)
	}
}

func TestCancelSession(t *testing.T) {
	_, _, sessionRepo, user, book, ctx := setupSessionTestData(t)

	session, err := sessionRepo.StartSession(ctx, book.ID, user.ID)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	err = sessionRepo.CancelSession(ctx, session.ID, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	active, err := sessionRepo.GetActiveSession(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if active != nil {
		t.Error("expected no active session after cancel, got one")
	}
}

func TestGetActiveSession(t *testing.T) {
	_, _, sessionRepo, user, book, ctx := setupSessionTestData(t)

	active, err := sessionRepo.GetActiveSession(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if active != nil {
		t.Error("expected nil active session, got one")
	}

	_, err = sessionRepo.StartSession(ctx, book.ID, user.ID)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	active, err = sessionRepo.GetActiveSession(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if active == nil {
		t.Fatal("expected active session, got nil")
	}
	if active.Status != "active" {
		t.Errorf("expected status active, got: %s", active.Status)
	}
}

func TestCreateAndGetNotes(t *testing.T) {
	_, _, sessionRepo, user, book, ctx := setupSessionTestData(t)

	session, err := sessionRepo.StartSession(ctx, book.ID, user.ID)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	note, err := sessionRepo.CreateNote(ctx, session.ID, user.ID, models.CreateNoteRequest{
		Chapter: "Chapter 1",
		Pages:   "1-20",
		Note:    "Interesting start",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if note.ID == "" {
		t.Error("expected note ID to be set")
	}
	if note.Note != "Interesting start" {
		t.Errorf("expected note text, got: %s", note.Note)
	}

	notes, err := sessionRepo.GetNotesBySessionID(ctx, session.ID, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(notes) != 1 {
		t.Errorf("expected 1 note, got: %d", len(notes))
	}
}

func TestDeleteNote(t *testing.T) {
	_, _, sessionRepo, user, book, ctx := setupSessionTestData(t)

	session, err := sessionRepo.StartSession(ctx, book.ID, user.ID)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	note, err := sessionRepo.CreateNote(ctx, session.ID, user.ID, models.CreateNoteRequest{
		Note: "Test note",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	err = sessionRepo.DeleteNote(ctx, note.ID, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	notes, err := sessionRepo.GetNotesBySessionID(ctx, session.ID, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(notes) != 0 {
		t.Errorf("expected 0 notes after delete, got: %d", len(notes))
	}
}

func TestUpdateProgress(t *testing.T) {
	_, _, sessionRepo, user, book, ctx := setupSessionTestData(t)

	err := sessionRepo.UpdateProgress(ctx, book.ID, user.ID, 50)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	progress, err := sessionRepo.GetProgress(ctx, book.ID, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if progress == nil {
		t.Fatal("expected progress, got nil")
	}
	if progress.CurrentPage != 50 {
		t.Errorf("expected page 50, got: %d", progress.CurrentPage)
	}

	err = sessionRepo.UpdateProgress(ctx, book.ID, user.ID, 100)
	if err != nil {
		t.Fatalf("expected no error on upsert, got: %v", err)
	}

	progress, err = sessionRepo.GetProgress(ctx, book.ID, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if progress.CurrentPage != 100 {
		t.Errorf("expected page 100 after upsert, got: %d", progress.CurrentPage)
	}
}