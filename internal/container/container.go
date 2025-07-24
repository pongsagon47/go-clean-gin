package container

import (
	"go-clean-gin/config"
	"go-clean-gin/internal/auth"
	"go-clean-gin/internal/product"

	"gorm.io/gorm"
)

type Container struct {
	Config *config.Config
	DB     *gorm.DB

	// Repositories
	AuthRepo    auth.AuthRepository
	ProductRepo product.ProductRepository

	// Usecases
	AuthUsecase    auth.AuthUsecase
	ProductUsecase product.ProductUsecase

	// Handlers
	AuthHandler    *auth.AuthHandler
	ProductHandler *product.ProductHandler
}

func NewContainer(cfg *config.Config, db *gorm.DB) *Container {
	container := &Container{
		Config: cfg,
		DB:     db,
	}

	container.setupRepositories()
	container.setupUsecases()
	container.setupHandlers()

	return container
}

func (c *Container) setupRepositories() {
	c.AuthRepo = auth.NewAuthRepository(c.DB)
	c.ProductRepo = product.NewProductRepository(c.DB)
}

func (c *Container) setupUsecases() {
	c.AuthUsecase = auth.NewAuthUsecase(c.AuthRepo, c.Config)
	c.ProductUsecase = product.NewProductUsecase(c.ProductRepo)
}

func (c *Container) setupHandlers() {
	c.AuthHandler = auth.NewAuthHandler(c.AuthUsecase)
	c.ProductHandler = product.NewProductHandler(c.ProductUsecase)
}
