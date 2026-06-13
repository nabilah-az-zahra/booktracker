package main

import (
	"booktracker/backend/internal/handlers"
	"booktracker/backend/internal/middleware"
	"booktracker/backend/internal/models"
	"booktracker/backend/internal/repositories"
	"booktracker/backend/internal/services"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("POSTGRES_HOST") == "" {
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("Warning: No .env file found, relying on environment variables")
		}
	}

	db := models.NewDB()
	defer db.Close()

	authRepo := repositories.NewAuthRepository(db)
	bookRepo := repositories.NewBookRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	statsRepo := repositories.NewStatsRepository(db)

	authService := services.NewAuthService(authRepo)
	bookService := services.NewBookService(bookRepo, sessionRepo)
	sessionService := services.NewSessionService(sessionRepo, bookRepo)
	statsService := services.NewStatsService(statsRepo)

	authHandler := handlers.NewAuthHandler(authService)
	bookHandler := handlers.NewBookHandler(bookService)
	sessionHandler := handlers.NewSessionHandler(sessionService)
	statsHandler := handlers.NewStatsHandler(statsService)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.SetTrustedProxies(nil)

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5173"
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		auth.Use(middleware.AuthRateLimiter())
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", authHandler.GetProfile)
			protected.PATCH("/profile/goal", authHandler.UpdateGoal)

			protected.GET("/books", bookHandler.GetAllBooks)
			protected.GET("/books/:id", bookHandler.GetBook)
			protected.POST("/books", bookHandler.CreateBook)
			protected.PUT("/books/:id", bookHandler.UpdateBook)
			protected.DELETE("/books/:id", bookHandler.DeleteBook)

			protected.POST("/sessions/start", sessionHandler.StartSession)
			protected.PATCH("/sessions/:id/pause", sessionHandler.PauseSession)
			protected.PATCH("/sessions/:id/resume", sessionHandler.ResumeSession)
			protected.PATCH("/sessions/:id/stop", sessionHandler.StopSession)
			protected.DELETE("/sessions/:id", sessionHandler.CancelSession)
			protected.GET("/sessions/active", sessionHandler.GetActiveSession)
			protected.GET("/sessions/book/:bookId", sessionHandler.GetSessionsByBook)
			protected.GET("/progress/:bookId", sessionHandler.GetProgress)

			protected.GET("/stats", statsHandler.GetStats)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
