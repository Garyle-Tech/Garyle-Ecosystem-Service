package category

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"ecosystem.garyle/service/internal/domain/model/wms/master-data/category"
	categoryRepo "ecosystem.garyle/service/internal/domain/repository/wms/master-data/category"
)

type categoryRepository struct {
	sql *sql.DB
}

func NewCategoryRepository(sql *sql.DB) categoryRepo.CategoryRepository {
	return &categoryRepository{sql: sql}
}

// Create implements category.CategoryRepository.
func (c *categoryRepository) Create(ctx context.Context, category *category.Category) (*category.Category, error) {
	query := `
		INSERT INTO categories (name, parent_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now

	err := c.sql.QueryRowContext(ctx, query, category.Name, category.ParentID, category.CreatedAt, category.UpdatedAt).Scan(&category.ID)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// List implements category.CategoryRepository.
func (c *categoryRepository) List(ctx context.Context, limit int, page int) ([]*category.Category, error) {
	query := `
		SELECT * FROM categories
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	offset := (page - 1) * limit
	rows, err := c.sql.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var categories []*category.Category
	for rows.Next() {
		var category category.Category
		err := rows.Scan(&category.ID, &category.Name, &category.ParentID, &category.CreatedAt, &category.UpdatedAt, &category.DeletedAt)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return []*category.Category{}, nil
	}

	return categories, nil
}

// Count implements category.CategoryRepository.
func (c *categoryRepository) Count(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*) FROM categories
		WHERE deleted_at IS NULL
	`

	var count int
	err := c.sql.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetByID implements category.CategoryRepository.
func (c *categoryRepository) GetByID(ctx context.Context, id int) (*category.Category, error) {
	query := `
		SELECT * FROM categories
		WHERE id = ? AND deleted_at IS NULL
	`
	var category category.Category
	err := c.sql.QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name, &category.ParentID, &category.CreatedAt, &category.UpdatedAt, &category.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("category not found") // Not found
		}
		return nil, err
	}

	return &category, nil
}

// Update implements category.CategoryRepository.
func (c *categoryRepository) Update(ctx context.Context, category *category.Category, id int) (*category.Category, error) {
	query := `
		UPDATE categories
		SET name = $1, parent_id = $2, updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	_, err := c.sql.ExecContext(ctx, query, category.Name, category.ParentID, category.UpdatedAt, id)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// Delete implements category.CategoryRepository.
func (c *categoryRepository) Delete(ctx context.Context, id int) error {
	query := `
		UPDATE categories
		SET deleted_at = $2, updated_at = $2
		WHERE id = $1 AND deleted_at IS NULL
	`

	if _, err := c.sql.ExecContext(ctx, query, id, time.Now()); err != nil {
		return err
	}

	return nil
}
