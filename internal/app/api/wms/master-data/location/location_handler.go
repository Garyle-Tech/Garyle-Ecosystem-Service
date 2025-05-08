package location

import (
	"strconv"

	locationService "ecosystem.garyle/service/internal/app/service/wms/master-data/location"
	locationModel "ecosystem.garyle/service/internal/domain/model/wms/master-data/location"
	"ecosystem.garyle/service/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	locationService locationService.LocationService
}

func NewLocationHandler(locationService locationService.LocationService) *Handler {
	return &Handler{locationService: locationService}
}

// create location
// check validation is valid
func (h *Handler) CreateLocation(c *gin.Context) {
	// validate request body
	var requestBody locationModel.Location
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		if err.Error() == "EOF" {
			response.BadRequest(c, "Missing request body. Please provide a valid JSON payload.")
			return
		}

		response.BadRequest(c, err.Error())
		return
	}

	// create location
	createdLocation, err := h.locationService.Create(c.Request.Context(), &requestBody)
	if err != nil {
		if isValidationCreateOrUpdateLocationError(err) {
			response.BadRequest(c, err.Error())
			return
		}

		if err.Error() == "location already exists" {
			response.BadRequest(c, err.Error())
			return
		}

		if err.Error() == "pq: duplicate key value violates unique constraint \"locations_code_key\"" {
			response.BadRequest(c, "code already exists, please use another code")
			return
		}

		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, createdLocation, "Location created successfully")
}

func isValidationCreateOrUpdateLocationError(err error) bool {
	validationErrors := []string{
		"code is required",
		"zone is required",
		"type is required",
		"capacity is required",
	}

	for _, validationError := range validationErrors {
		if err.Error() == validationError {
			return true
		}
	}

	return false
}

// Get Locations
func (h *Handler) GetLocations(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))

	if limit <= 0 {
		limit = 10
	}

	if page <= 0 {
		page = 1
	}

	locations, err := h.locationService.List(c.Request.Context(), limit, page)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	total, err := h.locationService.Count(c.Request.Context())
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithPagination(c, locations, "Locations fetched successfully", page, limit, total)
}

// Get Location By ID
func (h *Handler) GetLocationByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		response.BadRequest(c, "Invalid location ID")
		return
	}

	locationID, err := strconv.Atoi(id)
	if err != nil {
		response.BadRequest(c, "Invalid location ID")
		return
	}

	location, err := h.locationService.GetByID(c.Request.Context(), locationID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if location == nil {
		response.NotFound(c, "Location not found")
		return
	}

	response.Success(c, location, "Location fetched successfully")
}

// Update Location
func (h *Handler) UpdateLocation(c *gin.Context) {
	id := c.Param("id")
	idConvert, err := strconv.Atoi(id)
	if err != nil || idConvert <= 0 {
		response.BadRequest(c, "Invalid location ID")
		return
	}

	var requestBody locationModel.Location
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		if err.Error() == "EOF" {
			response.BadRequest(c, "Missing request body. Please provide a valid JSON payload.")
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	updatedLocation, err := h.locationService.Update(c.Request.Context(), &requestBody, idConvert)
	if err != nil {
		if err.Error() == "location not found" {
			response.NotFound(c, err.Error())
			return
		}

		if err.Error() == "pq: duplicate key value violates unique constraint \"locations_code_key\"" {
			response.BadRequest(c, "code already exists, please use another code")
			return
		}

		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, updatedLocation, "Location updated successfully")
}

// Delete Location
func (h *Handler) DeleteLocation(c *gin.Context) {
	id := c.Param("id")
	idConvert, err := strconv.Atoi(id)
	if err != nil || idConvert <= 0 {
		response.BadRequest(c, "Invalid location ID")
		return
	}

	if err = h.locationService.Delete(c.Request.Context(), idConvert); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil, "Location deleted successfully")
}

// location routes
func (h *Handler) RegisterLocationRoutes(router *gin.RouterGroup) {
	locationRouter := router.Group("/location")
	{
		locationRouter.POST("", h.CreateLocation)
		locationRouter.GET("", h.GetLocations)
		locationRouter.GET("/:id", h.GetLocationByID)
		locationRouter.PATCH("/:id", h.UpdateLocation)
		locationRouter.DELETE("/:id", h.DeleteLocation)
	}
}
