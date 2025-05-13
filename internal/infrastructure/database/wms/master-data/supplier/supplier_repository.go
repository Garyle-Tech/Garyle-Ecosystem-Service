package supplier

import (
	"context"
	"database/sql"
	"time"

	"ecosystem.garyle/service/internal/domain/model/wms/master-data/supplier"
	supplierRepo "ecosystem.garyle/service/internal/domain/repository/wms/master-data/supplier"
)

type supplierRepository struct {
	db *sql.DB
}

func NewSupplierRepository(db *sql.DB) supplierRepo.SupplierRepository {
	return &supplierRepository{db: db}
}

// Create implements supplier.SupplierRepository.
func (s *supplierRepository) Create(ctx context.Context, supplier *supplier.Supplier) (*supplier.Supplier, error) {
	query := `
		INSERT INTO suppliers (name, address, contact, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()
	supplier.CreatedAt = now
	supplier.UpdatedAt = now

	err := s.db.QueryRowContext(
		ctx,
		query,
		supplier.Name,
		supplier.Address,
		supplier.Contact,
		supplier.CreatedAt,
		supplier.UpdatedAt,
	).Scan(&supplier.ID)

	if err != nil {
		return nil, err
	}

	return supplier, nil
}

// List implements supplier.SupplierRepository.
func (s *supplierRepository) List(ctx context.Context, limit int, page int) ([]*supplier.Supplier, error) {
	query := `
		SELECT * FROM suppliers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var offset int = (page - 1) * limit
	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	suppliers := []*supplier.Supplier{}
	for rows.Next() {
		var supplier supplier.Supplier
		err := rows.Scan(&supplier.ID, &supplier.Name, &supplier.Address, &supplier.Contact, &supplier.CreatedAt, &supplier.UpdatedAt, &supplier.DeletedAt)
		if err != nil {
			return nil, err
		}

		suppliers = append(suppliers, &supplier)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(suppliers) == 0 || suppliers == nil {
		return []*supplier.Supplier{}, nil
	}

	return suppliers, nil
}

// Count implements supplier.SupplierRepository.
func (s *supplierRepository) Count(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*) FROM suppliers
		WHERE deleted_at IS NULL
	`

	var total int
	err := s.db.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// GetByID implements supplier.SupplierRepository.
func (s *supplierRepository) GetByID(ctx context.Context, id int) (*supplier.Supplier, error) {
	query := `
		SELECT * FROM suppliers
		WHERE id = $1 AND deleted_at IS NULL
	`

	var supplier supplier.Supplier
	err := s.db.QueryRowContext(ctx, query, id).Scan(&supplier.ID, &supplier.Name, &supplier.Address, &supplier.Contact, &supplier.CreatedAt, &supplier.UpdatedAt, &supplier.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &supplier, nil
}

// UpdateByID implements supplier.SupplierRepository.
func (s *supplierRepository) UpdateByID(ctx context.Context, supplier *supplier.Supplier, id int) error {
	query := `
		UPDATE suppliers
		SET name = $1, address = $2, contact = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := s.db.ExecContext(ctx, query, supplier.Name, supplier.Address, supplier.Contact, supplier.UpdatedAt, id)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByID implements supplier.SupplierRepository.
func (s *supplierRepository) DeleteByID(ctx context.Context, id int) error {
	query := `
		UPDATE suppliers
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	_, err := s.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}
