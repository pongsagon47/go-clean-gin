package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Log      LogConfig
	Env      string
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	LogLevel        string // 🆕 เพิ่มใหม่ - สำหรับ GORM logging
	MaxIdleConns    int    // 🆕 เพิ่มใหม่ - connection pool
	MaxOpenConns    int    // 🆕 เพิ่มใหม่ - connection pool
	ConnMaxLifetime int    // 🆕 เพิ่มใหม่ - connection lifetime (minutes)
}

type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

type LogConfig struct {
	Level  string
	Format string
}

func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvAsInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "password"),
			Name:            getEnv("DB_NAME", "go_clean_gin"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			LogLevel:        getEnv("DB_LOG_LEVEL", "warn"),          // 🆕 เพิ่มใหม่
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),    // 🆕 เพิ่มใหม่
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 100),   // 🆕 เพิ่มใหม่
			ConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME", 60), // 🆕 เพิ่มใหม่ (60 นาที)
		},
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvAsInt("SERVER_PORT", 8080), // 👆 เก็บ 8080 ตามเดิม
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
			ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		Env: getEnv("ENV", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
