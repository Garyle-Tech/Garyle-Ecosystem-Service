package category

import (
	"strconv"

	categoryService "ecosystem.garyle/service/internal/app/service/wms/master-data/category"
	"ecosystem.garyle/service/internal/domain/model/wms/master-data/category"
	"ecosystem.garyle/service/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type categoryHandler struct {
	categoryService categoryService.CategoryService
}

func NewCategoryHandler(categoryService categoryService.CategoryService) *categoryHandler {
	return &categoryHandler{
		categoryService: categoryService,
	}
}

func (h *categoryHandler) CreateNewCategory(c *gin.Context) {
	var category category.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		// validation error when json body empty
		if err.Error() == "EOF" {
			response.BadRequest(c, "Missing request body. Please provide a valid JSON payload.")
			return
		}

		response.BadRequest(c, err.Error())
		return
	}

	// check existing category
	existingCategory, err := h.categoryService.GetByID(c.Request.Context(), category.ID)
	if err != nil {
		response.Server(c, err.Error())
		return
	}

	if existingCategory != nil {
		response.BadRequest(c, "category already exists")
		return
	}

	createdCategory, err := h.categoryService.Create(c.Request.Context(), &category)
	if err != nil {
		if validateCreateOrUpdateCategory(err) {
			response.BadRequest(c, err.Error())
			return
		}

		response.Server(c, err.Error())
		return
	}

	response.Success(c, createdCategory, "Category created successfully")
}

func (h *categoryHandler) GetAllCategories(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	categories, err := h.categoryService.List(c.Request.Context(), limit, page)
	if err != nil {
		response.Server(c, err.Error())
		return
	}

	total, err := h.categoryService.Count(c.Request.Context())
	if err != nil {
		response.Server(c, err.Error())
		return
	}

	if len(categories) == 0 || categories == nil {
		response.SuccessWithPagination(c, categories, "Category retrieved successfully", limit, page, total)
		return
	}

	response.SuccessWithPagination(c, categories, "Category retrieved successfully", limit, page, total)
}

func (h *categoryHandler) GetCategoryByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid category ID")
		return
	}

	category, err := h.categoryService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "category not found" {
			response.NotFound(c, err.Error())
			return
		}
		response.Server(c, err.Error())
		return
	}

	response.Success(c, category, "Category retrieved successfully")
}

func (h *categoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid category ID")
		return
	}

	var category category.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		if err.Error() == "EOF" {
			response.BadRequest(c, "Missing request body. Please provide a valid JSON payload.")
			return
		}

		response.BadRequest(c, err.Error())
		return
	}

	updatedCategory, err := h.categoryService.Update(c.Request.Context(), &category, id)
	if err != nil {
		if validateCreateOrUpdateCategory(err) {
			response.BadRequest(c, err.Error())
			return
		}

		if err.Error() == "category not found" {
			response.NotFound(c, err.Error())
			return
		}

		response.Server(c, err.Error())
		return
	}

	response.Success(c, updatedCategory, "Category updated successfully")
}

func (h *categoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid category ID")
		return
	}

	err = h.categoryService.Delete(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "category not found" {
			response.NotFound(c, err.Error())
			return
		}

		response.Server(c, err.Error())
		return
	}

	response.Success(c, nil, "Category deleted successfully")
}

func validateCreateOrUpdateCategory(err error) bool {
	validations := []string{
		"category name is required",
	}

	for _, validation := range validations {
		if err.Error() == validation {
			return true
		}
	}
	return false
}

func (h *categoryHandler) RegisterCategoryRoutes(router *gin.RouterGroup) {
	// group routes "categories"
	categoryRouter := router.Group("/categories")
	{
		categoryRouter.POST("/", h.CreateNewCategory)
		categoryRouter.GET("/", h.GetAllCategories)
		categoryRouter.GET("/:id", h.GetCategoryByID)
		categoryRouter.PUT("/:id", h.UpdateCategory)
		categoryRouter.DELETE("/:id", h.DeleteCategory)
	}
}
