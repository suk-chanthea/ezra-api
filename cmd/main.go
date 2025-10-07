package main

import (
	"log"
	"os"

	"github.com/suk-chanthea/ezra/infrastructure/persistence"
	"github.com/suk-chanthea/ezra/interface/http/handler"
	"github.com/suk-chanthea/ezra/interface/http/router"
	"github.com/suk-chanthea/ezra/usecase"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Port        string
	PostgresURL string
	SecretKey   string
}

func loadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	pg := os.Getenv("POSTGRES_URL")
	if pg == "" {
		pg = "postgres://postgres:secret@postgres:5432/ezradb?sslmode=disable"
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "paracletus"
	}

	return &Config{
		Port:        port,
		PostgresURL: pg,
		SecretKey:   secret,
	}
}

func main() {
	// Load configuration
	config := loadConfig()

	// Connect to database
	db, err := gorm.Open(postgres.Open(config.PostgresURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ failed to connect database: %v", err)
	}
	log.Println("✅ PostgreSQL connected")

	// Initialize repositories (Infrastructure layer)
	userRepo := persistence.NewUserRepository(db)
	eventRepo := persistence.NewEventRepository(db)

	// Initialize use cases (Application layer)
	authUseCase := usecase.NewAuthUseCase(userRepo, config.SecretKey)
	eventUseCase := usecase.NewEventUseCase(eventRepo)

	// Initialize handlers (Interface layer)
	authHandler := handler.NewAuthHandler(authUseCase)
	eventHandler := handler.NewEventHandler(eventUseCase)

	// Setup router
	r := router.NewRouter(authHandler, eventHandler, authUseCase)
	engine := r.Setup()

	// Start server
	log.Printf("🚀 Server starting on port %s", config.Port)
	if err := engine.Run(":" + config.Port); err != nil {
		log.Fatalf("❌ failed to start server: %v", err)
	}
}
