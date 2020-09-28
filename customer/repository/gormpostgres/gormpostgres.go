package gormpostgres

import (
	"github.com/DaveChambers/gocustomerapp/dbconnection"
	"github.com/DaveChambers/gocustomerapp/domain"
	"github.com/DaveChambers/gocustomerapp/errors"
	"github.com/jinzhu/gorm"
)

type gormPostgresCustomerRepository struct {
	conn *gorm.DB
}

// NewCustomerRepository will create an object that represent the domain.CustomerRepository interface
func NewCustomerRepository() domain.CustomerRepository {
	db := dbconnection.Connect()
	return &gormPostgresCustomerRepository{db}
}

// Interface Implementation Functions:

func (m *gormPostgresCustomerRepository) CloseConnection() {
	m.conn.Close()
}

func (m *gormPostgresCustomerRepository) Create(customer *domain.Customer) error {
	if err := m.conn.Create(&customer).Error; err != nil {
		return err
	}
	return nil
}

func (m *gormPostgresCustomerRepository) Update(customer *domain.Customer) error {
	// No need for mux.Lock() etc since GORM performs write operations inside a
	// Transaction by default https://gorm.io/docs/transactions.html#Disable-Default-Transaction
	m.conn.Save(&customer)
	return nil
}

func (m *gormPostgresCustomerRepository) GetByEmail(email string) (domain.Customer, error) {
	var cust domain.Customer
	// Note SQL Injection is not possible due to FE checking for a valid email
	// address AND Gorm would escape a dangerous input anyway in this case: https://gorm.io/docs/security.html
	if err := m.conn.Where("email = ?", email).First(&cust).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return 404 to front end meaning that the email is available
			return domain.Customer{}, &errors.EmailNotFoundError{}
		}
	}
	return cust, nil
}

func (m *gormPostgresCustomerRepository) FetchAll() ([]domain.Customer, error) {
	var customers []domain.Customer
	if err := m.conn.Find(&customers).Error; err != nil {
		return []domain.Customer{}, err
	}
	return customers, nil
}

func (m *gormPostgresCustomerRepository) Delete(customer *domain.Customer) error {
	if err := m.conn.Delete(&customer).Error; err != nil {
		return err
	}
	return nil
}

func (m *gormPostgresCustomerRepository) GetByID(id int) (domain.Customer, error) {
	var cust domain.Customer
	if err := m.conn.First(&cust, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Customer{}, &errors.CustomerNotFoundError{}
		}
		return domain.Customer{}, err
	}
	return cust, nil
}
