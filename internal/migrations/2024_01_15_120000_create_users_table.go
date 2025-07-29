package migrations

import (
	"gorm.io/gorm"
)

// CreateUsersTable migration - Create users table
type CreateUsersTable struct{}

// Up creates the users table
func (m *CreateUsersTable) Up(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE
		)
	`).Error
}

// Down drops the users table
func (m *CreateUsersTable) Down(db *gorm.DB) error {
	return db.Exec("DROP TABLE IF EXISTS users").Error
}

// Description returns migration description
func (m *CreateUsersTable) Description() string {
	return "Create users table"
}

// Version returns migration version
func (m *CreateUsersTable) Version() string {
	return "2024_01_15_120000_create_users_table"
}

// Auto-register migration
func init() {
	Register(&CreateUsersTable{})
}
