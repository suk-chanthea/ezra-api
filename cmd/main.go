package main

import (
	"log"
	"os"
	"time"

	"github.com/suk-chanthea/ezra/infrastructure/firebase"
	"github.com/suk-chanthea/ezra/infrastructure/persistence"
	"github.com/suk-chanthea/ezra/interface/http/handler"
	"github.com/suk-chanthea/ezra/interface/http/router"
	"github.com/suk-chanthea/ezra/usecase"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Port                   string
	PostgresURL            string
	SecretKey              string
	GoogleClientID         string
	FirebaseCredentialPath string
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

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		googleClientID = "" // Set via environment variable
	}

	firebaseCredPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	// Optional: If not set, FCM will be disabled (dummy service)

	return &Config{
		Port:                   port,
		PostgresURL:            pg,
		SecretKey:              secret,
		GoogleClientID:         googleClientID,
		FirebaseCredentialPath: firebaseCredPath,
	}
}

func main() {
	// Set timezone for the application
	loc, err := time.LoadLocation("Asia/Phnom_Penh")
	if err != nil {
		log.Printf("⚠️  Warning: Could not load timezone, using default: %v", err)
	} else {
		time.Local = loc
		log.Printf("🌏 Timezone set to: %s", loc.String())
	}

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
	musicRepo := persistence.NewMusicRepository(db)
	eventRepo := persistence.NewEventRepository(db)
	bookingRepo := persistence.NewBookingRepository(db)
	favoriteRepo := persistence.NewFavoriteRepository(db)
	bandRepo := persistence.NewBandRepository(db)
	settingRepo := persistence.NewSettingRepository(db)
	notificationRepo := persistence.NewNotificationRepository(db)
	deviceTokenRepo := persistence.NewDeviceTokenRepository(db)

	// Initialize Firebase Cloud Messaging service
	fcmService, err := firebase.NewFCMService(config.FirebaseCredentialPath, deviceTokenRepo)
	if err != nil {
		log.Fatalf("❌ Failed to initialize FCM service: %v", err)
	}

	// Initialize use cases (Application layer)
	authUseCase := usecase.NewAuthUseCase(userRepo, config.SecretKey, config.GoogleClientID)
	musicUseCase := usecase.NewMusicUseCase(musicRepo)
	eventUseCase := usecase.NewEventUseCase(eventRepo, musicRepo, notificationRepo)
	bookingUseCase := usecase.NewBookingUseCase(bookingRepo, eventRepo)
	favoriteUseCase := usecase.NewFavoriteUseCase(favoriteRepo, musicRepo)
	bandUseCase := usecase.NewBandUseCase(bandRepo, musicRepo)
	settingUseCase := usecase.NewSettingUseCase(settingRepo)
	notificationUseCase := usecase.NewNotificationUseCase(notificationRepo, fcmService)

	// Initialize handlers (Interface layer)
	authHandler := handler.NewAuthHandler(authUseCase)
	musicHandler := handler.NewMusicHandler(musicUseCase)
	eventHandler := handler.NewEventHandler(eventUseCase)
	bookingHandler := handler.NewBookingHandler(bookingUseCase)
	favoriteHandler := handler.NewFavoriteHandler(favoriteUseCase)
	bandHandler := handler.NewBandHandler(bandUseCase)
	settingHandler := handler.NewSettingHandler(settingUseCase)
	notificationHandler := handler.NewNotificationHandler(notificationUseCase)
	deviceTokenHandler := handler.NewDeviceTokenHandler(deviceTokenRepo)

	// Setup router
	r := router.NewRouter(authHandler, musicHandler, eventHandler, bookingHandler, favoriteHandler, bandHandler, settingHandler, notificationHandler, deviceTokenHandler, authUseCase)
	engine := r.Setup()

	// Start server
	log.Printf("🚀 Server starting on port %s", config.Port)
	if err := engine.Run(":" + config.Port); err != nil {
		log.Fatalf("❌ failed to start server: %v", err)
	}
}
