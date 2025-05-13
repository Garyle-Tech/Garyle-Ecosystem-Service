package supplier

import (
	"context"
	"errors"

	supplierModel "ecosystem.garyle/service/internal/domain/model/wms/master-data/supplier"
	supplierRepo "ecosystem.garyle/service/internal/domain/repository/wms/master-data/supplier"
)

type SupplierService interface {
	Create(ctx context.Context, supplier *supplierModel.Supplier) (*supplierModel.Supplier, error)
	GetByID(ctx context.Context, id int) (*supplierModel.Supplier, error)
	List(ctx context.Context, limit, page int) ([]*supplierModel.Supplier, error)
	Count(ctx context.Context) (int, error)
	UpdateByID(ctx context.Context, supplier *supplierModel.Supplier, id int) error
	DeleteByID(ctx context.Context, id int) error
}

type supplierService struct {
	supplierRepo supplierRepo.SupplierRepository
}

func NewSupplierService(supplierRepo supplierRepo.SupplierRepository) SupplierService {
	return &supplierService{supplierRepo: supplierRepo}
}

// Create implements SupplierService.
func (s *supplierService) Create(ctx context.Context, supplier *supplierModel.Supplier) (*supplierModel.Supplier, error) {
	// validate supplier
	err := validateSupplier(supplier)
	if err != nil {
		return nil, err
	}

	// check if supplier already exists
	existingSupplier, err := s.supplierRepo.GetByID(ctx, supplier.ID)
	if err != nil {
		return nil, err
	}

	if existingSupplier != nil {
		return nil, errors.New("supplier already exists")
	}

	// create supplier
	createdSupplier, err := s.supplierRepo.Create(ctx, supplier)
	if err != nil {
		return nil, err
	}

	return createdSupplier, nil
}

// List implements SupplierService.
func (s *supplierService) List(ctx context.Context, limit int, page int) ([]*supplierModel.Supplier, error) {
	return s.supplierRepo.List(ctx, limit, page)
}

// Count implements SupplierService.
func (s *supplierService) Count(ctx context.Context) (int, error) {
	return s.supplierRepo.Count(ctx)
}

// GetByID implements SupplierService.
func (s *supplierService) GetByID(ctx context.Context, id int) (*supplierModel.Supplier, error) {
	if id <= 0 {
		return nil, errors.New("invalid supplier id")
	}

	return s.supplierRepo.GetByID(ctx, id)
}

// UpdateByID implements SupplierService.
func (s *supplierService) UpdateByID(ctx context.Context, supplier *supplierModel.Supplier, id int) error {
	if id <= 0 {
		return errors.New("invalid supplier id")
	}

	err := validateSupplier(supplier)
	if err != nil {
		return err
	}

	// check if supplier exists
	existingSupplier, err := s.supplierRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingSupplier == nil {
		return errors.New("supplier not found")
	}

	return s.supplierRepo.UpdateByID(ctx, supplier, id)
}

// DeleteByID implements SupplierService.
func (s *supplierService) DeleteByID(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid supplier id")
	}

	// check if supplier exists
	existingSupplier, err := s.supplierRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingSupplier == nil {
		return errors.New("supplier not found")
	}

	return s.supplierRepo.DeleteByID(ctx, id)
}

func validateSupplier(supplier *supplierModel.Supplier) error {
	if supplier.Name == "" {
		return errors.New("name is required")
	}

	if supplier.Address == "" {
		return errors.New("address is required")
	}

	if supplier.Contact == "" {
		return errors.New("contact is required")
	}

	return nil
}
