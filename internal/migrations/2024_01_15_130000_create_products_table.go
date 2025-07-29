package migrations

import (
	"gorm.io/gorm"
)

// CreateProductsTable migration - Create products table
type CreateProductsTable struct{}

// Up creates the products table
func (m *CreateProductsTable) Up(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE products (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price DECIMAL(10,2) NOT NULL,
			stock INTEGER NOT NULL DEFAULT 0,
			category VARCHAR(100) NOT NULL,
			is_active BOOLEAN DEFAULT true,
			created_by UUID NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE,
			FOREIGN KEY (created_by) REFERENCES users(id)
		)
	`).Error
}

// Down drops the products table
func (m *CreateProductsTable) Down(db *gorm.DB) error {
	return db.Exec("DROP TABLE IF EXISTS products").Error
}

// Description returns migration description
func (m *CreateProductsTable) Description() string {
	return "Create products table"
}

// Version returns migration version
func (m *CreateProductsTable) Version() string {
	return "2024_01_15_130000_create_products_table"
}

// Auto-register migration
func init() {
	Register(&CreateProductsTable{})
}
