package category

import (
	"database/sql"

	categoryHandler "ecosystem.garyle/service/internal/app/api/wms/master-data/category"
	categoryService "ecosystem.garyle/service/internal/app/service/wms/master-data/category"
	categoryRepoPostgres "ecosystem.garyle/service/internal/infrastructure/database/wms/master-data/category"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Module("category",
	fx.Provide(
		categoryRepoPostgres.NewCategoryRepository,
		categoryService.NewCategoryService,
		categoryHandler.NewCategoryHandler,
	),
)

func RegisterCategoryHandler(db *sql.DB, router *gin.RouterGroup) {
	repo := categoryRepoPostgres.NewCategoryRepository(db)
	service := categoryService.NewCategoryService(repo)
	handler := categoryHandler.NewCategoryHandler(service)

	handler.RegisterCategoryRoutes(router)
}
