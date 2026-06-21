package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	YearlyGoal   int       `json:"yearly_goal"`
	Timezone     string    `json:"timezone"`
	CreatedAt    time.Time `json:"created_at"`
}

type Book struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	Title      string     `json:"title"`
	Author     string     `json:"author"`
	CoverURL   string     `json:"cover_url"`
	TotalPages int        `json:"total_pages"`
	Status     string     `json:"status"`
	Rating     int        `json:"rating"`
	FinishedAt *time.Time `json:"finished_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

type ReadingSession struct {
	ID              string     `json:"id"`
	BookID          string     `json:"book_id"`
	UserID          string     `json:"user_id"`
	StartedAt       time.Time  `json:"started_at"`
	EndedAt         *time.Time `json:"ended_at"`
	DurationSeconds int        `json:"duration_seconds"`
	PagesRead       int        `json:"pages_read"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
}

type ReadingProgress struct {
	ID          string    `json:"id"`
	BookID      string    `json:"book_id"`
	UserID      string    `json:"user_id"`
	CurrentPage int       `json:"current_page"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RefreshToken struct {
	ID 			string 		`json:"id"`
	UserID 		string 		`json:"user_id"`
	Token 		string 		`json:"token"`
	FamilyID 	string 		`json:"family_id"`
	ParentToken *string		`json:"parent_token"`
	ExpiresAt 	time.Time 	`json:"expires_at"`
	CreatedAt 	time.Time 	`json:"created_at"`
}

type UsedRefreshToken struct {
    Token    string    `json:"token"`
    FamilyID string    `json:"family_id"`
    UserID   string    `json:"user_id"`
    UsedAt   time.Time `json:"used_at"`
}

type DailyReading struct {
	Date 	string 	`json:"date"`
	Pages 	int 	`json:"pages"`
	Seconds int 	`json:"seconds"`
}

type SessionNote struct {
	ID 			string 		`json:"id"`
	SessionID 	string 		`json:"session_id"`
	UserID 		string 		`json:"user_id"`
	Chapter 	string 		`json:"chapter"`
	Pages 		string 		`json:"pages"`
	Note 		string 		`json:"note"`
	CreatedAt 	time.Time	`json:"created_at"`
}

type CreateNoteRequest struct {
	Chapter string `json:"chapter"`
	Pages 	string `json:"pages"`
	Note 	string `json:"note" binding:"required"`
}
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,max=72"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	Timezone string `json:"timezone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateTimezoneRequest struct {
    Timezone string `json:"timezone" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateBookRequest struct {
	Title      string `json:"title" binding:"required"`
	Author     string `json:"author"`
	CoverURL   string `json:"cover_url" binding:"omitempty,url"`
	TotalPages int    `json:"total_pages" binding:"omitempty,min=1"`
	Status     string `json:"status" binding:"omitempty,oneof=want_to_read reading finished"`
}

type UpdateBookRequest struct {
	Title      *string `json:"title"`
	Author     *string `json:"author"`
	CoverURL   *string `json:"cover_url" binding:"omitempty,url"`
	TotalPages *int    `json:"total_pages" binding:"omitempty,min=1"`
	Status     *string `json:"status" binding:"omitempty,oneof=want_to_read reading finished"`
	Rating     *int    `json:"rating" binding:"omitempty,min=0,max=5"`
}

type StartSessionRequest struct {
	BookID string `json:"book_id" binding:"required"`
}

type UpdateGoalRequest struct {
	YearlyGoal int `json:"yearly_goal" binding:"required,min=1,max=1000"`
}

type StatsResult struct {
	TotalBooks       int `json:"total_books"`
	FinishedBooks    int `json:"finished_books"`
	TotalReadingTime int `json:"total_reading_time_seconds"`
	TotalPagesRead   int `json:"total_pages_read"`
	CurrentStreak    int `json:"current_streak"`
	YearlyGoal       int `json:"yearly_goal"`
	YearlyFinished   int `json:"yearly_finished"`
}
