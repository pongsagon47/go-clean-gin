// internal/seeders/manager.go - Seeder Manager สำหรับระบบ Laravel-style
package seeders

import (
	"fmt"
	"strings"

	"go-clean-gin/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Seeder interface ที่แต่ละไฟล์ต้อง implement
type Seeder interface {
	Run(db *gorm.DB) error
	Name() string
}

// SeederManager จัดการ seeders
type SeederManager struct {
	db      *gorm.DB
	seeders []Seeder
}

// Global seeder manager instance
var globalSeederManager *SeederManager
var registeredSeeders []Seeder

// NewSeederManager สร้าง seeder manager ใหม่
func NewSeederManager(db *gorm.DB) *SeederManager {
	manager := &SeederManager{
		db:      db,
		seeders: make([]Seeder, 0),
	}

	// Register all seeders that were registered during init()
	for _, seeder := range registeredSeeders {
		manager.RegisterSeeder(seeder)
	}

	return manager
}

// SetGlobalSeederManager ตั้งค่า global seeder manager
func SetGlobalSeederManager(manager *SeederManager) {
	globalSeederManager = manager
}

// Register ฟังก์ชันสำหรับให้แต่ละไฟล์เรียกใช้ใน init()
func Register(seeder Seeder) {
	registeredSeeders = append(registeredSeeders, seeder)

	// ถ้ามี global manager แล้ว ให้ register เลย
	if globalSeederManager != nil {
		globalSeederManager.RegisterSeeder(seeder)
	}
}

// RegisterSeeder ลงทะเบียน seeder
func (sm *SeederManager) RegisterSeeder(seeder Seeder) {
	sm.seeders = append(sm.seeders, seeder)
}

// RunSeeders รัน seeders ทั้งหมด
func (sm *SeederManager) RunSeeders(seederName string) error {
	if len(sm.seeders) == 0 {
		logger.Info("No seeders found")
		return nil
	}

	logger.Info("Starting database seeding...",
		zap.Int("total_seeders", len(sm.seeders)))

	successCount := 0
	if seederName != "" {
		if !strings.HasSuffix(seederName, "Seeder") {
			seederName += "Seeder"
		}

		if err := sm.RunSpecificSeeder(seederName); err != nil {
			logger.Error("Seeder failed",
				zap.String("name", seederName),
				zap.Error(err))
			return fmt.Errorf("seeder %s failed: %w", seederName, err)
		}

		logger.Info("Seeder completed successfully", zap.String("name", seederName))
	} else {
		for _, seeder := range sm.seeders {
			logger.Info("Running seeder", zap.String("name", seeder.Name()))

			if err := seeder.Run(sm.db); err != nil {
				logger.Error("Seeder failed",
					zap.String("name", seeder.Name()),
					zap.Error(err))
				return fmt.Errorf("seeder %s failed: %w", seeder.Name(), err)
			}

			successCount++
			logger.Info("Seeder completed successfully", zap.String("name", seeder.Name()))
		}
		logger.Info("All seeders completed successfully", zap.Int("count", successCount))
	}

	return nil
}

// RunSpecificSeeder รัน seeder เฉพาะ
func (sm *SeederManager) RunSpecificSeeder(seederName string) error {
	for _, seeder := range sm.seeders {
		if seeder.Name() == seederName {
			logger.Info("Running specific seeder",
				zap.String("name", seederName))

			if err := seeder.Run(sm.db); err != nil {
				logger.Error("Seeder failed",
					zap.String("name", seederName),
					zap.Error(err))
				return fmt.Errorf("seeder %s failed: %w", seederName, err)
			}

			logger.Info("Seeder completed successfully",
				zap.String("name", seederName))
			return nil
		}
	}

	return fmt.Errorf("seeder %s not found", seederName)
}

// ListSeeders แสดงรายการ seeders ทั้งหมด
func (sm *SeederManager) ListSeeders() {
	logger.Info("Registered Seeders:")
	logger.Info("==================")

	if len(sm.seeders) == 0 {
		logger.Info("No seeders registered")
		return
	}

	for i, seeder := range sm.seeders {
		logger.Info(fmt.Sprintf("%d. %s", i+1, seeder.Name()))
	}

	logger.Info("==================")
	logger.Info("Total seeders", zap.Int("count", len(sm.seeders)))
}
