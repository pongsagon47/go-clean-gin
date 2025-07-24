package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-clean-gin/config"
	"go-clean-gin/internal/entity"
	"go-clean-gin/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authUsecase struct {
	repo   AuthRepository
	config *config.Config
}

func NewAuthUsecase(repo AuthRepository, config *config.Config) AuthUsecase {
	return &authUsecase{
		repo:   repo,
		config: config,
	}
}

func (u *authUsecase) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.AuthResponse, error) {
	// Check if user already exists
	existingUser, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Failed to check existing user by email", zap.Error(err))
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Check username
	existingUser, err = u.repo.GetUserByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Failed to check existing user by username", zap.Error(err))
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with username %s already exists", req.Username)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password", zap.Error(err))
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &entity.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	if err := u.repo.CreateUser(ctx, user); err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate token
	token, err := u.generateToken(user.ID)
	if err != nil {
		logger.Error("Failed to generate token", zap.Error(err))
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	logger.Info("User registered successfully", zap.String("user_id", user.ID.String()))

	return &entity.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (u *authUsecase) Login(ctx context.Context, req *entity.LoginRequest) (*entity.AuthResponse, error) {
	// Get user by email
	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invalid email or password")
		}
		logger.Error("Failed to get user by email", zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate token
	token, err := u.generateToken(user.ID)
	if err != nil {
		logger.Error("Failed to generate token", zap.Error(err))
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	logger.Info("User logged in successfully", zap.String("user_id", user.ID.String()))

	return &entity.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (u *authUsecase) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return u.repo.GetUserByID(ctx, userID)
}

func (u *authUsecase) ValidateToken(ctx context.Context, tokenString string) (*entity.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(u.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid token claims")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID in token: %w", err)
		}

		user, err := u.repo.GetUserByID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("user not found: %w", err)
		}

		return user, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (u *authUsecase) generateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Duration(u.config.JWT.ExpirationHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.config.JWT.Secret))
}
