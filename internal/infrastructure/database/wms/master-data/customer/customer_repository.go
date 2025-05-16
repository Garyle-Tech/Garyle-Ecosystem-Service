package customer

import (
	"context"
	"database/sql"
	"time"

	customerModel "ecosystem.garyle/service/internal/domain/model/wms/master-data/customer"
	customerRepo "ecosystem.garyle/service/internal/domain/repository/wms/master-data/customer"
)

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) customerRepo.CustomerRepository {
	return &customerRepository{db: db}
}

func (c *customerRepository) Count(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*) FROM customers
		WHERE deleted_at IS NULL
	`

	var count int
	if err := c.db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (c *customerRepository) Create(ctx context.Context, customer *customerModel.Customer) (*customerModel.Customer, error) {
	query := `
		INSERT INTO customers (name, address, contact, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()
	customer.CreatedAt = now
	customer.UpdatedAt = now

	var newCustomer customerModel.Customer
	if err := c.db.QueryRowContext(
		ctx,
		query,
		customer.Name,
		customer.Address,
		customer.Contact,
		customer.CreatedAt,
		customer.UpdatedAt,
	).Scan(&newCustomer.ID); err != nil {
		return nil, err
	}

	return &newCustomer, nil
}

// DeleteByID implements customer.CustomerRepository.
func (c *customerRepository) DeleteByID(ctx context.Context, id int) error {
	query := `
		UPDATE customers
		SET deleted_at = $1, updated_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	if _, err := c.db.ExecContext(ctx, query, time.Now(), id); err != nil {
		return err
	}

	return nil
}

// GetByID implements customer.CustomerRepository.
func (c *customerRepository) GetByID(ctx context.Context, id int) (*customerModel.Customer, error) {
	query := `
		SELECT * FROM customers
		WHERE id = $1 AND deleted_at IS NULL
	`

	var customer customerModel.Customer
	if err := c.db.QueryRowContext(ctx, query, id).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Address,
		&customer.Contact,
		&customer.CreatedAt,
		&customer.UpdatedAt,
		&customer.DeletedAt,
	); err != nil {
		// data not found
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &customer, nil
}

// List implements customer.CustomerRepository.
func (c *customerRepository) List(ctx context.Context, limit int, page int) ([]*customerModel.Customer, error) {
	query := `
		SELECT * FROM customers
		WHERE deleted_at IS NULL
		LIMIT $1 OFFSET $2
	`

	offset := (page - 1) * limit
	rows, err := c.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	customers := []*customerModel.Customer{}
	// jika rows masih bisa di next (masih ada data di page selanjutnya dengan limit tertentu)
	for rows.Next() {
		var customer customerModel.Customer
		if err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Address,
			&customer.Contact,
			&customer.CreatedAt,
			&customer.UpdatedAt,
			&customer.DeletedAt,
		); err != nil {
			return nil, err
		}

		customers = append(customers, &customer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(customers) == 0 || customers == nil {
		return []*customerModel.Customer{}, nil
	}

	return customers, nil
}

// UpdateByID implements customer.CustomerRepository.
func (c *customerRepository) UpdateByID(ctx context.Context, customer *customerModel.Customer, id int) error {
	query := `
		UPDATE customers
		SET name = $1, address = $2, contact = $3, updated_at = $4
		WHERE id = $5 AND deleted_at IS NULL
	`

	if _, err := c.db.ExecContext(ctx, query, customer.Name, customer.Address, customer.Contact, customer.UpdatedAt, id); err != nil {
		return err
	}

	return nil
}
