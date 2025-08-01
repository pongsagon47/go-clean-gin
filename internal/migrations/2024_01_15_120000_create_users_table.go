package migrations

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=50"`
	Password  string         `json:"-" gorm:"not null" validate:"required,min=6"`
	FirstName string         `json:"first_name" gorm:"not null" validate:"required,min=1,max=100"`
	LastName  string         `json:"last_name" gorm:"not null" validate:"required,min=1,max=100"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (User) TableName() string {
	return "tb_users"
}

// CreateUsersTable migration - Create users table
type CreateUsersTable struct{}

// Up creates the users table
func (m *CreateUsersTable) Up(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}

// Down drops the users table
func (m *CreateUsersTable) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(&User{})
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
