package location

import (
	"context"
	"database/sql"
	"time"

	"ecosystem.garyle/service/internal/domain/model/wms/master-data/location"
	locationRepo "ecosystem.garyle/service/internal/domain/repository/wms/master-data/location"
)

type locationRepository struct {
	db *sql.DB
}

func NewLocationRepository(db *sql.DB) locationRepo.LocationRepository {
	return &locationRepository{db: db}
}

// Create implements location.LocationRepository.
func (l *locationRepository) Create(ctx context.Context, location *location.Location) (*location.Location, error) {
	query := `
		INSERT INTO locations (code, zone, type, capacity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, code, zone, type, capacity, created_at, updated_at
	`

	now := time.Now()
	location.CreatedAt = now
	location.UpdatedAt = now

	err := l.db.QueryRowContext(ctx, query, location.Code, location.Zone, location.Type, location.Capacity, location.CreatedAt, location.UpdatedAt).Scan(
		&location.ID,
		&location.Code,
		&location.Zone,
		&location.Type,
		&location.Capacity,
		&location.CreatedAt,
		&location.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return location, nil
}

// List implements location.LocationRepository.
func (l *locationRepository) List(ctx context.Context, limit int, page int) ([]*location.Location, error) {
	query := `
		SELECT id, code, zone, type, capacity, created_at, updated_at
		FROM locations
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	offset := (page - 1) * limit

	rows, err := l.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var locations []*location.Location
	for rows.Next() {
		var location location.Location
		err := rows.Scan(
			&location.ID,
			&location.Code,
			&location.Zone,
			&location.Type,
			&location.Capacity,
			&location.CreatedAt,
			&location.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		locations = append(locations, &location)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(locations) == 0 || locations == nil {
		return []*location.Location{}, nil
	}

	return locations, nil
}

// Count implements location.LocationRepository.
func (l *locationRepository) Count(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*) FROM locations
		WHERE deleted_at IS NULL
	`

	var count int
	err := l.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetByID implements location.LocationRepository.
func (l *locationRepository) GetByID(ctx context.Context, id int) (*location.Location, error) {
	query := `
		SELECT id, code, zone, type, capacity, created_at, updated_at
		FROM locations
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var location location.Location
	err := l.db.QueryRowContext(ctx, query, id).Scan(
		&location.ID,
		&location.Code,
		&location.Zone,
		&location.Type,
		&location.Capacity,
		&location.CreatedAt,
		&location.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &location, nil
}

// Update implements location.LocationRepository.
func (l *locationRepository) Update(ctx context.Context, dataLocation *location.Location, id int) (*location.Location, error) {
	// Allow updating all fields
	query := `
		UPDATE locations
		SET code = $1, zone = $2, type = $3, capacity = $4, updated_at = $5
		WHERE id = $6
		AND deleted_at IS NULL
		RETURNING id, code, zone, type, capacity, created_at, updated_at
	`

	now := time.Now()

	var updatedLocation location.Location
	if err := l.db.QueryRowContext(ctx, query, dataLocation.Code, dataLocation.Zone, dataLocation.Type, dataLocation.Capacity, now, id).Scan(
		&updatedLocation.ID,
		&updatedLocation.Code,
		&updatedLocation.Zone,
		&updatedLocation.Type,
		&updatedLocation.Capacity,
		&updatedLocation.CreatedAt,
		&updatedLocation.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &updatedLocation, nil
}

// Delete implements location.LocationRepository.
func (l *locationRepository) Delete(ctx context.Context, id int) error {
	query := `
		UPDATE locations
		SET deleted_at = $1
		WHERE id = $2
		AND deleted_at IS NULL
	`

	now := time.Now()

	if _, err := l.db.ExecContext(ctx, query, now, id); err != nil {
		return err
	}

	return nil
}
