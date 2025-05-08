package product

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	productHandler "ecosystem.garyle/service/internal/app/api/wms/master-data/product"
	productService "ecosystem.garyle/service/internal/app/service/wms/master-data/product"
	productRepoPostgres "ecosystem.garyle/service/internal/infrastructure/database/wms/master-data/product"
)

var Module = fx.Module("product",
	fx.Provide(
		productRepoPostgres.NewProductRepository,
		productService.NewProductService,
		productHandler.NewProductHandler,
	),
)

// RegisterProductHandler registers product routes with the router group
func RegisterProductHandler(db *sql.DB, router *gin.RouterGroup) {
	repo := productRepoPostgres.NewProductRepository(db)
	service := productService.NewProductService(repo)
	handler := productHandler.NewProductHandler(service)

	handler.RegisterProductRoutes(router)
}
