package seeders

import (
	"go-clean-gin/pkg/logger"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProductSeeder seeds the products table
type ProductSeeder struct{}

// Run executes the seeder
func (s *ProductSeeder) Run(db *gorm.DB) error {
	logger.Info("Running ProductSeeder...")

	// Check if data already exists
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM tb_products").Scan(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		logger.Info("products already exist, skipping ProductSeeder")
		return nil
	}

	// Get admin user ID for created_by field
	var adminUserID string
	if err := db.Raw("SELECT id FROM tb_users WHERE email = ? LIMIT 1", "admin@example.com").Scan(&adminUserID).Error; err != nil {
		logger.Error("Admin user not found for ProductSeeder", zap.Error(err))
		return err
	}

	// Create sample products
	products := []map[string]interface{}{
		{
			"id":          uuid.New().String(),
			"name":        "MacBook Pro 16",
			"description": "Apple MacBook Pro 16-inch with M2 Pro chip",
			"price":       2499.99,
			"stock":       10,
			"category":    "Electronics",
			"is_active":   true,
			"created_by":  adminUserID,
			"created_at":  time.Now().UTC(),
			"updated_at":  time.Now().UTC(),
		},
		{
			"id":          uuid.New().String(),
			"name":        "iPhone 15 Pro",
			"description": "Latest iPhone with titanium design",
			"price":       999.99,
			"stock":       25,
			"category":    "Electronics",
			"is_active":   true,
			"created_by":  adminUserID,
			"created_at":  time.Now().UTC(),
			"updated_at":  time.Now().UTC(),
		},
		{
			"id":          uuid.New().String(),
			"name":        "Nike Air Force 1",
			"description": "Classic white sneakers",
			"price":       90.00,
			"stock":       50,
			"category":    "Fashion",
			"is_active":   true,
			"created_by":  adminUserID,
			"created_at":  time.Now().UTC(),
			"updated_at":  time.Now().UTC(),
		},
		{
			"id":          uuid.New().String(),
			"name":        "The Go Programming Language",
			"description": "Comprehensive guide to Go programming",
			"price":       45.99,
			"stock":       100,
			"category":    "Books",
			"is_active":   true,
			"created_by":  adminUserID,
			"created_at":  time.Now().UTC(),
			"updated_at":  time.Now().UTC(),
		},
		{
			"id":          uuid.New().String(),
			"name":        "Wireless Mouse",
			"description": "Ergonomic wireless mouse with long battery life",
			"price":       29.99,
			"stock":       75,
			"category":    "Electronics",
			"is_active":   true,
			"created_by":  adminUserID,
			"created_at":  time.Now().UTC(),
			"updated_at":  time.Now().UTC(),
		},
	}

	// Insert products
	for _, product := range products {
		if err := db.Exec(`
			INSERT INTO tb_products (id, name, description, price, stock, category, is_active, created_by, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, product["id"], product["name"], product["description"], product["price"],
			product["stock"], product["category"], product["is_active"],
			product["created_by"], product["created_at"], product["updated_at"]).Error; err != nil {
			return err
		}
	}

	logger.Info("ProductSeeder completed successfully")
	return nil
}

// Name returns seeder name
func (s *ProductSeeder) Name() string {
	return "ProductSeeder"
}

// Dependencies returns list of seeders that must run before this seeder
func (s *ProductSeeder) Dependencies() []string {
	return []string{
		"UserSeeder",
	}
}

// Auto-register seeder
func init() {
	Register(&ProductSeeder{})
}
