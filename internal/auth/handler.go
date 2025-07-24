// internal/auth/handler.go (Updated)
package auth

import (
	"go-clean-gin/internal/entity"
	"go-clean-gin/pkg/errors"
	"go-clean-gin/pkg/logger"
	"go-clean-gin/pkg/response"
	"go-clean-gin/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthHandler struct {
	usecase AuthUsecase
}

func NewAuthHandler(usecase AuthUsecase) *AuthHandler {
	return &AuthHandler{
		usecase: usecase,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, username, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body entity.RegisterRequest true "Register user"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req entity.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		response.Error(c, 400, errors.ErrBadRequest, "Invalid request body", err.Error())
		return
	}

	if fieldErrors := validator.ValidateStruct(req); fieldErrors != nil {
		response.ValidationError(c, "Validation failed", fieldErrors)
		return
	}

	authResponse, err := h.usecase.Register(c.Request.Context(), &req)
	if err != nil {
		logger.Error("Failed to register user", zap.Error(err))

		// Handle specific errors
		if appErr, ok := err.(*errors.AppError); ok {
			response.Error(c, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details)
		} else {
			response.Error(c, 500, errors.ErrInternal, "Failed to register user", nil)
		}
		return
	}

	response.Success(c, 201, "User registered successfully", authResponse)
}

// Login godoc
// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body entity.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req entity.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		response.Error(c, 400, errors.ErrBadRequest, "Invalid request body", err.Error())
		return
	}

	if fieldErrors := validator.ValidateStruct(req); fieldErrors != nil {
		response.ValidationError(c, "Validation failed", fieldErrors)
		return
	}

	authResponse, err := h.usecase.Login(c.Request.Context(), &req)
	if err != nil {
		logger.Error("Failed to login", zap.Error(err))

		if appErr, ok := err.(*errors.AppError); ok {
			response.Error(c, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details)
		} else {
			response.Error(c, 401, errors.ErrInvalidCredentials, "Authentication failed", nil)
		}
		return
	}

	response.Success(c, 200, "Login successful", authResponse)
}

// Profile godoc
// @Summary Get user profile
// @Description Get current user profile
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/profile [get]
func (h *AuthHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 401, errors.ErrUnauthorized, "User not found in context", nil)
		return
	}

	userIDParsed, err := uuid.Parse(userID.(string))
	if err != nil {
		response.Error(c, 400, errors.ErrBadRequest, "Invalid user ID", err.Error())
		return
	}

	user, err := h.usecase.GetUserByID(c.Request.Context(), userIDParsed)
	if err != nil {
		logger.Error("Failed to get user profile", zap.Error(err))

		if appErr, ok := err.(*errors.AppError); ok {
			response.Error(c, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details)
		} else {
			response.Error(c, 500, errors.ErrInternal, "Failed to get user profile", nil)
		}
		return
	}

	response.Success(c, 200, "Profile retrieved successfully", user)
}
