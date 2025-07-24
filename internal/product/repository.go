package product

import (
	"context"
	"fmt"
	"go-clean-gin/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) CreateProduct(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetProductByID(ctx context.Context, productID uuid.UUID) (*entity.Product, error) {
	var product entity.Product
	err := r.db.WithContext(ctx).Preload("User").Where("id = ?", productID).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetProducts(ctx context.Context, filter *entity.ProductFilter) ([]*entity.Product, int64, error) {
	var products []*entity.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Product{}).Preload("User")

	// Apply filters
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}

	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}

	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	if filter.Search != "" {
		searchTerm := fmt.Sprintf("%%%s%%", filter.Search)
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	// Order by created_at desc
	query = query.Order("created_at DESC")

	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, product *entity.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Product{}, productID).Error
}

func (r *productRepository) GetProductsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Product, error) {
	var products []*entity.Product
	err := r.db.WithContext(ctx).Preload("User").Where("created_by = ?", userID).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}
