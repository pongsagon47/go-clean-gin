package auth

import (
	"context"
	"go-clean-gin/internal/entity"

	"github.com/google/uuid"
)

// AuthUsecase defines the business logic interface for authentication
type AuthUsecase interface {
	Register(ctx context.Context, req *entity.RegisterRequest) (*entity.AuthResponse, error)
	Login(ctx context.Context, req *entity.LoginRequest) (*entity.AuthResponse, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	ValidateToken(ctx context.Context, token string) (*entity.User, error)
}

// AuthRepository defines the data access interface for authentication
type AuthRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
}
