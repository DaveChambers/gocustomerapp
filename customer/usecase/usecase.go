package usecase

import (
	"github.com/DaveChambers/gocustomerapp/domain"
)

type gormPostgresCustomerUsecase struct {
	customerRepo domain.CustomerRepository
}

// NewCustomerUsecase will create an object that represent the domain.CustomerUsecase interface
func NewCustomerUsecase(c domain.CustomerRepository) domain.CustomerUsecase {
	return &gormPostgresCustomerUsecase{customerRepo: c}
}

// Interface Implementation Functions:

func (m *gormPostgresCustomerUsecase) FetchCustomer(id int) (domain.Customer, error) {
	customer, err := m.customerRepo.GetByID(id)
	if err != nil {
		return domain.Customer{}, err
	}
	return customer, nil
}

func (m *gormPostgresCustomerUsecase) GetByID(id int) (domain.Customer, error) {
	customer, err := m.customerRepo.GetByID(id)
	if err != nil {
		return domain.Customer{}, err
	}
	return customer, nil
}

func (m *gormPostgresCustomerUsecase) GetByEmail(email string) (domain.Customer, error) {
	customer, err := m.customerRepo.GetByEmail(email)
	if err != nil {
		return domain.Customer{}, err
	}
	return customer, nil
}

func (m *gormPostgresCustomerUsecase) FetchAll() ([]domain.Customer, error) {
	customers, err := m.customerRepo.FetchAll()
	if err != nil {
		return []domain.Customer{}, err
	}
	return customers, nil
}

func (m *gormPostgresCustomerUsecase) Delete(customer *domain.Customer) error {
	err := m.customerRepo.Delete(customer)
	if err != nil {
		return err
	}
	return nil
}

func (m *gormPostgresCustomerUsecase) Create(customer *domain.Customer) error {
	err := m.customerRepo.Create(customer)
	if err != nil {
		return err
	}
	return nil
}

func (m *gormPostgresCustomerUsecase) Update(customer *domain.Customer) error {
	err := m.customerRepo.Update(customer)
	if err != nil {
		return err
	}
	return nil
}
