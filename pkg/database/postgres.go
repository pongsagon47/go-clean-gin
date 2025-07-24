package database

import (
	"fmt"
	"time"

	"go-clean-gin/config"
	"go-clean-gin/internal/entity"
	"go-clean-gin/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode,
	)

	// Configure GORM logger
	gormConfig := &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	// Get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to get underlying sql.DB", zap.Error(err))
		return nil, err
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, err
	}

	logger.Info("Successfully connected to PostgreSQL database")
	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	logger.Info("Running database migrations...")

	err := db.AutoMigrate(
		&entity.User{},
		&entity.Product{},
	)

	if err != nil {
		logger.Error("Failed to run migrations", zap.Error(err))
		return err
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

func SeedData(db *gorm.DB) error {
	logger.Info("Seeding database...")

	// Check if admin user already exists
	var count int64
	db.Model(&entity.User{}).Where("email = ?", "admin@example.com").Count(&count)

	if count == 0 {
		// Create admin user (you should hash this password in production)
		admin := entity.User{
			Email:     "admin@example.com",
			Username:  "admin",
			Password:  "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password: password
			FirstName: "Admin",
			LastName:  "User",
			IsActive:  true,
		}

		if err := db.Create(&admin).Error; err != nil {
			logger.Error("Failed to create admin user", zap.Error(err))
			return err
		}

		logger.Info("Admin user created successfully")
	}

	logger.Info("Database seeding completed")
	return nil
}
