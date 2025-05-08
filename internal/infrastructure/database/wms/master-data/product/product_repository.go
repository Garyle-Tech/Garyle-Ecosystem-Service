package product

import (
	"context"
	"database/sql"
	"time"

	productModel "ecosystem.garyle/service/internal/domain/model/wms/master-data/product"
	productRepo "ecosystem.garyle/service/internal/domain/repository/wms/master-data/product"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) productRepo.ProductRepository {
	return &productRepository{db: db}
}

// endpoint create product
func (r *productRepository) Create(ctx context.Context, product *productModel.Product) (*productModel.Product, error) {
	// query insert product
	query := `
		INSERT INTO products (sku, name, description, unit, weight, dimension, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	// execute query
	err := r.db.QueryRowContext(
		ctx,
		query,
		product.Sku,
		product.Name,
		product.Description,
		product.Unit,
		product.Weight,
		product.Dimension,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID)

	if err != nil {
		return nil, err
	}

	return product, nil
}

// count total products
func (r *productRepository) Count(ctx context.Context) (int, error) {
	// query get total products
	query := `
		SELECT COUNT(*) FROM products
		WHERE deleted_at IS NULL
	`

	// execute query
	var total int
	err := r.db.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// get all products
func (r *productRepository) List(ctx context.Context, limit, page int) ([]*productModel.Product, error) {
	// query get all products
	query := `
		SELECT * FROM products
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	// execute query
	rows, err := r.db.QueryContext(ctx, query, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}

	// close rows
	defer rows.Close()

	// iterate rows
	products := []*productModel.Product{}
	for rows.Next() {
		var product productModel.Product
		err := rows.Scan(&product.ID, &product.Sku, &product.Name, &product.Description, &product.Unit, &product.Weight, &product.Dimension, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt)
		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(products) == 0 || products == nil {
		return []*productModel.Product{}, nil
	}

	// return products
	return products, nil
}

// get product by id
func (r *productRepository) GetByID(ctx context.Context, id int) (*productModel.Product, error) {
	// query get product by id
	query := `
		SELECT * FROM products
		WHERE id = $1 AND deleted_at IS NULL
	`

	// execute query
	var product productModel.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.Sku, &product.Name, &product.Description, &product.Unit, &product.Weight, &product.Dimension, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

// update product by id
func (r *productRepository) UpdateByID(ctx context.Context, product *productModel.Product, id int) error {
	// query update product by id
	query := `
		UPDATE products
		SET sku = $1, name = $2, description = $3, unit = $4, weight = $5, dimension = $6, updated_at = $7
		WHERE id = $8
	`

	// execute query
	_, err := r.db.ExecContext(ctx, query, product.Sku, product.Name, product.Description, product.Unit, product.Weight, product.Dimension, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

// delete product by id
func (r *productRepository) DeleteByID(ctx context.Context, id int) error {
	// query soft delete product by id
	query := `
		UPDATE products
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	// execute query
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}
