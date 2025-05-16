package customer

import (
	"strconv"

	customerService "ecosystem.garyle/service/internal/app/service/wms/master-data/customer"
	customerModel "ecosystem.garyle/service/internal/domain/model/wms/master-data/customer"
	"ecosystem.garyle/service/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type customerHandler struct {
	customerService customerService.CustomerService
}

func NewCustomerHandler(customerService customerService.CustomerService) *customerHandler {
	return &customerHandler{customerService: customerService}
}

// Create new customer
func (h *customerHandler) CreateCustomer(c *gin.Context) {
	var customer customerModel.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		// validation error when json body empty
		if err.Error() == "EOF" {
			response.BadRequest(c, "Missing request body. Please provide a valid JSON payload.")
			return
		}

		response.BadRequest(c, err.Error())
		return
	}

	createdCustomer, err := h.customerService.Create(c.Request.Context(), &customer)
	if err != nil {
		if validateCreateCustomerError(err) {
			response.BadRequest(c, err.Error())
			return
		}

		if err.Error() == "customer already exists" {
			response.BadRequest(c, err.Error())
			return
		}

		response.Server(c, err.Error())
		return
	}

	response.Success(c, createdCustomer, "Customer created successfully")
}

func (h *customerHandler) GetListCustomer(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	if limit <= 0 {
		limit = 10
	}

	if page <= 0 {
		page = 1
	}

	customers, err := h.customerService.List(c.Request.Context(), limit, page)
	if err != nil {
		response.Server(c, err.Error())
		return
	}

	if len(customers) == 0 || customers == nil {
		response.Success(c, customers, "No customers found")
		return
	}

	totalCustomers, err := h.customerService.Count(c.Request.Context())
	if err != nil {
		response.Server(c, err.Error())
		return
	}

	response.SuccessWithPagination(c, customers, "Customers retrieved successfully", limit, page, totalCustomers)
}

// Get Customer by id
func (h *customerHandler) GetCustomerByID(c *gin.Context) {
	custoemrID := c.Param("id")
	if custoemrID == "" {
		response.BadRequest(c, "Customer ID is required")
		return
	}

	convertCustomerID, err := strconv.Atoi(custoemrID)
	if err != nil {
		response.BadRequest(c, "Invalid customer ID")
		return
	}

	customer, err := h.customerService.GetByID(c.Request.Context(), convertCustomerID)
	if err != nil {
		if err.Error() == "customer not found" {
			response.NotFound(c, err.Error())
			return
		}

		response.Server(c, err.Error())
		return
	}

	response.Success(c, customer, "Customer retrieved successfully")
}

// update customer by id
func (h *customerHandler) UpdateCustomerByID(c *gin.Context) {
	custoemrID := c.Param("id")
	if custoemrID == "" {
		response.BadRequest(c, "Customer ID is required")
		return
	}

	convertCustomerID, err := strconv.Atoi(custoemrID)
	if err != nil {
		response.BadRequest(c, "Invalid customer ID")
		return
	}

	var customer customerModel.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		if err.Error() == "EOF" {
			response.BadRequest(c, "Missing request body. Please provide a valid JSON payload.")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	err = h.customerService.UpdateByID(c.Request.Context(), &customer, convertCustomerID)
	if err != nil {
		if validateCreateCustomerError(err) {
			response.BadRequest(c, err.Error())
			return
		}
		if err.Error() == "customer not found" {
			response.NotFound(c, err.Error())
			return
		}

		response.Server(c, err.Error())
		return
	}

	response.Success(c, customer, "Customer updated successfully")
}

func (h *customerHandler) DeleteCustomerByID(c *gin.Context) {
	custoemrID := c.Param("id")
	if custoemrID == "" {
		response.BadRequest(c, "Customer ID is required")
		return
	}

	convertCustomerID, err := strconv.Atoi(custoemrID)
	if err != nil {
		response.BadRequest(c, "Invalid customer ID")
		return
	}

	err = h.customerService.DeleteByID(c.Request.Context(), convertCustomerID)
	if err != nil {
		if err.Error() == "customer not found" {
			response.NotFound(c, err.Error())
			return
		}

		response.Server(c, err.Error())
		return
	}

	response.Success(c, nil, "Customer deleted successfully")
}

func validateCreateCustomerError(err error) bool {
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

// Customer Route Groups
func (h *customerHandler) RegisterCustomerRoutes(router *gin.RouterGroup) {
	customerRouter := router.Group("/customers")
	{
		customerRouter.POST("", h.CreateCustomer)
		customerRouter.GET("", h.GetListCustomer)
		customerRouter.GET("/:id", h.GetCustomerByID)
		customerRouter.PUT("/:id", h.UpdateCustomerByID)
		customerRouter.DELETE("/:id", h.DeleteCustomerByID)
	}
}
