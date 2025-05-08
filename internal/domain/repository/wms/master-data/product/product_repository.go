package product

import (
	"context"

	productModel "ecosystem.garyle/service/internal/domain/model/wms/master-data/product"
)

type ProductRepository interface {
	Create(ctx context.Context, product *productModel.Product) (*productModel.Product, error)
	GetByID(ctx context.Context, id int) (*productModel.Product, error)
	List(ctx context.Context, limit, page int) ([]*productModel.Product, error)
	Count(ctx context.Context) (int, error)
	UpdateByID(ctx context.Context, product *productModel.Product, id int) error
	DeleteByID(ctx context.Context, id int) error
}
