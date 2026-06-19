package repositories

import (
	"booktracker/backend/internal/models"
	"context"
	"database/sql"
	"time"
)

type StatsRepository struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) GetStats(ctx context.Context, userID string) (*models.StatsResult, error) {
	stats := &models.StatsResult{}

	err := r.db.QueryRowContext(
		ctx,
		`SELECT
            COUNT(*) as total_books,
            COUNT(CASE WHEN books.status='finished' THEN 1 END) as finished_books,
            users.yearly_goal,
            COUNT(CASE WHEN books.status='finished'
                AND EXTRACT(YEAR FROM books.finished_at) = EXTRACT(YEAR FROM NOW())
                THEN 1 END) as yearly_finished,
			COALESCE((
				SELECT SUM(duration_seconds)
				FROM reading_sessions
				WHERE user_id = users.id AND status = 'completed'
			), 0) as total_reading_time,
			COALESCE((
				SELECT SUM (pages_read)
				FROM reading_sessions
				WHERE user_id = users.id AND status = 'completed'
			), 0) as total_pages_read
        FROM users
        LEFT JOIN books ON books.user_id = users.id
        WHERE users.id = $1
        GROUP BY users.id, users.yearly_goal`,
		userID,
	).Scan(
		&stats.TotalBooks,
		&stats.FinishedBooks,
		&stats.YearlyGoal,
		&stats.YearlyFinished,
		&stats.TotalReadingTime,
		&stats.TotalPagesRead,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	stats.CurrentStreak = r.calculateStreak(ctx, userID)
	return stats, nil
}

func (r *StatsRepository) GetReadingHistory(ctx context.Context, userID string, days int) ([]models.DailyReading, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`WITH date_series AS (
			SELECT generate_series(
				(NOW() - ($2 - 1) * INTERVAL '1 day')::date,
				NOW()::date,
				INTERVAL '1 day'
			)::date AS reading_date
		),
		daily_sessions AS (
			SELECT
				DATE(started_at) AS reading_date,
				COALESCE(SUM(pages_read), 0) AS pages,
				COALESCE(SUM(duration_seconds), 0) AS seconds
			FROM reading_sessions
			WHERE user_id = $1
				AND status = 'completed'
				AND started_at >= NOW() - ($2 * INTERVAL '1 day')
			GROUP BY DATE(started_at)
		)
		SELECT
			ds.reading_date,
			COALESCE(s.pages, 0) AS pages,
			COALESCE(s.seconds, 0) AS seconds
		FROM date_series ds
		LEFT JOIN daily_sessions s ON s.reading_date = ds.reading_date
		ORDER BY ds.reading_date ASC`,
		userID, days,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := make([]models.DailyReading, 0 )
	for rows.Next() {
		var d models.DailyReading
		var date time.Time
		if err := rows.Scan(&date, &d.Pages, &d.Seconds); err != nil {
			return nil, err
		}
		d.Date = date.Format("2006-01-02")
		history = append(history, d)
	}
	return history, nil
}

func (r *StatsRepository) calculateStreak(ctx context.Context, userID string) int {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT DISTINCT DATE(started_at) as reading_date
		FROM reading_sessions
		WHERE user_id=$1 AND status='completed'
		ORDER BY reading_date DESC`,
		userID,
	)
	if err != nil {
		return 0
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var rd time.Time
		if err := rows.Scan(&rd); err == nil {
			dates = append(dates, rd.Truncate(24*time.Hour))
		}
	}

	if len(dates) == 0 {
		return 0
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)

	if !dates[0].Equal(today) && !dates[0].Equal(yesterday) {
		return 0
	}

	streak := 0
	expected := dates[0]

	for _, readingDate := range dates {
		if readingDate.Equal(expected) {
			streak++
			expected = readingDate.AddDate(0, 0, -1)
		} else {
			break
		}
	}
	return streak
}
