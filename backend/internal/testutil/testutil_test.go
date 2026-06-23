package testutil

import (
	"testing"
)

func TestDBConnection(t *testing.T) {
	db := NewTestDB(t)
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("expected db ping to succeed, got: %v", err)
	}

	t.Log("test database connection successful")
}