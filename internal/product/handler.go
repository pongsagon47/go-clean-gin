package product

import (
	"net/http"

	"go-clean-gin/internal/entity"
	"go-clean-gin/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProductHandler struct {
	usecase   ProductUsecase
	validator *validator.Validate
}

func NewProductHandler(usecase ProductUsecase) *ProductHandler {
	return &ProductHandler{
		usecase:   usecase,
		validator: validator.New(),
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
// @Success 201 {object} entity.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req entity.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"message": err.Error(),
		})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not found in context",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": err.Error(),
		})
		return
	}

	product, err := h.usecase.CreateProduct(c.Request.Context(), &req, userID)
	if err != nil {
		logger.Error("Failed to create product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create product",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"data":    product,
	})
}

// GetProduct godoc
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} entity.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid product ID",
			"message": err.Error(),
		})
		return
	}

	product, err := h.usecase.GetProductByID(c.Request.Context(), productID)
	if err != nil {
		logger.Error("Failed to get product", zap.Error(err))
		status := http.StatusInternalServerError
		if err.Error() == "product not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"error":   "Failed to get product",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product retrieved successfully",
		"data":    product,
	})
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
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	var filter entity.ProductFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		logger.Error("Failed to bind query", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"message": err.Error(),
		})
		return
	}

	if err := h.validator.Struct(filter); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"message": err.Error(),
		})
		return
	}

	products, total, err := h.usecase.GetProducts(c.Request.Context(), &filter)
	if err != nil {
		logger.Error("Failed to get products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get products",
			"message": err.Error(),
		})
		return
	}

	// Calculate pagination info
	totalPages := (int(total) + filter.Limit - 1) / filter.Limit
	if filter.Limit == 0 {
		totalPages = 1
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Products retrieved successfully",
		"data":    products,
		"meta": gin.H{
			"total":       total,
			"page":        filter.Page,
			"limit":       filter.Limit,
			"total_pages": totalPages,
			"has_next":    filter.Page < totalPages,
			"has_prev":    filter.Page > 1,
		},
	})
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
// @Success 200 {object} entity.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid product ID",
			"message": err.Error(),
		})
		return
	}

	var req entity.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"message": err.Error(),
		})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not found in context",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": err.Error(),
		})
		return
	}

	product, err := h.usecase.UpdateProduct(c.Request.Context(), productID, &req, userID)
	if err != nil {
		logger.Error("Failed to update product", zap.Error(err))
		status := http.StatusInternalServerError
		message := err.Error()

		switch message {
		case "product not found":
			status = http.StatusNotFound
		case "unauthorized: you can only update your own products":
			status = http.StatusForbidden
		}

		c.JSON(status, gin.H{
			"error":   "Failed to update product",
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"data":    product,
	})
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete product by ID
// @Tags products
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid product ID",
			"message": err.Error(),
		})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User not found in context",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": err.Error(),
		})
		return
	}

	err = h.usecase.DeleteProduct(c.Request.Context(), productID, userID)
	if err != nil {
		logger.Error("Failed to delete product", zap.Error(err))
		status := http.StatusInternalServerError
		message := err.Error()

		switch message {
		case "product not found":
			status = http.StatusNotFound
		case "unauthorized: you can only delete your own products":
			status = http.StatusForbidden
		}

		c.JSON(status, gin.H{
			"error":   "Failed to delete product",
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
