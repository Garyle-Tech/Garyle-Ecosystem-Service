package wms

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	categoryModule "ecosystem.garyle/service/internal/app/module/wms/master-data/category"
	customerModule "ecosystem.garyle/service/internal/app/module/wms/master-data/customer"
	locationModule "ecosystem.garyle/service/internal/app/module/wms/master-data/location"
	productModule "ecosystem.garyle/service/internal/app/module/wms/master-data/product"
	supplierModule "ecosystem.garyle/service/internal/app/module/wms/master-data/supplier"
)

var Module = fx.Module("wms",
	productModule.Module,
	locationModule.Module,
	supplierModule.Module,
	customerModule.Module,
	categoryModule.Module,
)

func RegisterWMSHandler(db *sql.DB, router *gin.RouterGroup) {
	wmsGroup := router.Group("/wms")
	masterDataGroup := wmsGroup.Group("/master-data")
	productModule.RegisterProductHandler(db, masterDataGroup)
	locationModule.RegisterLocationHandler(db, masterDataGroup)
	supplierModule.RegisterSupplierHandler(db, masterDataGroup)
	customerModule.RegisterCustomerHandler(db, masterDataGroup)
	categoryModule.RegisterCategoryHandler(db, masterDataGroup)
}
