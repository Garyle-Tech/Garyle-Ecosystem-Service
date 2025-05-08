package ota

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"ecosystem.garyle/service/internal/app/api/ota"
	otaService "ecosystem.garyle/service/internal/app/service/ota"
	otaRepoPostgres "ecosystem.garyle/service/internal/infrastructure/database/ota"
)

var Module = fx.Module("ota",
	fx.Provide(
		otaRepoPostgres.NewOTARepository,
		otaService.NewService,
		ota.NewHandler,
	),
)

// RegisterOTAHandler registers OTA routes with the router group
func RegisterOTAHandler(db *sql.DB, router *gin.RouterGroup) {
	repo := otaRepoPostgres.NewOTARepository(db)
	service := otaService.NewService(repo)
	handler := ota.NewHandler(service)

	handler.RegisterRoutes(router)
}
