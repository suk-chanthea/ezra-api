package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config aggregates all application configuration sections.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	OAuth    OAuthConfig
	Email    EmailConfig
	Firebase FirebaseConfig
	Payway   PaywayConfig
}

// AppConfig contains metadata about the running application.
type AppConfig struct {
	Environment string
	Version     string
	Port        string
}

// DatabaseConfig wraps PostgreSQL connection parameters.
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Timezone string
}

// DSN builds a GORM-compatible PostgreSQL DSN.
func (c DatabaseConfig) DSN() string {
	timezone := c.Timezone
	if timezone == "" {
		timezone = "UTC"
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Name,
		c.SSLMode,
		timezone,
	)
}

// JWTConfig stores JWT related settings.
type JWTConfig struct {
	Secret      string
	TokenExpiry int // minutes
}

// OAuthConfig stores OAuth provider configuration.
type OAuthConfig struct {
	GoogleClientID string
}

// EmailConfig configures SMTP credentials.
type EmailConfig struct {
	Enabled  bool
	Host     string
	Port     string
	Username string
	Password string
	From     string
	Secure   string // starttls (default), ssl, plain
}

// FirebaseConfig configures Firebase Cloud Messaging.
type FirebaseConfig struct {
	Enabled         bool
	CredentialsPath string
}

// PaywayConfig contains PayWay payment gateway credentials.
type PaywayConfig struct {
	Enabled     bool
	MerchantID  string
	APIKey      string
	APIUsername string
	BaseURL     string
	ReturnURL   string
	ContinueURL string
	CallbackURL string
}

// Load reads configuration from the environment with sensible defaults.
func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			Version:     getEnv("APP_VERSION", "1.0.0"),
			Port:        getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "ezradb"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Timezone: getEnv("DB_TIMEZONE", "UTC"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "super-secret-key"),
			TokenExpiry: getEnvAsInt("JWT_TOKEN_EXPIRY", 10),
		},
		OAuth: OAuthConfig{
			GoogleClientID: getEnv("GOOGLE_CLIENT_ID", ""),
		},
		Email: EmailConfig{
			Enabled:  getEnvAsBool("EMAIL_ENABLED", false),
			Host:     getEnv("SMTP_HOST", ""),
			Port:     getEnv("SMTP_PORT", "587"),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", ""),
			Secure:   getEnv("SMTP_SECURE", "starttls"),
		},
		Firebase: FirebaseConfig{
			Enabled:         getEnvAsBool("FIREBASE_ENABLED", false),
			CredentialsPath: getEnv("FIREBASE_CREDENTIALS_PATH", ""),
		},
		Payway: PaywayConfig{
			Enabled:     getEnvAsBool("PAYWAY_ENABLED", false),
			MerchantID:  getEnv("PAYWAY_MERCHANT_ID", ""),
			APIKey:      getEnv("PAYWAY_API_KEY", ""),
			APIUsername: getEnv("PAYWAY_API_USERNAME", ""),
			BaseURL:     getEnv("PAYWAY_BASE_URL", ""),
			ReturnURL:   getEnv("PAYWAY_RETURN_URL", ""),
			ContinueURL: getEnv("PAYWAY_CONTINUE_URL", ""),
			CallbackURL: getEnv("PAYWAY_CALLBACK_URL", ""),
		},
	}

	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		lower := strings.ToLower(value)
		switch lower {
		case "true", "1", "t", "yes", "y":
			return true
		case "false", "0", "f", "no", "n":
			return false
		default:
			return fallback
		}
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	}
	return fallback
}
