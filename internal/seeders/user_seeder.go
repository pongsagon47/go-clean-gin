package seeders

import (
	"go-clean-gin/pkg/logger"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserSeeder seeds the users table
type UserSeeder struct{}

// Run executes the seeder
func (s *UserSeeder) Run(db *gorm.DB) error {
	logger.Info("Running UserSeeder...")

	// Check if users already exist
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM users").Scan(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		logger.Info("Users already exist, skipping UserSeeder")
		return nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create sample users
	users := []map[string]interface{}{
		{
			"id":         uuid.New().String(),
			"email":      "admin@example.com",
			"username":   "admin",
			"password":   string(hashedPassword),
			"first_name": "Admin",
			"last_name":  "User",
			"is_active":  true,
			"created_at": time.Now().UTC(),
			"updated_at": time.Now().UTC(),
		},
		{
			"id":         uuid.New().String(),
			"email":      "john@example.com",
			"username":   "johndoe",
			"password":   string(hashedPassword),
			"first_name": "John",
			"last_name":  "Doe",
			"is_active":  true,
			"created_at": time.Now().UTC(),
			"updated_at": time.Now().UTC(),
		},
		{
			"id":         uuid.New().String(),
			"email":      "jane@example.com",
			"username":   "janedoe",
			"password":   string(hashedPassword),
			"first_name": "Jane",
			"last_name":  "Doe",
			"is_active":  true,
			"created_at": time.Now().UTC(),
			"updated_at": time.Now().UTC(),
		},
	}

	// Insert users
	for _, user := range users {
		if err := db.Exec(`
			INSERT INTO users (id, email, username, password, first_name, last_name, is_active, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, user["id"], user["email"], user["username"], user["password"],
			user["first_name"], user["last_name"], user["is_active"],
			user["created_at"], user["updated_at"]).Error; err != nil {
			return err
		}
	}

	logger.Info("UserSeeder completed successfully", zap.Int("users_created", len(users)))
	return nil
}

// Name returns seeder name
func (s *UserSeeder) Name() string {
	return "UserSeeder"
}

// Dependencies returns list of seeders that must run before this seeder
func (s *UserSeeder) Dependencies() []string {
	return []string{} // UserSeeder ไม่มี dependencies
}

// Auto-register seeder
func init() {
	Register(&UserSeeder{})
}
