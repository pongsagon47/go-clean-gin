package product

import (
	"context"
	"go-clean-gin/internal/entity"
	"go-clean-gin/pkg/errors"
	"go-clean-gin/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type productUsecase struct {
	repo ProductRepository
}

func NewProductUsecase(repo ProductRepository) ProductUsecase {
	return &productUsecase{
		repo: repo,
	}
}

func (u *productUsecase) CreateProduct(ctx context.Context, req *entity.CreateProductRequest, userID uuid.UUID) (*entity.Product, error) {
	product := &entity.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		IsActive:    true,
		CreatedBy:   userID,
	}

	if err := u.repo.CreateProduct(ctx, product); err != nil {
		logger.Error("Failed to create product", zap.Error(err))
		return nil, errors.Wrap(err, errors.ErrInternal, "Failed to create product", 500)
	}

	// Get the created product with user data
	createdProduct, err := u.repo.GetProductByID(ctx, product.ID)
	if err != nil {
		logger.Error("Failed to get created product", zap.Error(err))
		return nil, errors.Wrap(err, errors.ErrInternal, "Failed to get created product", 500)
	}

	logger.Info("Product created successfully", zap.String("product_id", product.ID.String()))
	return createdProduct, nil
}

func (u *productUsecase) GetProductByID(ctx context.Context, productID uuid.UUID) (*entity.Product, error) {
	product, err := u.repo.GetProductByID(ctx, productID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrProductNotFoundError
		}
		logger.Error("Failed to get product", zap.Error(err))
		return nil, errors.Wrap(err, errors.ErrInternal, "Failed to get product", 500)
	}

	return product, nil
}

func (u *productUsecase) GetProducts(ctx context.Context, filter *entity.ProductFilter) ([]*entity.Product, int64, error) {
	// Set default pagination if not provided
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	products, total, err := u.repo.GetProducts(ctx, filter)
	if err != nil {
		logger.Error("Failed to get products", zap.Error(err))
		return nil, 0, errors.Wrap(err, errors.ErrInternal, "Failed to get products", 500)
	}

	return products, total, nil
}

func (u *productUsecase) UpdateProduct(ctx context.Context, productID uuid.UUID, req *entity.UpdateProductRequest, userID uuid.UUID) (*entity.Product, error) {
	// Get existing product
	existingProduct, err := u.repo.GetProductByID(ctx, productID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrProductNotFoundError
		}
		logger.Error("Failed to get product for update", zap.Error(err))
		return nil, errors.Wrap(err, errors.ErrInternal, "Failed to get product", 500)
	}

	// Check if user is the owner of the product
	if existingProduct.CreatedBy != userID {
		return nil, errors.ErrInvalidOwnerError
	}

	// Update fields if provided
	if req.Name != nil {
		existingProduct.Name = *req.Name
	}
	if req.Description != nil {
		existingProduct.Description = *req.Description
	}
	if req.Price != nil {
		existingProduct.Price = *req.Price
	}
	if req.Stock != nil {
		existingProduct.Stock = *req.Stock
	}
	if req.Category != nil {
		existingProduct.Category = *req.Category
	}
	if req.IsActive != nil {
		existingProduct.IsActive = *req.IsActive
	}

	if err := u.repo.UpdateProduct(ctx, existingProduct); err != nil {
		logger.Error("Failed to update product", zap.Error(err))
		return nil, errors.Wrap(err, errors.ErrInternal, "Failed to update product", 500)
	}

	logger.Info("Product updated successfully", zap.String("product_id", productID.String()))
	return existingProduct, nil
}

func (u *productUsecase) DeleteProduct(ctx context.Context, productID uuid.UUID, userID uuid.UUID) error {
	// Get existing product
	existingProduct, err := u.repo.GetProductByID(ctx, productID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrProductNotFoundError
		}
		logger.Error("Failed to get product for deletion", zap.Error(err))
		return errors.Wrap(err, errors.ErrInternal, "Failed to get product", 500)
	}

	// Check if user is the owner of the product
	if existingProduct.CreatedBy != userID {
		return errors.ErrInvalidOwnerError
	}

	if err := u.repo.DeleteProduct(ctx, productID); err != nil {
		logger.Error("Failed to delete product", zap.Error(err))
		return errors.Wrap(err, errors.ErrInternal, "Failed to delete product", 500)
	}

	logger.Info("Product deleted successfully", zap.String("product_id", productID.String()))
	return nil
}
