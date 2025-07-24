package auth

import (
	"context"
	"testing"

	"go-clean-gin/config"
	"go-clean-gin/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock repository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockAuthRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockAuthRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockAuthRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockAuthRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestAuthUsecase_Register_Success(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:          "test-secret",
			ExpirationHours: 24,
		},
	}
	usecase := NewAuthUsecase(mockRepo, cfg)

	req := &entity.RegisterRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	// Mock expectations
	mockRepo.On("GetUserByEmail", mock.Anything, req.Email).Return((*entity.User)(nil), gorm.ErrRecordNotFound)
	mockRepo.On("GetUserByUsername", mock.Anything, req.Username).Return((*entity.User)(nil), gorm.ErrRecordNotFound)
	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

	// Test
	result, err := usecase.Register(context.Background(), req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Email, result.User.Email)
	assert.NotEmpty(t, result.Token)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_Register_EmailExists(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:          "test-secret",
			ExpirationHours: 24,
		},
	}
	usecase := NewAuthUsecase(mockRepo, cfg)

	req := &entity.RegisterRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	existingUser := &entity.User{
		ID:    uuid.New(),
		Email: req.Email,
	}

	// Mock expectations
	mockRepo.On("GetUserByEmail", mock.Anything, req.Email).Return(existingUser, nil)

	// Test
	result, err := usecase.Register(context.Background(), req)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}
