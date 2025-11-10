package database

import (
	"time"
	"log"
	"strings"

	"github.com/suk-chanthea/ezra/config"
	"github.com/suk-chanthea/ezra/infrastructure/persistence"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// NewPostgresDB establishes a new PostgreSQL connection using GORM.
func NewPostgresDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	gormCfg := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), gormCfg)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Basic pooling configuration with sensible defaults.
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}

// AutoMigrate runs database schema migrations for all persistence models.
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&persistence.UserModel{},
		&persistence.MusicModel{},
		&persistence.BandModel{},
		&persistence.BandMusicModel{},
		&persistence.EventModel{},
		&persistence.EventMusicModel{},
		&persistence.BookingModel{},
		&persistence.FavoriteModel{},
		&persistence.NotificationModel{},
		&persistence.DeviceTokenModel{},
		&persistence.OTPModel{},
		&persistence.DonationModel{},
		&persistence.SupporterModel{},
		&persistence.SettingModel{},
		&persistence.ChurchModel{},
	)

	// Gracefully ignore missing constraint errors
	if err != nil {
		if strings.Contains(err.Error(), `constraint "uni_supporters_email"`) {
			log.Println("⚠️  Ignoring missing constraint 'uni_supporters_email' on 'supporters' table.")
			return nil
		}
		return err
	}

	return nil
}
