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
	"github.com/suk-chanthea/ezra/infrastructure/cache"
	"github.com/suk-chanthea/ezra/infrastructure/database"
	"github.com/suk-chanthea/ezra/infrastructure/email"
	"github.com/suk-chanthea/ezra/infrastructure/firebase"
	"github.com/suk-chanthea/ezra/infrastructure/logger"
	"github.com/suk-chanthea/ezra/infrastructure/payway"
	"github.com/suk-chanthea/ezra/infrastructure/persistence"
	"github.com/suk-chanthea/ezra/interface/http/handler"
	"github.com/suk-chanthea/ezra/interface/http/router"
	"github.com/suk-chanthea/ezra/usecase"

	"go.uber.org/zap"
)

func main() {
	// Run application
	if err := run(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}

func run() error {
	// ==================== Configuration ====================
	log.Println("📋 Loading configuration...")
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// ==================== Logger ====================
	log.Println("📝 Initializing logger...")
	if err := logger.InitLogger(cfg.App.Environment); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	appLogger := logger.GetLogger()
	appLogger.Info("Logger initialized",
		zap.String("environment", cfg.App.Environment),
		zap.String("version", cfg.App.Version),
	)

	// ==================== Database ====================
	appLogger.Info("🗄️  Connecting to database...")
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	appLogger.Info("✅ Database connected")

	// Auto-migrate if in development
	if cfg.App.Environment == "development" {
		appLogger.Info("Running database migrations...")
		if err := database.AutoMigrate(db); err != nil {
			appLogger.Warn("Failed to run migrations", zap.Error(err))
		}
	}

	// ==================== Cache (Optional) ====================
	var cacheService cache.Cache
	if cfg.Redis.Enabled {
		appLogger.Info("💾 Connecting to Redis...")
		redisClient, err := cache.NewRedisClient(&cfg.Redis)
		if err != nil {
			appLogger.Warn("Failed to connect to Redis, continuing without cache", zap.Error(err))
			cacheService = cache.NewNoOpCache() // Fallback to no-op cache
		} else {
			cacheService = cache.NewRedisCache(redisClient)
			appLogger.Info("✅ Redis connected")
		}
	} else {
		cacheService = cache.NewNoOpCache()
		appLogger.Info("ℹ️  Cache disabled")
	}

	// ==================== Repositories ====================
	appLogger.Info("🏗️  Initializing repositories...")
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
	appLogger.Info("✅ Repositories initialized")

	// ==================== External Services ====================
	
	// Firebase (Optional)
	var fcmService firebase.FCMService
	if cfg.Firebase.Enabled && cfg.Firebase.CredentialsPath != "" {
		appLogger.Info("🔥 Initializing Firebase...")
		fcmService = firebase.NewFCMService(cfg.Firebase.CredentialsPath, deviceTokenRepo)
		appLogger.Info("✅ Firebase initialized")
	} else {
		appLogger.Info("ℹ️  Firebase disabled, using dummy service")
		fcmService = firebase.NewDummyFCMService()
	}

	// Email Service (Optional)
	var emailService email.EmailService
	if cfg.Email.Enabled {
		appLogger.Info("📧 Initializing email service...")
		emailService = email.NewSMTPEmailService(&cfg.Email)
		appLogger.Info("✅ Email service initialized")
	} else {
		appLogger.Info("ℹ️  Email service disabled")
		emailService = email.NewDummyEmailService()
	}

	// PayWay Service (Optional)
	var paywayService payway.PayWayService
	if cfg.PayWay.Enabled {
		appLogger.Info("💳 Initializing PayWay service...")
		paywayService = payway.NewPayWayService(&cfg.PayWay)
		appLogger.Info("✅ PayWay service initialized")
	} else {
		appLogger.Info("ℹ️  PayWay service disabled")
		paywayService = payway.NewDummyPayWayService()
	}

	// ==================== Use Cases ====================
	appLogger.Info("⚙️  Initializing use cases...")
	authUseCase := usecase.NewAuthUseCase(
		userRepo,
		otpRepo,
		&cfg.JWT,
		&cfg.OAuth,
		cacheService,
		appLogger,
	)
	musicUseCase := usecase.NewMusicUseCase(
		musicRepo,
		cacheService,
		appLogger,
	)
	eventUseCase := usecase.NewEventUseCase(
		eventRepo,
		musicRepo,
		notificationRepo,
		appLogger,
	)
	bookingUseCase := usecase.NewBookingUseCase(
		bookingRepo,
		eventRepo,
		appLogger,
	)
	favoriteUseCase := usecase.NewFavoriteUseCase(
		favoriteRepo,
		musicRepo,
		appLogger,
	)
	bandUseCase := usecase.NewBandUseCase(
		bandRepo,
		musicRepo,
		appLogger,
	)
	settingUseCase := usecase.NewSettingUseCase(
		settingRepo,
		appLogger,
	)
	notificationUseCase := usecase.NewNotificationUseCase(
		notificationRepo,
		fcmService,
		appLogger,
	)
	donationUseCase := usecase.NewDonationUseCase(
		donationRepo,
		userRepo,
		eventRepo,
		paywayService,
		appLogger,
	)
	supporterUseCase := usecase.NewSupporterUseCase(
		supporterRepo,
		donationRepo,
		appLogger,
	)
	churchUseCase := usecase.NewChurchUseCase(
		churchRepo,
		userRepo,
		appLogger,
	)
	otpUseCase := usecase.NewOTPUseCase(
		otpRepo,
		userRepo,
		emailService,
		cfg.JWT.TokenExpiry,
		appLogger,
	)
	appLogger.Info("✅ Use cases initialized")

	// ==================== Handlers ====================
	appLogger.Info("🎯 Initializing handlers...")
	authHandler := handler.NewAuthHandler(authUseCase, appLogger)
	musicHandler := handler.NewMusicHandler(musicUseCase, appLogger)
	eventHandler := handler.NewEventHandler(eventUseCase, appLogger)
	bookingHandler := handler.NewBookingHandler(bookingUseCase, appLogger)
	favoriteHandler := handler.NewFavoriteHandler(favoriteUseCase, appLogger)
	bandHandler := handler.NewBandHandler(bandUseCase, appLogger)
	settingHandler := handler.NewSettingHandler(settingUseCase, appLogger)
	notificationHandler := handler.NewNotificationHandler(notificationUseCase, appLogger)
	deviceTokenHandler := handler.NewDeviceTokenHandler(deviceTokenRepo, appLogger)
	donationHandler := handler.NewDonationHandler(donationUseCase, appLogger)
	supporterHandler := handler.NewSupporterHandler(supporterUseCase, appLogger)
	churchHandler := handler.NewChurchHandler(churchUseCase, appLogger)
	otpHandler := handler.NewOTPHandler(otpUseCase, appLogger)
	appLogger.Info("✅ Handlers initialized")

	// ==================== Router ====================
	appLogger.Info("🛣️  Setting up routes...")
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
		cfg,
		appLogger,
	)
	engine := r.Setup()
	appLogger.Info("✅ Routes configured")

	// ==================== HTTP Server ====================
	srv := &http.Server{
		Addr:           ":" + cfg.App.Port,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in goroutine
	go func() {
		appLogger.Info("🚀 Server starting",
			zap.String("port", cfg.App.Port),
			zap.String("environment", cfg.App.Environment),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// ==================== Graceful Shutdown ====================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", zap.Error(err))
		return err
	}

	// Close database connection
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	appLogger.Info("✅ Server gracefully stopped")
	return nil
}