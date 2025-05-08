package location

import (
	"context"

	locationModel "ecosystem.garyle/service/internal/domain/model/wms/master-data/location"
)

type LocationRepository interface {
	Create(ctx context.Context, location *locationModel.Location) (*locationModel.Location, error)
	GetByID(ctx context.Context, id int) (*locationModel.Location, error)
	List(ctx context.Context, limit, page int) ([]*locationModel.Location, error)
	Count(ctx context.Context) (int, error)
	Update(ctx context.Context, location *locationModel.Location, id int) (*locationModel.Location, error)
	Delete(ctx context.Context, id int) error
}
