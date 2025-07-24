package router

import (
	"go-clean-gin/internal/container"
	"go-clean-gin/internal/middleware"

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

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Server is running",
			"version": "1.0.0",
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
