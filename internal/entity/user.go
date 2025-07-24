package entity

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

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" validate:"required,min=1,max=100"`
}

type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
