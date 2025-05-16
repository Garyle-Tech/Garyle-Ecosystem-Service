package customer

import (
	"database/sql"

	customerHandler "ecosystem.garyle/service/internal/app/api/wms/master-data/customer"
	customerService "ecosystem.garyle/service/internal/app/service/wms/master-data/customer"
	customerRepoPostgres "ecosystem.garyle/service/internal/infrastructure/database/wms/master-data/customer"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Module("customer",
	fx.Provide(
		customerRepoPostgres.NewCustomerRepository,
		customerService.NewCustomerService,
		customerHandler.NewCustomerHandler,
	),
)

func RegisterCustomerHandler(db *sql.DB, router *gin.RouterGroup) {
	repo := customerRepoPostgres.NewCustomerRepository(db)
	service := customerService.NewCustomerService(repo)
	handler := customerHandler.NewCustomerHandler(service)

	handler.RegisterCustomerRoutes(router)
}
