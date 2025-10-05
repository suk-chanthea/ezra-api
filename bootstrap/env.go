package bootstrap

import (
    "os"
)

type Config struct {
    Port         string
    PostgresURL  string
    RedisAddress string
    RedisPass    string
	SecretKey 	 string
}

func LoadConfig() *Config {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    pg := os.Getenv("POSTGRES_URL")
    if pg == "" {
        pg = "postgres://postgres:secret@postgres:5432/ezradb?sslmode=disable"
    }

    redisAddr := os.Getenv("REDIS_ADDRESS")
    if redisAddr == "" {
        redisAddr = "redis:6379"
    }

	scKey := os.Getenv("SECRET")
	if scKey == "" {
		scKey = "paracletus"
	}

    redisPass := os.Getenv("REDIS_PASSWORD")
    return &Config{
        Port:         port,
        PostgresURL:  pg,
        RedisAddress: redisAddr,
        RedisPass:    redisPass,
		SecretKey: 	  scKey,
    }
}
