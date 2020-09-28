package domain

import (
	"time"
)

// Customer struct
type Customer struct {
	ID        int
	FirstName string
	LastName  string
	BirthDate time.Time
	Gender    string
	Email     string
	Address   string
}

// CustomerUsecase represent the Customer's usecases contract
type CustomerUsecase interface {
	Create(customer *Customer) error
	Update(customer *Customer) error
	GetByEmail(email string) (Customer, error)
	FetchAll() ([]Customer, error)
	Delete(customer *Customer) error
	GetByID(id int) (Customer, error)
}

// CustomerRepository represent the Customer's repository contract
type CustomerRepository interface {
	CloseConnection()
	Create(customer *Customer) error
	Update(customer *Customer) error
	GetByEmail(email string) (Customer, error)
	FetchAll() ([]Customer, error)
	Delete(customer *Customer) error
	GetByID(id int) (Customer, error)
}
