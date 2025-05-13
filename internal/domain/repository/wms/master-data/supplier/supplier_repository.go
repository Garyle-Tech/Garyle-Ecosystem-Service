package supplier

import (
	"context"

	supplierModel "ecosystem.garyle/service/internal/domain/model/wms/master-data/supplier"
)

type SupplierRepository interface {
	Create(ctx context.Context, supplier *supplierModel.Supplier) (*supplierModel.Supplier, error)
	GetByID(ctx context.Context, id int) (*supplierModel.Supplier, error)
	List(ctx context.Context, limit, page int) ([]*supplierModel.Supplier, error)
	Count(ctx context.Context) (int, error)
	UpdateByID(ctx context.Context, supplier *supplierModel.Supplier, id int) error
	DeleteByID(ctx context.Context, id int) error
}
