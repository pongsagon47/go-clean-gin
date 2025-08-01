package migrations

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null" validate:"required,min=1,max=255"`
	Description string         `json:"description" gorm:"type:text"`
	Price       float64        `json:"price" gorm:"not null" validate:"required,min=0"`
	Stock       int            `json:"stock" gorm:"not null;default:0" validate:"min=0"`
	Category    string         `json:"category" gorm:"not null" validate:"required"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedBy   uuid.UUID      `json:"created_by" gorm:"type:uuid;not null"`
	User        User           `json:"user,omitempty" gorm:"foreignKey:CreatedBy"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Product) TableName() string {
	return "tb_products"
}

// CreateProductsTable migration - Create products table
type CreateProductsTable struct{}

// Up creates the products table
func (m *CreateProductsTable) Up(db *gorm.DB) error {
	return db.AutoMigrate(&Product{})
}

// Down drops the products table
func (m *CreateProductsTable) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(&Product{})
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
