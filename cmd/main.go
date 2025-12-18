package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"

	"github.com/elghazx/go-vault/internal/adapters/handlers"
	"github.com/elghazx/go-vault/internal/adapters/repositories"
	"github.com/elghazx/go-vault/internal/adapters/storage"
	"github.com/elghazx/go-vault/internal/core/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Get configuration from environment
	dbPath := getEnv("DB_PATH", "./go-vault.db")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	port := getEnv("PORT", "8080")
	uploadsPath := getEnv("UPLOADS_PATH", "./uploads")

	// Initialize database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	userRepo := repositories.NewSQLiteUserRepository(db)
	fileRepo := repositories.NewSQLiteFileRepository(db)

	// Initialize storage
	fileStorage := storage.NewLocalFileStorage(uploadsPath)

	// Initialize services
	authService := services.NewAuthService(userRepo, jwtSecret)
	fileService := services.NewFileService(fileRepo, fileStorage)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	fileHandler := handlers.NewFileHandler(fileService, authService)

	// Setup Echo
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// Auth routes
	e.POST("/api/login", authHandler.Login)
	e.POST("/api/register", authHandler.Register)
	e.POST("/api/logout", authHandler.Logout)
	e.GET("/api/check-auth", authHandler.CheckAuth)

	// File routes
	e.POST("/api/upload", fileHandler.Upload)
	e.GET("/api/my-files", fileHandler.GetMyFiles)
	e.GET("/f/*", fileHandler.Preview)
	e.GET("/d/*", fileHandler.Download)

	// Static files
	e.Static("/", "./static")

	log.Printf("Go-Vault server starting on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}

func runMigrations(db *sql.DB) error {
	migration, err := ioutil.ReadFile("./migrations/001_initial.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(migration))
	return err
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
