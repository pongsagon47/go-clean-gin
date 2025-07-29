// internal/seeders/manager.go - Enhanced Seeder Manager with Dependencies
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
	Dependencies() []string // เพิ่ม method สำหรับ dependencies
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

// RunSeeders รัน seeders ทั้งหมด (จัดเรียงตาม dependencies)
func (sm *SeederManager) RunSeeders(seederName string) error {
	if len(sm.seeders) == 0 {
		logger.Info("No seeders found")
		return nil
	}

	logger.Info("Starting database seeding...",
		zap.Int("total_seeders", len(sm.seeders)))

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
		return nil
	}

	// เรียงลำดับ seeders ตาม dependencies
	orderedSeeders, err := sm.resolveDependencies()
	if err != nil {
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	successCount := 0
	for _, seeder := range orderedSeeders {
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
	return nil
}

// RunSpecificSeeder รัน seeder เฉพาะ พร้อม dependencies
func (sm *SeederManager) RunSpecificSeeder(seederName string) error {
	// หา seeder ที่ต้องการ
	var targetSeeder Seeder
	for _, seeder := range sm.seeders {
		if seeder.Name() == seederName {
			targetSeeder = seeder
			break
		}
	}

	if targetSeeder == nil {
		return fmt.Errorf("seeder %s not found", seederName)
	}

	// สร้าง dependency graph สำหรับ seeder นี้
	toRun, err := sm.resolveDependenciesFor(targetSeeder)
	if err != nil {
		return fmt.Errorf("failed to resolve dependencies for %s: %w", seederName, err)
	}

	// รัน seeders ตามลำดับ
	for _, seeder := range toRun {
		logger.Info("Running seeder", zap.String("name", seeder.Name()))

		if err := seeder.Run(sm.db); err != nil {
			logger.Error("Seeder failed",
				zap.String("name", seeder.Name()),
				zap.Error(err))
			return fmt.Errorf("seeder %s failed: %w", seeder.Name(), err)
		}

		logger.Info("Seeder completed successfully", zap.String("name", seeder.Name()))
	}

	return nil
}

// resolveDependencies เรียงลำดับ seeders ตาม dependencies
func (sm *SeederManager) resolveDependencies() ([]Seeder, error) {
	// สร้าง map สำหรับการค้นหา seeder
	seederMap := make(map[string]Seeder)
	for _, seeder := range sm.seeders {
		seederMap[seeder.Name()] = seeder
	}

	// ตรวจสอบว่าทุก dependency มีอยู่จริง
	for _, seeder := range sm.seeders {
		for _, dep := range seeder.Dependencies() {
			if _, exists := seederMap[dep]; !exists {
				return nil, fmt.Errorf("seeder %s depends on %s but %s not found",
					seeder.Name(), dep, dep)
			}
		}
	}

	// Topological sort
	return sm.topologicalSort(seederMap)
}

// resolveDependenciesFor แก้ไข dependencies สำหรับ seeder เฉพาะ
func (sm *SeederManager) resolveDependenciesFor(targetSeeder Seeder) ([]Seeder, error) {
	seederMap := make(map[string]Seeder)
	for _, seeder := range sm.seeders {
		seederMap[seeder.Name()] = seeder
	}

	visited := make(map[string]bool)
	visiting := make(map[string]bool)
	var result []Seeder

	var visit func(string) error
	visit = func(name string) error {
		if visiting[name] {
			return fmt.Errorf("circular dependency detected involving %s", name)
		}
		if visited[name] {
			return nil
		}

		seeder, exists := seederMap[name]
		if !exists {
			return fmt.Errorf("seeder %s not found", name)
		}

		visiting[name] = true

		// Visit dependencies first
		for _, dep := range seeder.Dependencies() {
			if err := visit(dep); err != nil {
				return err
			}
		}

		visiting[name] = false
		visited[name] = true
		result = append(result, seeder)
		return nil
	}

	if err := visit(targetSeeder.Name()); err != nil {
		return nil, err
	}

	return result, nil
}

// topologicalSort ใช้ Kahn's algorithm
func (sm *SeederManager) topologicalSort(seederMap map[string]Seeder) ([]Seeder, error) {
	// สร้าง adjacency list และ in-degree count
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize
	for name := range seederMap {
		graph[name] = []string{}
		inDegree[name] = 0
	}

	// Build graph และ count in-degrees
	for name, seeder := range seederMap {
		for _, dep := range seeder.Dependencies() {
			graph[dep] = append(graph[dep], name)
			inDegree[name]++
		}
	}

	// Queue สำหรับ nodes ที่ไม่มี dependencies
	var queue []string
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	var result []Seeder
	for len(queue) > 0 {
		// Dequeue
		current := queue[0]
		queue = queue[1:]

		result = append(result, seederMap[current])

		// ลด in-degree ของ neighbors
		for _, neighbor := range graph[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// ตรวจสอบ circular dependency
	if len(result) != len(seederMap) {
		return nil, fmt.Errorf("circular dependency detected in seeders")
	}

	return result, nil
}

// ListSeeders แสดงรายการ seeders ทั้งหมด พร้อม dependencies
func (sm *SeederManager) ListSeeders() {
	logger.Info("Registered Seeders:")
	logger.Info("==================")

	if len(sm.seeders) == 0 {
		logger.Info("No seeders registered")
		return
	}

	// เรียงลำดับตาม dependencies
	orderedSeeders, err := sm.resolveDependencies()
	if err != nil {
		logger.Error("Failed to resolve dependencies", zap.Error(err))
		// Fallback to original order
		orderedSeeders = sm.seeders
	}

	for i, seeder := range orderedSeeders {
		deps := seeder.Dependencies()
		if len(deps) > 0 {
			logger.Info(fmt.Sprintf("%d. %s (depends on: %s)",
				i+1, seeder.Name(), strings.Join(deps, ", ")))
		} else {
			logger.Info(fmt.Sprintf("%d. %s", i+1, seeder.Name()))
		}
	}

	logger.Info("==================")
	logger.Info("Total seeders", zap.Int("count", len(sm.seeders)))
}
