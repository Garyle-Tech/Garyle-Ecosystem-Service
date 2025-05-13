package supplier

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	supplierHandler "ecosystem.garyle/service/internal/app/api/wms/master-data/supplier"
	supplierService "ecosystem.garyle/service/internal/app/service/wms/master-data/supplier"
	supplierRepoPostgres "ecosystem.garyle/service/internal/infrastructure/database/wms/master-data/supplier"
)

var Module = fx.Module("supplier",
	fx.Provide(
		supplierRepoPostgres.NewSupplierRepository,
		supplierService.NewSupplierService,
		supplierHandler.NewSupplierHandler,
	),
)

func RegisterSupplierHandler(db *sql.DB, router *gin.RouterGroup) {
	repo := supplierRepoPostgres.NewSupplierRepository(db)
	service := supplierService.NewSupplierService(repo)
	handler := supplierHandler.NewSupplierHandler(service)

	handler.RegisterSupplierRoutes(router)
}
