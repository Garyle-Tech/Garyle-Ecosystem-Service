package product

import (
	"context"
	"errors"

	productModel "ecosystem.garyle/service/internal/domain/model/wms/master-data/product"
	productRepo "ecosystem.garyle/service/internal/domain/repository/wms/master-data/product"
)

type ProductService interface {
	Create(ctx context.Context, product *productModel.Product) (*productModel.Product, error)
	GetByID(ctx context.Context, id int) (*productModel.Product, error)
	List(ctx context.Context, limit, page int) ([]*productModel.Product, error)
	Count(ctx context.Context) (int, error)
	UpdateByID(ctx context.Context, product *productModel.Product, id int) error
	DeleteByID(ctx context.Context, id int) error
}

type productService struct {
	productRepo productRepo.ProductRepository
}

func NewProductService(productRepo productRepo.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) Create(ctx context.Context, product *productModel.Product) (*productModel.Product, error) {
	// validate product
	err := validateProduct(product)
	if err != nil {
		return nil, err
	}

	// check if data product exists
	existingProduct, err := s.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	if existingProduct != nil {
		return nil, errors.New("product already exists")
	}

	return s.productRepo.Create(ctx, product)
}

// validate product
func validateProduct(product *productModel.Product) error {
	if product.Sku == "" {
		return errors.New("sku is required")
	}

	if product.Name == "" {
		return errors.New("name is required")
	}

	if product.Unit == "" {
		return errors.New("unit is required")
	}

	if product.Weight <= 0 {
		return errors.New("weight is required")
	}

	if product.Dimension == "" {
		return errors.New("dimension is required")
	}

	return nil
}

func (s *productService) GetByID(ctx context.Context, id int) (*productModel.Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid product id")
	}

	return s.productRepo.GetByID(ctx, id)
}

func (s *productService) List(ctx context.Context, limit, page int) ([]*productModel.Product, error) {
	return s.productRepo.List(ctx, limit, page)
}

func (s *productService) Count(ctx context.Context) (int, error) {
	return s.productRepo.Count(ctx)
}

func (s *productService) UpdateByID(ctx context.Context, product *productModel.Product, id int) error {
	// validate product
	err := validateProduct(product)
	if err != nil {
		return err
	}

	// check if data product exists
	existingProduct, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingProduct == nil {
		return errors.New("product not found")
	}

	return s.productRepo.UpdateByID(ctx, product, id)
}

func (s *productService) DeleteByID(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid product id")
	}

	// check if data product exists
	existingProduct, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingProduct == nil {
		return errors.New("product not found")
	}

	return s.productRepo.DeleteByID(ctx, id)
}
