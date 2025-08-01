package container

import (
	"go-clean-gin/config"
	"go-clean-gin/internal/auth"
	"go-clean-gin/internal/product"
	"go-clean-gin/pkg/logger"
	"go-clean-gin/pkg/mail"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	Config *config.Config
	DB     *gorm.DB
	Mail   *mail.Mailer

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

	mail, err := mail.NewGomail(&cfg.Email)
	if err != nil {
		logger.Fatal("Failed to initialize email", zap.Error(err))
	}

	if err := mail.TestConnection(); err != nil {
		logger.Fatal("Failed to test email connection", zap.Error(err))
	}

	logger.Info("Email connection successful")

	// Auth
	authRepo := auth.NewAuthRepository(db)
	authUsecase := auth.NewAuthUsecase(authRepo, cfg, mail)
	authHandler := auth.NewAuthHandler(authUsecase)

	// Product
	productRepo := product.NewProductRepository(db)
	productUsecase := product.NewProductUsecase(productRepo)
	productHandler := product.NewProductHandler(productUsecase)

	return &Container{
		Config: cfg,
		DB:     db,
		Mail:   mail,

		// Repositories
		AuthRepo:    authRepo,
		ProductRepo: productRepo,

		// Usecases
		AuthUsecase:    authUsecase,
		ProductUsecase: productUsecase,

		// Handlers
		AuthHandler:    authHandler,
		ProductHandler: productHandler,
	}
}
