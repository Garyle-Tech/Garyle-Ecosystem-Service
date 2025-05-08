package product

import (
	"strconv"

	productService "ecosystem.garyle/service/internal/app/service/wms/master-data/product"
	productModel "ecosystem.garyle/service/internal/domain/model/wms/master-data/product"
	"ecosystem.garyle/service/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	productService productService.ProductService
}

func NewProductHandler(productService productService.ProductService) *Handler {
	return &Handler{
		productService: productService,
	}
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var product productModel.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		if err.Error() == "EOF" {
			response.BadRequest(c, "Missing request body. Please provide a valid JSON payload.")
			return
		}

		response.BadRequest(c, err.Error())
		return
	}

	createdProduct, err := h.productService.Create(c.Request.Context(), &product)
	if err != nil {
		if isValidationCreateOrUpdateProductError(err) {
			response.BadRequest(c, err.Error())
			return
		}

		if err.Error() == "pq: duplicate key value violates unique constraint \"products_sku_key\"" {
			response.BadRequest(c, "Product with this SKU already exists, please use another SKU")
			return
		}

		response.Server(c, err.Error())
		return
	}

	response.Success(c, createdProduct, "Product created successfully")
}

func isValidationCreateOrUpdateProductError(err error) bool {
	validationErrors := []string{
		"sku is required",
		"name is required",
		"unit is required",
		"weight is required",
		"dimension is required",
	}

	for _, validationError := range validationErrors {
		if err.Error() == validationError {
			return true
		}
	}

	return false
}

// Get List Products
func (h *Handler) GetListProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	if limit <= 0 {
		limit = 10
	}

	if page <= 0 {
		page = 1
	}

	products, err := h.productService.List(c.Request.Context(), limit, page)
	if err != nil {
		response.Server(c, err.Error())
		return
	}

	total, err := h.productService.Count(c.Request.Context())
	if err != nil {
		response.Server(c, err.Error())
		return
	}

	response.SuccessWithPagination(c, products, "Products retrieved successfully", page, limit, total)
}

// Get Product By ID
func (h *Handler) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		response.BadRequest(c, "Invalid product ID")
		return
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		response.BadRequest(c, "Invalid product ID")
		return
	}

	product, err := h.productService.GetByID(c.Request.Context(), productID)
	if err != nil {
		response.Server(c, err.Error())
		return
	}

	if product == nil {
		response.NotFound(c, "Product not found")
		return
	}

	response.Success(c, product, "Product retrieved successfully")
}

// Update Product By ID
func (h *Handler) UpdateProductByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		response.BadRequest(c, "Invalid product ID")
		return
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		response.BadRequest(c, "Invalid product ID")
		return
	}

	var product productModel.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		if err.Error() == "EOF" {
			response.BadRequest(c, "Missing request body. Please provide a valid JSON payload.")
			return
		}

		response.BadRequest(c, err.Error())
		return
	}

	if err := h.productService.UpdateByID(c.Request.Context(), &product, productID); err != nil {
		if isValidationCreateOrUpdateProductError(err) {
			response.BadRequest(c, err.Error())
			return
		}

		if err.Error() == "product not found" {
			response.NotFound(c, err.Error())
			return
		}

		if err.Error() == "pq: duplicate key value violates unique constraint \"products_sku_key\"" {
			response.BadRequest(c, "Product with this SKU already exists, please use another SKU")
			return
		}

		response.Server(c, err.Error())
		return
	}

	response.Success(c, product, "Product updated successfully")
}

// Delete Product By ID
func (h *Handler) DeleteProductByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		response.BadRequest(c, "Invalid product ID")
		return
	}

	productID, err := strconv.Atoi(id)
	if err != nil {
		response.BadRequest(c, "Invalid product ID")
		return
	}

	if err := h.productService.DeleteByID(c.Request.Context(), productID); err != nil {
		if err.Error() == "product not found" {
			response.NotFound(c, err.Error())
			return
		}

		response.Server(c, err.Error())
		return
	}

	response.Success(c, nil, "Product deleted successfully")
}

func (h *Handler) RegisterProductRoutes(router *gin.RouterGroup) {
	productRoutes := router.Group("/product")
	{
		productRoutes.POST("", h.CreateProduct)
		productRoutes.GET("", h.GetListProducts)
		productRoutes.GET("/:id", h.GetProductByID)
		productRoutes.PUT("/:id", h.UpdateProductByID)
		productRoutes.DELETE("/:id", h.DeleteProductByID)
	}
}
