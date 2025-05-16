package category

import (
	"context"
	"errors"

	categoryModel "ecosystem.garyle/service/internal/domain/model/wms/master-data/category"
	categoryRepo "ecosystem.garyle/service/internal/domain/repository/wms/master-data/category"
)

type CategoryService interface {
	Create(ctx context.Context, category *categoryModel.Category) (*categoryModel.Category, error)
	List(ctx context.Context, limit, page int) ([]*categoryModel.Category, error)
	Count(ctx context.Context) (int, error)
	GetByID(ctx context.Context, id int) (*categoryModel.Category, error)
	Update(ctx context.Context, category *categoryModel.Category, id int) (*categoryModel.Category, error)
	Delete(ctx context.Context, id int) error
}

type categoryService struct {
	categoryRepository categoryRepo.CategoryRepository
}

func NewCategoryService(categoryRepository categoryRepo.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}

// Create implements CategoryService.
func (c *categoryService) Create(ctx context.Context, category *categoryModel.Category) (*categoryModel.Category, error) {
	if err := validateCreateOrUpdateCategory(category); err != nil {
		return nil, err
	}

	existingCategory, err := c.categoryRepository.GetByID(ctx, category.ID)
	if err != nil {
		return nil, err
	}
	if existingCategory != nil {
		return nil, errors.New("category already exists")
	}

	newCategory, err := c.categoryRepository.Create(ctx, category)
	if err != nil {
		return nil, err
	}

	return newCategory, nil
}

// List implements CategoryService.
func (c *categoryService) List(ctx context.Context, limit int, page int) ([]*categoryModel.Category, error) {
	categories, err := c.categoryRepository.List(ctx, limit, page)
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 || categories == nil {
		return nil, errors.New("no categories found")
	}

	return categories, nil
}

// Count implements CategoryService.
func (c *categoryService) Count(ctx context.Context) (int, error) {
	count, err := c.categoryRepository.Count(ctx)
	if err != nil {
		return 0, err
	}

	if count == 0 {
		return 0, nil
	}

	return count, nil
}

// GetByID implements CategoryService.
func (c *categoryService) GetByID(ctx context.Context, id int) (*categoryModel.Category, error) {
	if id <= 0 {
		return nil, errors.New("invalid category ID")
	}

	category, err := c.categoryRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New("category not found")
	}

	return category, nil
}

// Update implements CategoryService.
func (c *categoryService) Update(ctx context.Context, category *categoryModel.Category, id int) (*categoryModel.Category, error) {
	if id <= 0 {
		return nil, errors.New("invalid category ID")
	}

	if err := validateCreateOrUpdateCategory(category); err != nil {
		return nil, err
	}

	existingCategory, err := c.categoryRepository.GetByID(ctx, category.ID)
	if err != nil {
		return nil, err
	}
	if existingCategory == nil {
		return nil, errors.New("category not found")
	}

	updatedCategory, err := c.categoryRepository.Update(ctx, category, id)
	if err != nil {
		return nil, err
	}

	return updatedCategory, nil
}

// Delete implements CategoryService.
func (c *categoryService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid category ID")
	}

	existingCategory, err := c.categoryRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingCategory == nil {
		return errors.New("category not found")
	}

	err = c.categoryRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func validateCreateOrUpdateCategory(category *categoryModel.Category) error {
	if category.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
