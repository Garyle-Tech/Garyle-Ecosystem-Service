package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	otaModel "ecosystem.garyle/service/internal/domain/model/ota"
	otaRepo "ecosystem.garyle/service/internal/domain/repository/ota"
)

type otaRepository struct {
	db *sql.DB
}

// NewOTARepository creates a new OTA repository
func NewOTARepository(db *sql.DB) otaRepo.OTARepository {
	return &otaRepository{
		db: db,
	}
}

func (r *otaRepository) Create(ctx context.Context, ota *otaModel.OTA) (*otaModel.OTA, error) {
	query := `
		INSERT INTO otas (app_id, version_name, version_code, url, release_notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	now := time.Now()
	ota.CreatedAt = now
	ota.UpdatedAt = now

	err := r.db.QueryRowContext(
		ctx,
		query,
		ota.AppID,
		ota.VersionName,
		ota.VersionCode,
		ota.URL,
		ota.ReleaseNotes,
		ota.CreatedAt,
		ota.UpdatedAt,
	).Scan(&ota.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to create OTA: %w", err)
	}

	return ota, nil
}

func (r *otaRepository) GetByAppID(ctx context.Context, appID string) (*otaModel.OTA, error) {
	query := `
		SELECT id, app_id, version_name, version_code, url, release_notes, created_at, updated_at
		FROM otas
		WHERE app_id = $1
		ORDER BY version_code DESC
		LIMIT 1
	`

	ota := &otaModel.OTA{}
	err := r.db.QueryRowContext(ctx, query, appID).Scan(
		&ota.ID,
		&ota.AppID,
		&ota.VersionName,
		&ota.VersionCode,
		&ota.URL,
		&ota.ReleaseNotes,
		&ota.CreatedAt,
		&ota.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get OTA by app ID: %w", err)
	}

	return ota, nil
}

func (r *otaRepository) List(ctx context.Context, limit, page int) ([]*otaModel.OTA, error) {
	query := `
		SELECT id, app_id, version_name, version_code, url, release_notes, created_at, updated_at
		FROM otas
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, (page-1)*limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list OTAs: %w", err)
	}

	defer rows.Close()

	var otas []*otaModel.OTA
	for rows.Next() {
		ota := &otaModel.OTA{}
		if err := rows.Scan(
			&ota.ID,
			&ota.AppID,
			&ota.VersionName,
			&ota.VersionCode,
			&ota.URL,
			&ota.ReleaseNotes,
			&ota.CreatedAt,
			&ota.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan OTA row: %w", err)
		}
		otas = append(otas, ota)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating OTA rows: %w", err)
	}

	if len(otas) == 0 || otas == nil {
		return []*otaModel.OTA{}, nil
	}

	return otas, nil
}

func (r *otaRepository) UpdateByAppID(ctx context.Context, ota *otaModel.OTA) error {
	query := `
		UPDATE otas
		SET version_name = $1, 
			version_code = $2, 
			url = $3, 
			release_notes = $4, 
			updated_at = $5
		WHERE app_id = $6
	`

	ota.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(
		ctx,
		query,
		ota.VersionName,
		ota.VersionCode,
		ota.URL,
		ota.ReleaseNotes,
		ota.UpdatedAt,
		ota.AppID,
	)

	if err != nil {
		return fmt.Errorf("failed to update OTA by app ID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("OTA for app ID %s not found", ota.AppID)
	}

	return nil
}

func (r *otaRepository) DeleteByAppID(ctx context.Context, appID string) error {
	query := `DELETE FROM otas WHERE app_id = $1`

	result, err := r.db.ExecContext(ctx, query, appID)
	if err != nil {
		return fmt.Errorf("failed to delete OTA by app ID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("OTA for app ID %s not found", appID)
	}

	return nil
}

func (r *otaRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM otas`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count OTAs: %w", err)
	}

	return count, nil
}
