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
	var timezone string

	err := r.db.QueryRowContext(
    	ctx,
		`SELECT
			COUNT(books.id) AS total_books,
			COUNT(books.id) FILTER (WHERE books.status = 'finished') AS finished_books,
			users.yearly_goal,
			COUNT(books.id) FILTER (
				WHERE books.status = 'finished'
				AND EXTRACT(YEAR FROM books.finished_at AT TIME ZONE COALESCE(users.timezone, 'UTC')) = 
					EXTRACT(YEAR FROM NOW() AT TIME ZONE COALESCE(users.timezone, 'UTC'))
			) AS yearly_finished,
			COALESCE(rs.total_time, 0) AS total_reading_time,
			COALESCE(rs.total_pages, 0) AS total_pages_read,
			COALESCE(users.timezone, 'UTC') AS timezone
		FROM users
		LEFT JOIN books ON books.user_id = users.id
		LEFT JOIN (
			SELECT user_id,
				SUM(duration_seconds) AS total_time,
				SUM(pages_read) AS total_pages
			FROM reading_sessions
			WHERE user_id = $1 AND status = 'completed'
			GROUP BY user_id
		) rs ON rs.user_id = users.id
		WHERE users.id = $1
		GROUP BY users.id, users.yearly_goal, rs.total_time, rs.total_pages, users.timezone`,
		userID,
	).Scan(
		&stats.TotalBooks,
		&stats.FinishedBooks,
		&stats.YearlyGoal,
		&stats.YearlyFinished,
		&stats.TotalReadingTime,
		&stats.TotalPagesRead,
		&timezone,
	)

	if err != nil && err != sql.ErrNoRows {
        return nil, err
    }

    streak, err := r.CalculateStreak(ctx, userID, timezone)
    if err != nil {
        return nil, err
    }

    stats.CurrentStreak = streak
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

func (r *StatsRepository) CalculateStreak(ctx context.Context, userID string, userTimeZone string) (int, error) {
    loc, err := time.LoadLocation(userTimeZone)
    if err != nil {
        loc = time.UTC
    }

    rows, err := r.db.QueryContext(
        ctx,
        `SELECT DISTINCT DATE(started_at AT TIME ZONE $2) as reading_date
        FROM reading_sessions
        WHERE user_id=$1 AND status='completed'
        ORDER BY reading_date DESC`,
        userID,
        userTimeZone,
    )
    if err != nil {
        return 0, err
    }
    defer rows.Close()

    var dates []time.Time
    for rows.Next() {
        var rd time.Time
        if err := rows.Scan(&rd); err != nil {
            return 0, err
        }
        normalizedDate := time.Date(rd.Year(), rd.Month(), rd.Day(), 0, 0, 0, 0, loc)
        dates = append(dates, normalizedDate)
    }

    if len(dates) == 0 {
        return 0, nil
    }

    now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	yesterday := today.AddDate(0, 0, -1)

    if !dates[0].Equal(today) && !dates[0].Equal(yesterday) {
        return 0, nil
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

    return streak, nil
}
