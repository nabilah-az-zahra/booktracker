package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("POSTGRES_HOST") == "" {
		if err := godotenv.Load("../.env.local"); err != nil {
			godotenv.Load("../.env")
		}
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		getEnvOrDefault("POSTGRES_SSLMODE", "disable"),
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatalf("failed to create migrations table: %v", err)
	}

	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		log.Fatalf("failed to find migration files: %v", err)
	}
	sort.Strings(files)

	if len(files) == 0 {
		log.Println("no migration files found")
		return
	}

	for _, file := range files {
		filename := filepath.Base(file)

		var count int
		err := db.QueryRow(
			`SELECT COUNT(*) FROM schema_migrations WHERE filename = $1`,
			filename,
		).Scan(&count)
		if err != nil {
			log.Fatalf("failed to check migration status: %v", err)
		}

		if count > 0 {
			log.Printf("skipping %s (already applied)", filename)
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("failed to read migration file %s: %v", filename, err)
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("failed to begin transaction: %v", err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			log.Fatalf("failed to apply migration %s: %v", filename, err)
		}

		if _, err := tx.Exec(
			`INSERT INTO schema_migrations (filename) VALUES ($1)`,
			filename,
		); err != nil {
			tx.Rollback()
			log.Fatalf("failed to record migration %s: %v", filename, err)
		}

		if err := tx.Commit(); err != nil {
			log.Fatalf("failed to commit migration %s: %v", filename, err)
		}

		log.Printf("applied %s", filename)
	}

	log.Println("migrations complete")
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}