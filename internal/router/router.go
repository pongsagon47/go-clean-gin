package router

import (
	"go-clean-gin/internal/container"
	"go-clean-gin/internal/middleware"
	"go-clean-gin/pkg/response"

	"github.com/gin-gonic/gin"
)

func SetupRouter(container *container.Container) *gin.Engine {
	// Set Gin mode based on environment
	if container.Config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Recovery())
	router.Use(middleware.Logging())
	router.Use(middleware.ErrorHandler()) // Add error handler middleware

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		response.Success(c, 200, "Server is running", gin.H{
			"status":  "OK",
			"version": "1.0.0",
			"env":     container.Config.Env,
		})
	})

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		response.Error(c, 404, "NOT_FOUND", "Route not found", gin.H{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", container.AuthHandler.Register)
			authRoutes.POST("/login", container.AuthHandler.Login)

			// Protected auth routes
			authProtected := authRoutes.Group("/")
			authProtected.Use(middleware.AuthMiddleware(container.AuthUsecase))
			{
				authProtected.GET("/profile", container.AuthHandler.Profile)
			}
		}

		// Product routes
		productRoutes := v1.Group("/products")
		{
			// Public product routes
			productRoutes.GET("", container.ProductHandler.GetProducts)
			productRoutes.GET("/:id", container.ProductHandler.GetProduct)

			// Protected product routes
			productProtected := productRoutes.Group("/")
			productProtected.Use(middleware.AuthMiddleware(container.AuthUsecase))
			{
				productProtected.POST("", container.ProductHandler.CreateProduct)
				productProtected.PUT("/:id", container.ProductHandler.UpdateProduct)
				productProtected.DELETE("/:id", container.ProductHandler.DeleteProduct)
			}
		}
	}

	return router
}
