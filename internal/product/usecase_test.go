package product

import (
	"context"
	"testing"

	"go-clean-gin/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mock repository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) CreateProduct(ctx context.Context, product *entity.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) GetProductByID(ctx context.Context, productID uuid.UUID) (*entity.Product, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockProductRepository) GetProducts(ctx context.Context, filter *entity.ProductFilter) ([]*entity.Product, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entity.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) UpdateProduct(ctx context.Context, product *entity.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	args := m.Called(ctx, productID)
	return args.Error(0)
}

func (m *MockProductRepository) GetProductsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Product, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entity.Product), args.Error(1)
}

func TestProductUsecase_CreateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)

	userID := uuid.New()
	req := &entity.CreateProductRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		Stock:       10,
		Category:    "electronics",
	}

	createdProduct := &entity.Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		CreatedBy:   userID,
		IsActive:    true,
	}

	// Mock expectations
	mockRepo.On("CreateProduct", mock.Anything, mock.AnythingOfType("*entity.Product")).Return(nil)
	mockRepo.On("GetProductByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(createdProduct, nil)

	// Test
	result, err := usecase.CreateProduct(context.Background(), req, userID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, userID, result.CreatedBy)
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_GetProductByID_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)

	productID := uuid.New()
	product := &entity.Product{
		ID:       productID,
		Name:     "Test Product",
		Price:    99.99,
		IsActive: true,
	}

	// Mock expectations
	mockRepo.On("GetProductByID", mock.Anything, productID).Return(product, nil)

	// Test
	result, err := usecase.GetProductByID(context.Background(), productID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, productID, result.ID)
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_GetProductByID_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)

	productID := uuid.New()

	// Mock expectations
	mockRepo.On("GetProductByID", mock.Anything, productID).Return((*entity.Product)(nil), gorm.ErrRecordNotFound)

	// Test
	result, err := usecase.GetProductByID(context.Background(), productID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)
}

func TestProductUsecase_UpdateProduct_Unauthorized(t *testing.T) {
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)

	productID := uuid.New()
	userID := uuid.New()
	ownerID := uuid.New() // Different from userID

	existingProduct := &entity.Product{
		ID:        productID,
		Name:      "Existing Product",
		CreatedBy: ownerID, // Different owner
	}

	req := &entity.UpdateProductRequest{
		Name: stringPtr("Updated Product"),
	}

	// Mock expectations
	mockRepo.On("GetProductByID", mock.Anything, productID).Return(existingProduct, nil)

	// Test
	result, err := usecase.UpdateProduct(context.Background(), productID, req, userID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unauthorized")
	mockRepo.AssertExpectations(t)
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
