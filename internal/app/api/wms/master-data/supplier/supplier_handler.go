package supplier

import (
	"strconv"

	supplierService "ecosystem.garyle/service/internal/app/service/wms/master-data/supplier"
	supplierModel "ecosystem.garyle/service/internal/domain/model/wms/master-data/supplier"
	"ecosystem.garyle/service/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type SupplierHandler struct {
	supplierService supplierService.SupplierService
}

func NewSupplierHandler(supplierService supplierService.SupplierService) *SupplierHandler {
	return &SupplierHandler{supplierService: supplierService}
}

// create supplier
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var supplier supplierModel.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {

		// validation error when json body empty
		if err.Error() == "EOF" {
			response.BadRequest(c, "Missing request body. Please provide a valid JSON payload.")
			return
		}

		response.BadRequest(c, err.Error())
		return
	}

	// create supplier
	createdSupplier, err := h.supplierService.Create(c.Request.Context(), &supplier)
	if err != nil {
		// validation error
		if validateCreateOrUpdateSupplierError(err) {
			response.BadRequest(c, err.Error())
			return
		}

		response.Server(c, err.Error())
		return
	}

	response.Success(c, createdSupplier, "Supplier created successfully")
}

// get list supplier
func (h *SupplierHandler) GetListSupplier(c *gin.Context) {
	// limit & page
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	if limit <= 0 {
		limit = 10
	}

	if page <= 0 {
		page = 1
	}

	suppliers, err := h.supplierService.List(c.Request.Context(), limit, page)
	if err != nil {
		response.Server(c, err.Error())
		return
	}

	total, err := h.supplierService.Count(c.Request.Context())
	if err != nil {
		response.Server(c, err.Error())
		return
	}

	response.SuccessWithPagination(c, suppliers, "Supplier list retrieved successfully", page, limit, total)
}

// get supplier by id
func (h *SupplierHandler) GetSupplierByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		response.BadRequest(c, "Invalid supplier ID")
		return
	}

	supplierId, err := strconv.Atoi(id)
	if err != nil {
		response.BadRequest(c, "Invalid supplier ID")
		return
	}

	// get supplier by id
	supplier, err := h.supplierService.GetByID(c.Request.Context(), supplierId)
	if err != nil {

		response.Server(c, err.Error())
		return
	}

	if supplier == nil {
		response.NotFound(c, "Supplier not found")
		return
	}

	response.Success(c, supplier, "Supplier retrieved successfully")
}

// update supplier by id
func (h *SupplierHandler) UpdateSupplierByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		response.BadRequest(c, "Invalid supplier ID")
		return
	}

	supplierId, err := strconv.Atoi(id)
	if err != nil {
		response.BadRequest(c, "Invalid supplier ID")
		return
	}

	var supplier supplierModel.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		// validation error when json body empty
		if err.Error() == "EOF" {
			response.BadRequest(c, "Missing request body. Please provide a valid JSON payload.")
			return
		}

		response.BadRequest(c, err.Error())
		return
	}

	// update supplier by id
	if err = h.supplierService.UpdateByID(c.Request.Context(), &supplier, supplierId); err != nil {
		// validation error
		if validateCreateOrUpdateSupplierError(err) {
			response.BadRequest(c, err.Error())
			return
		}

		if err.Error() == "supplier not found" {
			response.NotFound(c, err.Error())
			return
		}

		response.Server(c, err.Error())
		return
	}

	response.Success(c, supplier, "Supplier updated successfully")
}

// delete supplier by id
func (h *SupplierHandler) DeleteSupplierByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		response.BadRequest(c, "Invalid supplier ID")
		return
	}

	supplierId, err := strconv.Atoi(id)
	if err != nil {
		response.BadRequest(c, "Invalid supplier ID")
		return
	}

	if err := h.supplierService.DeleteByID(c.Request.Context(), supplierId); err != nil {
		if err.Error() == "supplier not found" {
			response.NotFound(c, err.Error())
			return
		}

		response.Server(c, err.Error())
		return
	}

	response.Success(c, nil, "Supplier deleted successfully")
}

// validate error when create or update supplier
func validateCreateOrUpdateSupplierError(err error) bool {
	validationErrors := []string{
		"name is required",
		"address is required",
		"contact is required",
	}

	for _, validationError := range validationErrors {
		if err.Error() == validationError {
			return true
		}
	}

	return false
}

func (h *SupplierHandler) RegisterSupplierRoutes(router *gin.RouterGroup) {
	supplierRoutes := router.Group("/supplier")
	{
		supplierRoutes.POST("", h.CreateSupplier)
		supplierRoutes.GET("", h.GetListSupplier)
		supplierRoutes.GET("/:id", h.GetSupplierByID)
		supplierRoutes.PUT("/:id", h.UpdateSupplierByID)
		supplierRoutes.DELETE("/:id", h.DeleteSupplierByID)
	}
}
