package ota

import (
	"strconv"

	"github.com/gin-gonic/gin"

	otaService "ecosystem.garyle/service/internal/app/service/ota"
	otaModel "ecosystem.garyle/service/internal/domain/model/ota"
	"ecosystem.garyle/service/pkg/utils/response"
)

type Handler struct {
	otaService otaService.Service
}

func NewHandler(otaService otaService.Service) *Handler {
	return &Handler{
		otaService: otaService,
	}
}

func (h *Handler) CreateOTA(c *gin.Context) {
	var ota otaModel.OTA
	if err := c.ShouldBindJSON(&ota); err != nil {
		if err.Error() == "EOF" {
			response.BadRequest(c, "Missing request body. Please provide a valid JSON payload.")
			return
		}

		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.otaService.Create(c.Request.Context(), &ota)
	if err != nil {
		if isValidationCreateOTAError(err) {
			response.BadRequest(c, err.Error())
			return
		}

		response.Server(c, err.Error())
		return
	}

	response.Success(c, result, "OTA created successfully")
}

func isValidationCreateOTAError(err error) bool {
	validationErrors := []string{
		"app ID is required",
		"version name is required",
		"version code must be a positive number",
		"URL is required",
		"an OTA update already exists for this app",
	}

	for _, validationErr := range validationErrors {
		if err.Error() == validationErr {
			return true
		}
	}
	return false
}

func (h *Handler) GetOTA(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		response.BadRequest(c, "invalid app_id")
		return
	}

	ota, err := h.otaService.GetByAppID(c.Request.Context(), appID)
	if err != nil {
		response.Server(c, err.Error())
		return
	}

	if ota == nil {
		response.NotFound(c, "OTA not found for this app")
		return
	}

	response.Success(c, ota, "OTA retrieved successfully")
}

func (h *Handler) ListOTAs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	if limit <= 0 {
		limit = 10
	}

	if page <= 0 {
		page = 1
	}

	otas, err := h.otaService.List(c.Request.Context(), limit, page)
	if err != nil {
		response.Server(c, err.Error())
		return
	}

	total, err := h.otaService.Count(c.Request.Context())
	if err != nil {
		response.Server(c, err.Error())
		return
	}

	response.SuccessWithPagination(c, otas, "OTAs retrieved successfully", page, limit, total)
}

func (h *Handler) UpdateOTA(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		response.BadRequest(c, "invalid app_id")
		return
	}

	var ota otaModel.OTA
	if err := c.ShouldBindJSON(&ota); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	ota.AppID = appID
	if err := h.otaService.UpdateByAppID(c.Request.Context(), &ota); err != nil {
		response.Server(c, err.Error())
		return
	}

	response.Success(c, ota, "OTA updated successfully")
}

func (h *Handler) DeleteOTA(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		response.BadRequest(c, "invalid app_id")
		return
	}

	if err := h.otaService.DeleteByAppID(c.Request.Context(), appID); err != nil {
		response.Server(c, err.Error())
		return
	}

	response.Success(c, nil, "OTA deleted successfully")
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	otaRoutes := router.Group("/ota")
	{
		otaRoutes.POST("", h.CreateOTA)
		otaRoutes.GET("", h.ListOTAs)
		otaRoutes.GET("/detail", h.GetOTA)
		otaRoutes.PUT("/edit", h.UpdateOTA)
		otaRoutes.DELETE("/delete", h.DeleteOTA)
	}
}
