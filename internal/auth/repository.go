package auth

import (
	"context"
	"go-clean-gin/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) CreateUser(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ? AND is_active = ?", email, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ? AND is_active = ?", userID, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ? AND is_active = ?", username, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
