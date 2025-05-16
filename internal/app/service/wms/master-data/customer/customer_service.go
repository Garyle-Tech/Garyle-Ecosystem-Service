package customer

import (
	"context"
	"errors"

	customerModel "ecosystem.garyle/service/internal/domain/model/wms/master-data/customer"
	customerRepo "ecosystem.garyle/service/internal/domain/repository/wms/master-data/customer"
)

type CustomerService interface {
	Create(ctx context.Context, customer *customerModel.Customer) (*customerModel.Customer, error)
	GetByID(ctx context.Context, id int) (*customerModel.Customer, error)
	List(ctx context.Context, limit int, page int) ([]*customerModel.Customer, error)
	UpdateByID(ctx context.Context, customer *customerModel.Customer, id int) error
	DeleteByID(ctx context.Context, id int) error
	Count(ctx context.Context) (int, error)
}

type customerService struct {
	customerRepository customerRepo.CustomerRepository
}

func NewCustomerService(customerRepository customerRepo.CustomerRepository) CustomerService {
	return &customerService{customerRepository: customerRepository}
}

// Create implements CustomerService.
func (c *customerService) Create(ctx context.Context, customer *customerModel.Customer) (*customerModel.Customer, error) {
	err := validateCustomer(customer)
	if err != nil {
		return nil, err
	}

	existingCustomer, err := c.GetByID(ctx, customer.ID)
	if err != nil {
		return nil, err
	}

	if existingCustomer != nil {
		return nil, errors.New("customer already exists")
	}

	createdCustomer, err := c.customerRepository.Create(ctx, customer)
	if err != nil {
		return nil, err
	}

	return createdCustomer, nil
}

// DeleteByID implements CustomerService.
func (c *customerService) DeleteByID(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid customer id")
	}

	existingCustomer, err := c.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingCustomer == nil {
		return errors.New("customer not found")
	}

	return c.customerRepository.DeleteByID(ctx, existingCustomer.ID)
}

// GetByID implements CustomerService.
func (c *customerService) GetByID(ctx context.Context, id int) (*customerModel.Customer, error) {
	if id <= 0 {
		return nil, errors.New("invalid customer id")
	}

	customer, err := c.customerRepository.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("customer not found")
		}

		return nil, err
	}

	return customer, nil
}

// List implements CustomerService.
func (c *customerService) List(ctx context.Context, limit int, page int) ([]*customerModel.Customer, error) {
	customers, err := c.customerRepository.List(ctx, limit, page)
	if err != nil {
		return nil, err
	}

	if len(customers) == 0 || customers == nil {
		return []*customerModel.Customer{}, nil
	}

	return customers, nil
}

// UpdateByID implements CustomerService.
func (c *customerService) UpdateByID(ctx context.Context, customer *customerModel.Customer, id int) error {
	if id <= 0 {
		return errors.New("invalid customer id")
	}

	err := validateCustomer(customer)
	if err != nil {
		return err
	}

	existingCustomer, err := c.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existingCustomer == nil {
		return errors.New("customer not found")
	}

	return c.customerRepository.UpdateByID(ctx, customer, id)
}

// Count implements CustomerService.
func (c *customerService) Count(ctx context.Context) (int, error) {
	return c.customerRepository.Count(ctx)
}

func validateCustomer(customer *customerModel.Customer) error {
	if customer.Name == "" {
		return errors.New("name is required")
	}

	if customer.Address == "" {
		return errors.New("address is required")
	}

	if customer.Contact == "" {
		return errors.New("contact is required")
	}

	return nil
}
