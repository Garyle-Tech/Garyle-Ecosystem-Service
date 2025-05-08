package location

import (
	"database/sql"

	locationHandler "ecosystem.garyle/service/internal/app/api/wms/master-data/location"
	locationService "ecosystem.garyle/service/internal/app/service/wms/master-data/location"
	locationRepo "ecosystem.garyle/service/internal/infrastructure/database/wms/master-data/location"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Module("location",
	fx.Provide(
		locationService.NewLocationService,
		locationRepo.NewLocationRepository,
		locationHandler.NewLocationHandler,
	),
)

func RegisterLocationHandler(db *sql.DB, router *gin.RouterGroup) {
	repo := locationRepo.NewLocationRepository(db)
	service := locationService.NewLocationService(repo)
	handler := locationHandler.NewLocationHandler(service)

	handler.RegisterLocationRoutes(router)
}
