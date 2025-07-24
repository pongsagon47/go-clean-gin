package auth

import (
	"context"
	"fmt"
	"time"

	"go-clean-gin/config"
	"go-clean-gin/internal/entity"
	"go-clean-gin/pkg/errors"
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
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error("Failed to check existing user by email", zap.Error(err))
		return nil, errors.Wrap(err, errors.ErrInternal, "Failed to check existing user", 500)
	}
	if existingUser != nil {
		return nil, errors.New(errors.ErrUserExists,
			fmt.Sprintf("User with email %s already exists", req.Email), 409)
	}

	// Check username
	existingUser, err = u.repo.GetUserByUsername(ctx, req.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error("Failed to check existing user by username", zap.Error(err))
		return nil, errors.Wrap(err, errors.ErrInternal, "Failed to check existing user", 500)
	}
	if existingUser != nil {
		return nil, errors.New(errors.ErrUserExists,
			fmt.Sprintf("User with username %s already exists", req.Username), 409)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password", zap.Error(err))
		return nil, errors.Wrap(err, errors.ErrInternal, "Failed to hash password", 500)
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
		return nil, errors.Wrap(err, errors.ErrInternal, "Failed to create user", 500)
	}

	// Generate token
	token, err := u.generateToken(user.ID)
	if err != nil {
		logger.Error("Failed to generate token", zap.Error(err))
		return nil, errors.Wrap(err, errors.ErrInternal, "Failed to generate token", 500)
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
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrInvalidCredentialsError
		}
		logger.Error("Failed to get user by email", zap.Error(err))
		return nil, errors.Wrap(err, errors.ErrInternal, "Failed to get user", 500)
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.ErrInvalidCredentialsError
	}

	// Generate token
	token, err := u.generateToken(user.ID)
	if err != nil {
		logger.Error("Failed to generate token", zap.Error(err))
		return nil, errors.Wrap(err, errors.ErrInternal, "Failed to generate token", 500)
	}

	logger.Info("User logged in successfully", zap.String("user_id", user.ID.String()))

	return &entity.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (u *authUsecase) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFoundError
		}
		logger.Error("Failed to get user by ID", zap.Error(err))
		return nil, errors.Wrap(err, errors.ErrInternal, "Failed to get user", 500)
	}
	return user, nil
}

func (u *authUsecase) ValidateToken(ctx context.Context, tokenString string) (*entity.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(u.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, errors.ErrTokenInvalidError.WithDetails(err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.ErrTokenInvalidError.WithDetails("Invalid token claims")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, errors.ErrTokenInvalidError.WithDetails("Invalid user ID in token")
		}

		user, err := u.repo.GetUserByID(ctx, userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.ErrUserNotFoundError
			}
			return nil, errors.Wrap(err, errors.ErrInternal, "Failed to get user", 500)
		}

		return user, nil
	}

	return nil, errors.ErrTokenInvalidError
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
