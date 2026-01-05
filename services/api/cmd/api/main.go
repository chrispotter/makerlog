package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chrispotter/makerlog/services/api/internal/database"
	"github.com/chrispotter/makerlog/services/api/internal/handlers"
	"github.com/chrispotter/makerlog/services/api/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

func main() {
	// Get configuration from environment
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/makerlog?sslmode=disable")
	sessionSecret := getEnv("SESSION_SECRET", "your-secret-key-change-this-in-production")
	port := getEnv("PORT", "8080")
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")

	// Warn if using default session secret
	if sessionSecret == "your-secret-key-change-this-in-production" {
		log.Println("WARNING: Using default session secret! This is insecure for production.")
		log.Println("Please set SESSION_SECRET environment variable to a strong random value.")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Test database connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to database successfully")

	// Initialize session store
	sessionStore := sessions.NewCookieStore([]byte(sessionSecret))
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}

	// Setup database queries
	queries := database.New(db)

	// Setup handlers
	authHandler := handlers.NewAuthHandler(queries, sessionStore)
	projectHandler := handlers.NewProjectHandler(queries)
	taskHandler := handlers.NewTaskHandler(queries)
	logEntryHandler := handlers.NewLogEntryHandler(queries)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{frontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(sessionStore))

		// Auth routes
		r.Post("/api/auth/logout", authHandler.Logout)
		r.Get("/api/auth/me", authHandler.Me)

		// Projects routes
		r.Get("/api/projects", projectHandler.List)
		r.Post("/api/projects", projectHandler.Create)
		r.Get("/api/projects/{id}", projectHandler.Get)
		r.Put("/api/projects/{id}", projectHandler.Update)
		r.Delete("/api/projects/{id}", projectHandler.Delete)

		// Tasks routes
		r.Get("/api/tasks", taskHandler.List)
		r.Post("/api/tasks", taskHandler.Create)
		r.Get("/api/tasks/{id}", taskHandler.Get)
		r.Put("/api/tasks/{id}", taskHandler.Update)
		r.Delete("/api/tasks/{id}", taskHandler.Delete)

		// Log entries routes
		r.Get("/api/log-entries", logEntryHandler.List)
		r.Post("/api/log-entries", logEntryHandler.Create)
		r.Get("/api/log-entries/{id}", logEntryHandler.Get)
		r.Put("/api/log-entries/{id}", logEntryHandler.Update)
		r.Delete("/api/log-entries/{id}", logEntryHandler.Delete)

		// Today route - get today's log entries
		r.Get("/api/today", logEntryHandler.Today)
	})

	// Start server with timeouts
	log.Printf("Server starting on port %s", port)
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
