// pkg/database/postgres.go - อัพเดทให้รองรับ Laravel-style migrations
package database

import (
	"fmt"
	"time"

	"go-clean-gin/config"
	"go-clean-gin/internal/migrations"
	"go-clean-gin/internal/seeders"
	"go-clean-gin/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// NewPostgresDB creates a new PostgreSQL database connection
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

	// Configure GORM logger based on config
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

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		DisableForeignKeyConstraintWhenMigrating: false,
		CreateBatchSize:                          1000,
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
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, err
	}

	logger.Info("Successfully connected to PostgreSQL database",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.String("database", cfg.Name),
		zap.Int("max_idle_conns", cfg.MaxIdleConns),
		zap.Int("max_open_conns", cfg.MaxOpenConns))

	return db, nil
}

// RunMigrations runs database migrations using Laravel-style migration system
func RunMigrations(db *gorm.DB) error {
	logger.Info("Starting Laravel-style migrations...")

	// Create migration manager
	migrationManager := migrations.NewMigrationManager(db)
	migrations.SetGlobalManager(migrationManager)

	// Run migrations
	if err := migrationManager.RunMigrations(); err != nil {
		logger.Error("Failed to run migrations", zap.Error(err))
		return err
	}

	logger.Info("Laravel-style migrations completed successfully")
	return nil
}

// RollbackMigrations rolls back the specified number of migrations
func RollbackMigrations(db *gorm.DB, count int) error {
	logger.Info("Starting migration rollback...", zap.Int("count", count))

	// Create migration manager
	migrationManager := migrations.NewMigrationManager(db)
	migrations.SetGlobalManager(migrationManager)

	// Rollback migrations
	if err := migrationManager.RollbackMigrations(count); err != nil {
		logger.Error("Failed to rollback migrations", zap.Error(err))
		return err
	}

	logger.Info("Migration rollback completed successfully")
	return nil
}

// GetMigrationStatus returns the current migration status
func GetMigrationStatus(db *gorm.DB) error {
	// Create migration manager
	migrationManager := migrations.NewMigrationManager(db)
	migrations.SetGlobalManager(migrationManager)

	// Get migration status
	if err := migrationManager.GetMigrationStatus(); err != nil {
		logger.Error("Failed to get migration status", zap.Error(err))
		return err
	}

	return nil
}

// SeedData seeds the database with initial data using Laravel-style seeders
func SeedData(db *gorm.DB, seederName string) error {
	logger.Info("Starting Laravel-style database seeding...")

	// Create seeder manager
	seederManager := seeders.NewSeederManager(db)
	seeders.SetGlobalSeederManager(seederManager)

	// Run seeders
	if err := seederManager.RunSeeders(seederName); err != nil {
		logger.Error("Failed to run seeders", zap.Error(err))
		return err
	}

	logger.Info("Laravel-style database seeding completed successfully")
	return nil
}

// RunSpecificSeeder runs a specific seeder
func RunSpecificSeeder(db *gorm.DB, seederName string) error {
	logger.Info("Running specific seeder...", zap.String("seeder", seederName))

	// Create seeder manager
	seederManager := seeders.NewSeederManager(db)
	seeders.SetGlobalSeederManager(seederManager)

	// Run specific seeder
	if err := seederManager.RunSpecificSeeder(seederName); err != nil {
		logger.Error("Failed to run specific seeder", zap.Error(err))
		return err
	}

	logger.Info("Specific seeder completed successfully")
	return nil
}

// ListSeeders lists all registered seeders
func ListSeeders(db *gorm.DB) error {
	// Create seeder manager
	seederManager := seeders.NewSeederManager(db)
	seeders.SetGlobalSeederManager(seederManager)

	// List seeders
	seederManager.ListSeeders()
	return nil
}

// HealthCheck checks the database connection health
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

// GetDatabaseStats returns database connection statistics
func GetDatabaseStats(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	stats := sqlDB.Stats()

	logger.Info("Database Connection Statistics",
		zap.Int("max_open_connections", stats.MaxOpenConnections),
		zap.Int("open_connections", stats.OpenConnections),
		zap.Int("in_use", stats.InUse),
		zap.Int("idle", stats.Idle),
		zap.Int64("wait_count", stats.WaitCount),
		zap.Duration("wait_duration", stats.WaitDuration),
		zap.Int64("max_idle_closed", stats.MaxIdleClosed),
		zap.Int64("max_idle_time_closed", stats.MaxIdleTimeClosed),
		zap.Int64("max_lifetime_closed", stats.MaxLifetimeClosed))

	return nil
}
