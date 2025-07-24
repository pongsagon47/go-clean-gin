package product

import (
	"go-clean-gin/internal/entity"
	"go-clean-gin/pkg/errors"
	"go-clean-gin/pkg/logger"
	"go-clean-gin/pkg/response"
	"go-clean-gin/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProductHandler struct {
	usecase ProductUsecase
}

func NewProductHandler(usecase ProductUsecase) *ProductHandler {
	return &ProductHandler{
		usecase: usecase,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Security Bearer
// @Param product body entity.CreateProductRequest true "Create product"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req entity.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		response.Error(c, 400, errors.ErrBadRequest, "Invalid request body", err.Error())
		return
	}

	if fieldErrors := validator.ValidateStruct(req); fieldErrors != nil {
		response.ValidationError(c, "Validation failed", fieldErrors)
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 401, errors.ErrUnauthorized, "User not found in context", nil)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		response.Error(c, 400, errors.ErrBadRequest, "Invalid user ID", err.Error())
		return
	}

	product, err := h.usecase.CreateProduct(c.Request.Context(), &req, userID)
	if err != nil {
		logger.Error("Failed to create product", zap.Error(err))

		if appErr, ok := err.(*errors.AppError); ok {
			response.Error(c, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details)
		} else {
			response.Error(c, 500, errors.ErrInternal, "Failed to create product", nil)
		}
		return
	}

	response.Success(c, 201, "Product created successfully", product)
}

// GetProducts godoc
// @Summary Get products with filters
// @Description Get products with optional filters and pagination
// @Tags products
// @Accept json
// @Produce json
// @Param category query string false "Filter by category"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Param is_active query boolean false "Filter by active status"
// @Param search query string false "Search in name and description"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	var filter entity.ProductFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		logger.Error("Failed to bind query", zap.Error(err))
		response.Error(c, 400, errors.ErrBadRequest, "Invalid query parameters", err.Error())
		return
	}

	if fieldErrors := validator.ValidateStruct(filter); fieldErrors != nil {
		response.ValidationError(c, "Validation failed", fieldErrors)
		return
	}

	products, total, err := h.usecase.GetProducts(c.Request.Context(), &filter)
	if err != nil {
		logger.Error("Failed to get products", zap.Error(err))

		if appErr, ok := err.(*errors.AppError); ok {
			response.Error(c, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details)
		} else {
			response.Error(c, 500, errors.ErrInternal, "Failed to get products", nil)
		}
		return
	}

	meta := response.Pagination(filter.Page, filter.Limit, total)
	response.SuccessWithMeta(c, 200, "Products retrieved successfully", products, meta)
}

// GetProduct godoc
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		response.Error(c, 400, errors.ErrBadRequest, "Invalid product ID", err.Error())
		return
	}

	product, err := h.usecase.GetProductByID(c.Request.Context(), productID)
	if err != nil {
		logger.Error("Failed to get product", zap.Error(err))

		if appErr, ok := err.(*errors.AppError); ok {
			response.Error(c, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details)
		} else {
			response.Error(c, 500, errors.ErrInternal, "Failed to get product", nil)
		}
		return
	}

	response.Success(c, 200, "Product retrieved successfully", product)
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update product by ID
// @Tags products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Product ID"
// @Param product body entity.UpdateProductRequest true "Update product"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		response.Error(c, 400, errors.ErrBadRequest, "Invalid product ID", err.Error())
		return
	}

	var req entity.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		response.Error(c, 400, errors.ErrBadRequest, "Invalid request body", err.Error())
		return
	}

	if fieldErrors := validator.ValidateStruct(req); fieldErrors != nil {
		response.ValidationError(c, "Validation failed", fieldErrors)
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 401, errors.ErrUnauthorized, "User not found in context", nil)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		response.Error(c, 400, errors.ErrBadRequest, "Invalid user ID", err.Error())
		return
	}

	product, err := h.usecase.UpdateProduct(c.Request.Context(), productID, &req, userID)
	if err != nil {
		logger.Error("Failed to update product", zap.Error(err))

		if appErr, ok := err.(*errors.AppError); ok {
			response.Error(c, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details)
		} else {
			response.Error(c, 500, errors.ErrInternal, "Failed to update product", nil)
		}
		return
	}

	response.Success(c, 200, "Product updated successfully", product)
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete product by ID
// @Tags products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Product ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		response.Error(c, 400, errors.ErrBadRequest, "Invalid product ID", err.Error())
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 401, errors.ErrUnauthorized, "User not found in context", nil)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		response.Error(c, 400, errors.ErrBadRequest, "Invalid user ID", err.Error())
		return
	}

	err = h.usecase.DeleteProduct(c.Request.Context(), productID, userID)
	if err != nil {
		logger.Error("Failed to delete product", zap.Error(err))

		if appErr, ok := err.(*errors.AppError); ok {
			response.Error(c, appErr.StatusCode, appErr.Code, appErr.Message, appErr.Details)
		} else {
			response.Error(c, 500, errors.ErrInternal, "Failed to delete product", nil)
		}
		return
	}

	response.Success(c, 200, "Product deleted successfully", nil)
}
