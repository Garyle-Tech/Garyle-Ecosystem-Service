package customer

import (
	"context"

	customerModel "ecosystem.garyle/service/internal/domain/model/wms/master-data/customer"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *customerModel.Customer) (*customerModel.Customer, error)
	GetByID(ctx context.Context, id int) (*customerModel.Customer, error)
	List(ctx context.Context, limit, page int) ([]*customerModel.Customer, error)
	Count(ctx context.Context) (int, error)
	UpdateByID(ctx context.Context, customer *customerModel.Customer, id int) error
	DeleteByID(ctx context.Context, id int) error
}
