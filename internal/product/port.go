package product

import (
	"context"
	"go-clean-gin/internal/entity"

	"github.com/google/uuid"
)

// ProductUsecase defines the business logic interface for products
type ProductUsecase interface {
	CreateProduct(ctx context.Context, req *entity.CreateProductRequest, userID uuid.UUID) (*entity.Product, error)
	GetProductByID(ctx context.Context, productID uuid.UUID) (*entity.Product, error)
	GetProducts(ctx context.Context, filter *entity.ProductFilter) ([]*entity.Product, int64, error)
	UpdateProduct(ctx context.Context, productID uuid.UUID, req *entity.UpdateProductRequest, userID uuid.UUID) (*entity.Product, error)
	DeleteProduct(ctx context.Context, productID uuid.UUID, userID uuid.UUID) error
}

// ProductRepository defines the data access interface for products
type ProductRepository interface {
	CreateProduct(ctx context.Context, product *entity.Product) error
	GetProductByID(ctx context.Context, productID uuid.UUID) (*entity.Product, error)
	GetProducts(ctx context.Context, filter *entity.ProductFilter) ([]*entity.Product, int64, error)
	UpdateProduct(ctx context.Context, product *entity.Product) error
	DeleteProduct(ctx context.Context, productID uuid.UUID) error
	GetProductsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Product, error)
}
