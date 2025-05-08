package location

import (
	"context"
	"errors"

	locationModel "ecosystem.garyle/service/internal/domain/model/wms/master-data/location"
	locationRepo "ecosystem.garyle/service/internal/domain/repository/wms/master-data/location"
)

type LocationService interface {
	Create(ctx context.Context, location *locationModel.Location) (*locationModel.Location, error)
	GetByID(ctx context.Context, id int) (*locationModel.Location, error)
	List(ctx context.Context, limit, page int) ([]*locationModel.Location, error)
	Count(ctx context.Context) (int, error)
	Update(ctx context.Context, location *locationModel.Location, id int) (*locationModel.Location, error)
	Delete(ctx context.Context, id int) error
}

type locationService struct {
	repo locationRepo.LocationRepository
}

func NewLocationService(repo locationRepo.LocationRepository) LocationService {
	return &locationService{repo: repo}
}

// Create implements LocationService.
func (l *locationService) Create(ctx context.Context, location *locationModel.Location) (*locationModel.Location, error) {
	// validate location
	err := validatorCreateOrUpdateLocation(location)
	if err != nil {
		return nil, err
	}

	// check if location already exists
	existingLocation, err := l.repo.GetByID(ctx, location.ID)
	if err != nil {
		return nil, err
	}

	if existingLocation != nil {
		return nil, errors.New("location already exists")
	}

	return l.repo.Create(ctx, location)
}

func validatorCreateOrUpdateLocation(location *locationModel.Location) error {
	if location.Code == "" {
		return errors.New("code is required")
	}

	if location.Zone == "" {
		return errors.New("zone is required")
	}

	if location.Type == "" {
		return errors.New("type is required")
	}

	if location.Capacity <= 0 {
		return errors.New("capacity is required")
	}

	return nil
}

// List implements LocationService.
func (l *locationService) List(ctx context.Context, limit int, page int) ([]*locationModel.Location, error) {
	return l.repo.List(ctx, limit, page)
}

// Count implements LocationService.
func (l *locationService) Count(ctx context.Context) (int, error) {
	return l.repo.Count(ctx)
}

// GetByID implements LocationService.
func (l *locationService) GetByID(ctx context.Context, id int) (*locationModel.Location, error) {
	if id <= 0 {
		return nil, errors.New("invalid location id")
	}

	return l.repo.GetByID(ctx, id)
}

// Update implements LocationService.
func (l *locationService) Update(ctx context.Context, location *locationModel.Location, id int) (*locationModel.Location, error) {
	if id <= 0 {
		return nil, errors.New("invalid location id")
	}

	// check if location already exists
	existingLocation, err := l.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingLocation == nil {
		return nil, errors.New("location not found")
	}

	return l.repo.Update(ctx, location, id)
}

// Delete implements LocationService.
func (l *locationService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid location id")
	}

	// check if location already exists
	existingLocation, err := l.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingLocation == nil {
		return errors.New("location not found")
	}

	return l.repo.Delete(ctx, id)
}
