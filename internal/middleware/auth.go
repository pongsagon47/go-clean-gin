package middleware

import (
	"net/http"
	"strings"

	"go-clean-gin/internal/auth"
	"go-clean-gin/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthMiddleware(authUsecase auth.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if token starts with "Bearer "
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]
		user, err := authUsecase.ValidateToken(c.Request.Context(), token)
		if err != nil {
			logger.Error("Token validation failed", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", user.ID.String())
		c.Set("user", user)
		c.Next()
	}
}
