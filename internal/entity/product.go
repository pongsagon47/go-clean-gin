package entity

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

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Stock       int     `json:"stock" validate:"min=0"`
	Category    string  `json:"category" validate:"required"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	Stock       *int     `json:"stock,omitempty" validate:"omitempty,min=0"`
	Category    *string  `json:"category,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

type ProductFilter struct {
	Category string  `form:"category"`
	MinPrice float64 `form:"min_price"`
	MaxPrice float64 `form:"max_price"`
	IsActive *bool   `form:"is_active"`
	Search   string  `form:"search"`
	Page     int     `form:"page" validate:"min=1"`
	Limit    int     `form:"limit" validate:"min=1,max=100"`
}
