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

	// Configure GORM logger based on config 🆕
	var logLevel gormLogger.LogLevel
	switch cfg.LogLevel {
	case "debug":
		logLevel = gormLogger.Info
	case "info":
		logLevel = gormLogger.Warn
	case "warn":
		logLevel = gormLogger.Error
	default:
		logLevel = gormLogger.Error
	}

	// Configure GORM - ปรับปรุงจากเดิม 🔧
	gormConfig := &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLevel), // 🆕 ใช้ dynamic log level
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		DisableForeignKeyConstraintWhenMigrating: false, // 🆕 เพิ่ม FK support
		CreateBatchSize:                          1000,  // 🆕 ปรับปรุง performance
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

	// Configure connection pool - ใช้ settings จาก config 🆕
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, err
	}

	// ปรับปรุง logging ให้มีข้อมูลมากขึ้น 🆕
	logger.Info("Successfully connected to PostgreSQL database",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Name),
		zap.Int("max_idle_conns", cfg.MaxIdleConns),
		zap.Int("max_open_conns", cfg.MaxOpenConns))

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
	} else {
		logger.Info("Admin user already exists, skipping creation") // 🆕 เพิ่ม log
	}

	logger.Info("Database seeding completed")
	return nil
}

// 🆕 เพิ่ม utility functions ใหม่

// HealthCheck - ตรวจสอบการเชื่อมต่อ database
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
