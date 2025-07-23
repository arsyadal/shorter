package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"shorter-backend/models"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
	Ctx = context.Background()
)

func ConnectDatabase() {
	// Database connection
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "shorter_db")
	port := getEnv("DB_PORT", "5432")
	sslmode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		host, user, password, dbname, port, sslmode)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = database

	// Auto migrate the schema
	err = DB.AutoMigrate(&models.URL{}, &models.Click{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected successfully")
}

func ConnectRedis() {
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")

	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0,
	})

	// Test connection
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		log.Println("Continuing without Redis cache...")
		RDB = nil
		return
	}

	log.Println("Redis connected successfully")
}

func CacheSet(key string, value string, expiration time.Duration) error {
	if RDB == nil {
		return nil // Skip if Redis is not available
	}
	return RDB.Set(Ctx, key, value, expiration).Err()
}

func CacheGet(key string) (string, error) {
	if RDB == nil {
		return "", redis.Nil // Return redis.Nil to indicate cache miss
	}
	return RDB.Get(Ctx, key).Result()
}

func CacheDelete(key string) error {
	if RDB == nil {
		return nil
	}
	return RDB.Del(Ctx, key).Err()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 