package category

import (
	"context"

	categoryModel "ecosystem.garyle/service/internal/domain/model/wms/master-data/category"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *categoryModel.Category) (*categoryModel.Category, error)
	List(ctx context.Context, limit, page int) ([]*categoryModel.Category, error)
	Count(ctx context.Context) (int, error)
	GetByID(ctx context.Context, id int) (*categoryModel.Category, error)
	Update(ctx context.Context, category *categoryModel.Category, id int) (*categoryModel.Category, error)
	Delete(ctx context.Context, id int) error
}
