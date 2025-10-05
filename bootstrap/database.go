package bootstrap

import (
    "context"
    "github.com/go-redis/redis/v8"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
)

func NewDatabase(config *Config) *gorm.DB {
    // Use PostgresURL from config
    db, err := gorm.Open(postgres.Open(config.PostgresURL), &gorm.Config{})
    if err != nil {
        log.Fatalf("❌ failed to connect database: %v", err)
    }

    log.Println("✅ PostgreSQL connected")
    return db
}

func ConnectPostgres(dsn string) *gorm.DB {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect postgres:", err)
    }
    return db
}

func ConnectRedis(addr, password string) *redis.Client {
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       0,
    })
    if err := rdb.Ping(context.Background()).Err(); err != nil {
        log.Fatal("failed to connect redis:", err)
    }
    return rdb
}
