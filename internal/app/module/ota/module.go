package ota

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"ecosystem.garyle/service/internal/app/api/ota"
	otaService "ecosystem.garyle/service/internal/app/service/ota"
	"ecosystem.garyle/service/internal/infrastructure/database"
)

var Module = fx.Module("ota",
	fx.Provide(
		database.NewOTARepository,
		otaService.NewService,
		ota.NewHandler,
	),
)

// RegisterOTAHandlers registers OTA routes with the router group
func RegisterOTAHandlers(db *sql.DB, router *gin.RouterGroup) {
	repo := database.NewOTARepository(db)
	service := otaService.NewService(repo)
	handler := ota.NewHandler(service)

	handler.RegisterRoutes(router)
}
