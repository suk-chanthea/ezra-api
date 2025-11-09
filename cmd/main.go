package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/suk-chanthea/ezra/config"
	"github.com/suk-chanthea/ezra/infrastructure/database"
	"github.com/suk-chanthea/ezra/infrastructure/email"
	"github.com/suk-chanthea/ezra/infrastructure/firebase"
	"github.com/suk-chanthea/ezra/infrastructure/payment"
	"github.com/suk-chanthea/ezra/infrastructure/persistence"
	"github.com/suk-chanthea/ezra/interface/http/handler"
	"github.com/suk-chanthea/ezra/interface/http/router"
	"github.com/suk-chanthea/ezra/usecase"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("failed to start application: %v", err)
	}
}

func run() error {
	log.Println("📋 Loading configuration...")
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	log.Println("🗄️  Connecting to database...")
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	log.Println("🛠️  Running database migrations...")
	if err := database.AutoMigrate(db); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	// Repositories
	log.Println("🏗️  Initializing repositories...")
	userRepo := persistence.NewUserRepository(db)
	musicRepo := persistence.NewMusicRepository(db)
	eventRepo := persistence.NewEventRepository(db)
	bookingRepo := persistence.NewBookingRepository(db)
	favoriteRepo := persistence.NewFavoriteRepository(db)
	bandRepo := persistence.NewBandRepository(db)
	settingRepo := persistence.NewSettingRepository(db)
	notificationRepo := persistence.NewNotificationRepository(db)
	deviceTokenRepo := persistence.NewDeviceTokenRepository(db)
	otpRepo := persistence.NewOTPRepository(db)
	donationRepo := persistence.NewDonationRepository(db)
	supporterRepo := persistence.NewSupporterRepository(db)
	churchRepo := persistence.NewChurchRepository(db)

	// External services
	var emailService email.EmailService
	if cfg.Email.Enabled {
		log.Println("📧 Initializing SMTP email service...")
		emailService = email.NewSMTPEmailService(
			cfg.Email.Host,
			cfg.Email.Port,
			cfg.Email.Username,
			cfg.Email.Password,
			cfg.Email.From,
			cfg.Email.Secure,
		)
	} else {
		emailService = email.NewDummyEmailService()
		log.Println("ℹ️  Email service disabled; using dummy implementation.")
	}

	var fcmService firebase.FCMService
	if cfg.Firebase.Enabled && cfg.Firebase.CredentialsPath != "" {
		log.Println("🔥 Initializing Firebase Cloud Messaging...")
		fcmService, err = firebase.NewFCMService(cfg.Firebase.CredentialsPath, deviceTokenRepo)
		if err != nil {
			log.Printf("⚠️  Failed to initialize Firebase: %v. Falling back to dummy service.", err)
			fcmService = firebase.NewDummyFCMService()
		}
	} else {
		fcmService = firebase.NewDummyFCMService()
		log.Println("ℹ️  Firebase disabled; using dummy implementation.")
	}

	var paywayService payment.PaywayService
	if cfg.Payway.Enabled {
		log.Println("💳 Initializing PayWay service...")
		paywayService = payment.NewPaywayService(&payment.PaywayConfig{
			MerchantID:  cfg.Payway.MerchantID,
			APIKey:      cfg.Payway.APIKey,
			APIUsername: cfg.Payway.APIUsername,
			BaseURL:     cfg.Payway.BaseURL,
			ReturnURL:   cfg.Payway.ReturnURL,
			ContinueURL: cfg.Payway.ContinueURL,
			CallbackURL: cfg.Payway.CallbackURL,
		})
	} else {
		paywayService = payment.NewDummyPayWayService()
		log.Println("ℹ️  PayWay disabled; using dummy implementation.")
	}

	// Use cases
	log.Println("⚙️  Initializing use cases...")
	authUseCase := usecase.NewAuthUseCase(userRepo, otpRepo, cfg.JWT.Secret, cfg.OAuth.GoogleClientID)
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
	otpUseCase := usecase.NewOTPUseCase(otpRepo, userRepo, emailService, cfg.JWT.TokenExpiry)

	// Handlers
	log.Println("🎯 Wiring HTTP handlers...")
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

	log.Println("🛣️  Configuring router...")
	r := router.NewRouter(
		authHandler,
		musicHandler,
		eventHandler,
		bookingHandler,
		favoriteHandler,
		bandHandler,
		settingHandler,
		notificationHandler,
		deviceTokenHandler,
		donationHandler,
		supporterHandler,
		churchHandler,
		otpHandler,
		authUseCase,
	)
	engine := r.Setup()

	addr := ":" + cfg.App.Port
	if cfg.App.Port == "" {
		addr = ":8080"
	}

	srv := &http.Server{
		Addr:           addr,
		Handler:        engine,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("🚀 Server starting on %s (%s)", addr, cfg.App.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	log.Println("✅ Shutdown complete.")
	return nil
}
