package auth

import (
	"net/http"

	"go-clean-gin/internal/entity"
	"go-clean-gin/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthHandler struct {
	usecase   AuthUsecase
	validator *validator.Validate
}

func NewAuthHandler(usecase AuthUsecase) *AuthHandler {
	return &AuthHandler{
		usecase:   usecase,
		validator: validator.New(),
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, username, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body entity.RegisterRequest true "Register user"
// @Success 201 {object} entity.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req entity.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"message": err.Error(),
		})
		return
	}

	response, err := h.usecase.Register(c.Request.Context(), &req)
	if err != nil {
		logger.Error("Failed to register user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to register user",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    response,
	})
}

// Login godoc
// @Summary Login user
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body entity.LoginRequest true "Login credentials"
// @Success 200 {object} entity.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req entity.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"message": err.Error(),
		})
		return
	}

	response, err := h.usecase.Login(c.Request.Context(), &req)
	if err != nil {
		logger.Error("Failed to login", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    response,
	})
}

// Profile godoc
// @Summary Get user profile
// @Description Get current user profile
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} entity.User
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/profile [get]
func (h *AuthHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not found in context",
		})
		return
	}

	userIDParsed, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": err.Error(),
		})
		return
	}

	user, err := h.usecase.GetUserByID(c.Request.Context(), userIDParsed)
	if err != nil {
		logger.Error("Failed to get user profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user profile",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data":    user,
	})
}
