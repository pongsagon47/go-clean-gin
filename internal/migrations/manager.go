// internal/migrations/manager.go - Migration Manager สำหรับระบบ Laravel-style
package migrations

import (
	"fmt"
	"sort"
	"time"

	"go-clean-gin/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Migration interface ที่แต่ละไฟล์ต้อง implement
type Migration interface {
	Up(db *gorm.DB) error
	Down(db *gorm.DB) error
	Version() string
	Description() string
}

// MigrationRecord represents migration history in database
type MigrationRecord struct {
	ID          uint      `gorm:"primaryKey"`
	Version     string    `gorm:"uniqueIndex;not null"`
	Description string    `gorm:"not null"`
	AppliedAt   time.Time `gorm:"not null"`
}

// MigrationManager จัดการ migrations
type MigrationManager struct {
	db         *gorm.DB
	migrations map[string]Migration
}

// Global migration manager instance
var globalManager *MigrationManager
var registeredMigrations []Migration

// NewMigrationManager สร้าง manager ใหม่
func NewMigrationManager(db *gorm.DB) *MigrationManager {
	manager := &MigrationManager{
		db:         db,
		migrations: make(map[string]Migration),
	}

	// Register all migrations that were registered during init()
	for _, migration := range registeredMigrations {
		manager.RegisterMigration(migration)
	}

	return manager
}

// SetGlobalManager ตั้งค่า global manager
func SetGlobalManager(manager *MigrationManager) {
	globalManager = manager
}

// Register ฟังก์ชันสำหรับให้แต่ละไฟล์เรียกใช้ใน init()
func Register(migration Migration) {
	registeredMigrations = append(registeredMigrations, migration)

	// ถ้ามี global manager แล้ว ให้ register เลย
	if globalManager != nil {
		globalManager.RegisterMigration(migration)
	}
}

// RegisterMigration ลงทะเบียน migration
func (mm *MigrationManager) RegisterMigration(migration Migration) {
	mm.migrations[migration.Version()] = migration
}

// RunMigrations รัน migrations ที่ยังไม่ได้ apply
func (mm *MigrationManager) RunMigrations() error {
	// Create migrations table if not exists
	if err := mm.db.AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	var appliedRecords []MigrationRecord
	if err := mm.db.Find(&appliedRecords).Error; err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedMap := make(map[string]bool)
	for _, record := range appliedRecords {
		appliedMap[record.Version] = true
	}

	// Sort migrations by version
	var versions []string
	for version := range mm.migrations {
		versions = append(versions, version)
	}
	sort.Strings(versions)

	// Run pending migrations
	pendingCount := 0
	for _, version := range versions {
		if appliedMap[version] {
			logger.Debug("Migration already applied",
				zap.String("version", version))
			continue
		}

		pendingCount++
		migration := mm.migrations[version]

		logger.Info("Running migration",
			zap.String("version", version),
			zap.String("description", migration.Description()))

		if err := mm.runSingleMigration(migration); err != nil {
			return fmt.Errorf("migration %s failed: %w", version, err)
		}

		logger.Info("Migration completed",
			zap.String("version", version))
	}

	if pendingCount == 0 {
		logger.Info("No pending migrations found")
	} else {
		logger.Info("All migrations completed successfully",
			zap.Int("count", pendingCount))
	}

	return nil
}

// RollbackMigrations rollback specified number of migrations
func (mm *MigrationManager) RollbackMigrations(count int) error {
	if count <= 0 {
		return fmt.Errorf("rollback count must be greater than 0")
	}

	// Get applied migrations in reverse order
	var appliedRecords []MigrationRecord
	if err := mm.db.Order("applied_at DESC").Limit(count).Find(&appliedRecords).Error; err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(appliedRecords) == 0 {
		logger.Info("No migrations to rollback")
		return nil
	}

	if len(appliedRecords) < count {
		logger.Warn("Only found migrations to rollback",
			zap.Int("requested", count),
			zap.Int("available", len(appliedRecords)))
	}

	// Rollback each migration
	for _, record := range appliedRecords {
		migration, exists := mm.migrations[record.Version]
		if !exists {
			return fmt.Errorf("migration %s not found in registered migrations", record.Version)
		}

		logger.Info("Rolling back migration",
			zap.String("version", record.Version),
			zap.String("description", record.Description))

		if err := mm.rollbackSingleMigration(migration, record); err != nil {
			return fmt.Errorf("rollback failed for migration %s: %w", record.Version, err)
		}

		logger.Info("Migration rolled back successfully",
			zap.String("version", record.Version))
	}

	logger.Info("Rollback completed successfully",
		zap.Int("count", len(appliedRecords)))
	return nil
}

// GetMigrationStatus แสดงสถานะ migrations
func (mm *MigrationManager) GetMigrationStatus() error {
	// Create migrations table if not exists
	if err := mm.db.AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	var appliedRecords []MigrationRecord
	if err := mm.db.Order("applied_at ASC").Find(&appliedRecords).Error; err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedMap := make(map[string]MigrationRecord)
	for _, record := range appliedRecords {
		appliedMap[record.Version] = record
	}

	// Sort all migrations by version
	var versions []string
	for version := range mm.migrations {
		versions = append(versions, version)
	}
	sort.Strings(versions)

	// Show status
	appliedCount := 0
	pendingCount := 0

	logger.Info("Migration Status:")
	logger.Info("================")

	for _, version := range versions {
		migration := mm.migrations[version]
		if record, applied := appliedMap[version]; applied {
			appliedCount++
			logger.Info("✅ APPLIED",
				zap.String("version", version),
				zap.String("description", migration.Description()),
				zap.Time("applied_at", record.AppliedAt))
		} else {
			pendingCount++
			logger.Info("⏳ PENDING",
				zap.String("version", version),
				zap.String("description", migration.Description()))
		}
	}

	logger.Info("==================")
	logger.Info("Summary",
		zap.Int("applied", appliedCount),
		zap.Int("pending", pendingCount),
		zap.Int("total", len(versions)))

	return nil
}

// runSingleMigration รัน migration เดียวใน transaction
func (mm *MigrationManager) runSingleMigration(migration Migration) error {
	// Start transaction
	tx := mm.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	// Run migration
	if err := migration.Up(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("migration failed: %w", err)
	}

	// Record migration
	record := MigrationRecord{
		Version:     migration.Version(),
		Description: migration.Description(),
		AppliedAt:   time.Now().UTC(),
	}

	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	return nil
}

// rollbackSingleMigration rollback migration เดียว
func (mm *MigrationManager) rollbackSingleMigration(migration Migration, record MigrationRecord) error {
	// Start transaction
	tx := mm.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	// Run rollback
	if err := migration.Down(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("rollback failed: %w", err)
	}

	// Remove migration record
	if err := tx.Delete(&record).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit rollback: %w", err)
	}

	return nil
}
