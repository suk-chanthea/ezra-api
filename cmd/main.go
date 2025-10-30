package main

import (
	"log"
	"os"
	"time"

	"github.com/suk-chanthea/ezra/infrastructure/email"
	"github.com/suk-chanthea/ezra/infrastructure/firebase"
	"github.com/suk-chanthea/ezra/infrastructure/payment"
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
	PaywayMerchantID       string
	PaywayAPIKey           string
	PaywayAPIUsername      string
	PaywayBaseURL          string
	PaywayReturnURL        string
	PaywayContinueURL      string
	PaywayCallbackURL      string
	SMTPHost               string
	SMTPPort               string
	SMTPUsername           string
	SMTPPassword           string
	SMTPFrom               string
	OTPExpiry              int // in minutes
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

	// Payway configuration
	paywayMerchantID := os.Getenv("PAYWAY_MERCHANT_ID")
	if paywayMerchantID == "" {
		paywayMerchantID = "your_merchant_id" // Replace with actual merchant ID
	}

	paywayAPIKey := os.Getenv("PAYWAY_API_KEY")
	if paywayAPIKey == "" {
		paywayAPIKey = "your_api_key" // Replace with actual API key
	}

	paywayAPIUsername := os.Getenv("PAYWAY_API_USERNAME")
	if paywayAPIUsername == "" {
		paywayAPIUsername = "your_api_username" // Replace with actual username
	}

	paywayBaseURL := os.Getenv("PAYWAY_BASE_URL")
	if paywayBaseURL == "" {
		// Use sandbox by default
		paywayBaseURL = "https://api-sandbox.payway.com.kh"
	}

	paywayReturnURL := os.Getenv("PAYWAY_RETURN_URL")
	if paywayReturnURL == "" {
		paywayReturnURL = "http://localhost:3000/donation/complete" // Frontend URL
	}

	paywayContinueURL := os.Getenv("PAYWAY_CONTINUE_URL")
	if paywayContinueURL == "" {
		paywayContinueURL = "http://localhost:3000/donation/success" // Frontend URL
	}

	paywayCallbackURL := os.Getenv("PAYWAY_CALLBACK_URL")
	if paywayCallbackURL == "" {
		paywayCallbackURL = "http://localhost:8080/webhooks/payway" // Backend webhook URL
	}

	// SMTP configuration for sending OTP emails
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com" // Default Gmail SMTP
	}

	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587" // Default SMTP port for TLS
	}

	smtpUsername := os.Getenv("SMTP_USERNAME")
	// Required: Set via environment variable (your Gmail address)

	smtpPassword := os.Getenv("SMTP_PASSWORD")
	// Required: Set via environment variable (Gmail App Password)

	smtpFrom := os.Getenv("SMTP_FROM")
	if smtpFrom == "" {
		smtpFrom = smtpUsername // Use same as username if not specified
	}

	otpExpiry := 10 // Default 10 minutes
	if otpExpiryEnv := os.Getenv("OTP_EXPIRY_MINUTES"); otpExpiryEnv != "" {
		// Parse OTP expiry from env if provided
		if exp, err := time.ParseDuration(otpExpiryEnv + "m"); err == nil {
			otpExpiry = int(exp.Minutes())
		}
	}

	return &Config{
		Port:                   port,
		PostgresURL:            pg,
		SecretKey:              secret,
		GoogleClientID:         googleClientID,
		FirebaseCredentialPath: firebaseCredPath,
		PaywayMerchantID:       paywayMerchantID,
		PaywayAPIKey:           paywayAPIKey,
		PaywayAPIUsername:      paywayAPIUsername,
		PaywayBaseURL:          paywayBaseURL,
		PaywayReturnURL:        paywayReturnURL,
		PaywayContinueURL:      paywayContinueURL,
		PaywayCallbackURL:      paywayCallbackURL,
		SMTPHost:               smtpHost,
		SMTPPort:               smtpPort,
		SMTPUsername:           smtpUsername,
		SMTPPassword:           smtpPassword,
		SMTPFrom:               smtpFrom,
		OTPExpiry:              otpExpiry,
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
	donationRepo := persistence.NewDonationRepository(db)
	supporterRepo := persistence.NewSupporterRepository(db)
	churchRepo := persistence.NewChurchRepository(db)
	otpRepo := persistence.NewOTPRepository(db)

	// Initialize Firebase Cloud Messaging service
	fcmService, err := firebase.NewFCMService(config.FirebaseCredentialPath, deviceTokenRepo)
	if err != nil {
		log.Fatalf("❌ Failed to initialize FCM service: %v", err)
	}

	// Initialize Payway service
	paywayConfig := &payment.PaywayConfig{
		MerchantID:   config.PaywayMerchantID,
		APIKey:       config.PaywayAPIKey,
		APIUsername:  config.PaywayAPIUsername,
		BaseURL:      config.PaywayBaseURL,
		ReturnURL:    config.PaywayReturnURL,
		ContinueURL:  config.PaywayContinueURL,
		CallbackURL:  config.PaywayCallbackURL,
	}
	paywayService := payment.NewPaywayService(paywayConfig)
	log.Println("✅ Payway service initialized")

	// Initialize Email service for OTP
	emailService := email.NewSMTPEmailService(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPUsername,
		config.SMTPPassword,
		config.SMTPFrom,
	)
	if config.SMTPUsername == "" || config.SMTPPassword == "" {
		log.Println("⚠️  Warning: SMTP credentials not set. OTP emails will not be sent.")
	} else {
		log.Println("✅ Email service initialized")
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
	donationUseCase := usecase.NewDonationUseCase(donationRepo, userRepo, eventRepo, paywayService)
	supporterUseCase := usecase.NewSupporterUseCase(supporterRepo, donationRepo)
	churchUseCase := usecase.NewChurchUseCase(churchRepo, userRepo)
	otpUseCase := usecase.NewOTPUseCase(otpRepo, userRepo, emailService, config.OTPExpiry)

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
	donationHandler := handler.NewDonationHandler(donationUseCase)
	supporterHandler := handler.NewSupporterHandler(supporterUseCase)
	churchHandler := handler.NewChurchHandler(churchUseCase)
	otpHandler := handler.NewOTPHandler(otpUseCase)

	// Setup router
	r := router.NewRouter(authHandler, musicHandler, eventHandler, bookingHandler, favoriteHandler, bandHandler, settingHandler, notificationHandler, deviceTokenHandler, donationHandler, supporterHandler, churchHandler, otpHandler, authUseCase)
	engine := r.Setup()

	// Start server
	log.Printf("🚀 Server starting on port %s", config.Port)
	if err := engine.Run(":" + config.Port); err != nil {
		log.Fatalf("❌ failed to start server: %v", err)
	}
}
