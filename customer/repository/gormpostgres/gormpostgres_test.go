package gormpostgres_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/DaveChambers/gocustomerapp/customer/repository/gormpostgres"
	"github.com/DaveChambers/gocustomerapp/domain"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
)

var (
	testRepo domain.CustomerRepository
)

func TestMain(m *testing.M) {
	log.Println("Setup our repo before the tests...")
	testRepo = gormpostgres.NewCustomerRepository()
	exitVal := m.Run()
	log.Println("Delete all customers after the tests...")
	testRepo.Delete(&domain.Customer{})
	testRepo.CloseConnection()
	os.Exit(exitVal)
}

func TestCreate(t *testing.T) {
	// Create a customer in the DB:
	customer := &domain.Customer{
		FirstName: "Dave",
		LastName:  "Chambers",
		BirthDate: time.Date(1976, time.Month(2), 14, 0, 0, 0, 0, time.UTC),
		Gender:    "male",
		Email:     "dave@davechambers.co.uk",
		Address:   "Metsavälu, Põllküla, Lääne-Harju Vald, Harjumaa, 76712, Estonia"}
	err := testRepo.Create(customer)
	assert.NoError(t, err)

	// Read the customer from the DB to ensure the data is the data we wrote:
	customers, err := testRepo.FetchAll()
	assert.NoError(t, err)
	customerId := customers[0].ID
	fetchedCustomer, err := testRepo.GetByID(customerId)
	assert.NoError(t, err)
	assert.Equal(t, fetchedCustomer.FirstName, "Dave")
	assert.Equal(t, fetchedCustomer.LastName, "Chambers")
	assert.Equal(t, fetchedCustomer.BirthDate.Day(), 14)
	assert.Equal(t, fetchedCustomer.BirthDate.Month(), time.Month(2))
	assert.Equal(t, fetchedCustomer.BirthDate.Year(), 1976)
	assert.Equal(t, fetchedCustomer.Gender, "male")
	assert.Equal(t, fetchedCustomer.Email, "dave@davechambers.co.uk")
	assert.Equal(t, fetchedCustomer.Address, "Metsavälu, Põllküla, Lääne-Harju Vald, Harjumaa, 76712, Estonia")
}

func TestUpdate(t *testing.T) {
	// Read the customer from the DB:
	customers, err := testRepo.FetchAll()
	assert.NoError(t, err)
	customerId := customers[0].ID
	fetchedCustomer, err := testRepo.GetByID(customerId)
	assert.NoError(t, err)

	// Update the customer's email address and Address:
	updatedCustomer := &domain.Customer{
		ID:        fetchedCustomer.ID,
		FirstName: "Dave",
		LastName:  "Chambers",
		BirthDate: time.Date(1976, time.Month(2), 14, 0, 0, 0, 0, time.UTC),
		Gender:    "male",
		Email:     "dave@davechambers.com",
		Address:   "Some cool place",
	}

	updateErr := testRepo.Update(updatedCustomer)
	assert.NoError(t, updateErr)

	// Read the customer from the DB to ensure the data is the data we wrote:
	customers2, err2 := testRepo.FetchAll()
	assert.NoError(t, err2)
	customerId2 := customers2[0].ID
	fetchedCustomer, err3 := testRepo.GetByID(customerId2)
	assert.NoError(t, err3)
	assert.Equal(t, fetchedCustomer.FirstName, "Dave")
	assert.Equal(t, fetchedCustomer.LastName, "Chambers")
	assert.Equal(t, fetchedCustomer.BirthDate.Day(), 14)
	assert.Equal(t, fetchedCustomer.BirthDate.Month(), time.Month(2))
	assert.Equal(t, fetchedCustomer.BirthDate.Year(), 1976)
	assert.Equal(t, fetchedCustomer.Gender, "male")
	assert.Equal(t, fetchedCustomer.Email, "dave@davechambers.com")
	assert.Equal(t, fetchedCustomer.Address, "Some cool place")
}

func TestGetByEmail(t *testing.T) {
	_, err := testRepo.GetByEmail("dave@davechambers.co.uk")
	if err == nil {
		t.Error("Expected an error since a user with this email address doesn't exist")
	}

	customerCoDotCom, err := testRepo.GetByEmail("dave@davechambers.com")

	assert.NoError(t, err)
	assert.Equal(t, customerCoDotCom.FirstName, "Dave")
	assert.Equal(t, customerCoDotCom.LastName, "Chambers")
	assert.Equal(t, customerCoDotCom.BirthDate.Day(), 14)
	assert.Equal(t, customerCoDotCom.BirthDate.Month(), time.Month(2))
	assert.Equal(t, customerCoDotCom.BirthDate.Year(), 1976)
	assert.Equal(t, customerCoDotCom.Gender, "male")
	assert.Equal(t, customerCoDotCom.Email, "dave@davechambers.com")
	assert.Equal(t, customerCoDotCom.Address, "Some cool place")
}

func TestFetchAll(t *testing.T) {
	// Read the customers from the DB:
	customers, err := testRepo.FetchAll()
	assert.NoError(t, err)
	assert.Equal(t, len(customers), 1)
}

func TestDelete(t *testing.T) {
	// Read the customer from the DB:
	customers, err := testRepo.FetchAll()
	assert.NoError(t, err)
	assert.Equal(t, len(customers), 1)
	customerId := customers[0].ID

	// Delete them
	customer := &domain.Customer{ID: customerId}
	err2 := testRepo.Delete(customer)
	assert.NoError(t, err2)

	// Assert there are none left:
	// Read the customer from the DB:
	customers2, err := testRepo.FetchAll()
	assert.NoError(t, err)
	assert.Equal(t, len(customers2), 0)
}

func TestGetByID(t *testing.T) {
	// Create a customer in the DB:
	customer := &domain.Customer{
		FirstName: "Jackie",
		LastName:  "Jones",
		BirthDate: time.Date(1986, time.Month(10), 11, 0, 0, 0, 0, time.UTC),
		Gender:    "female",
		Email:     "jackie@jones.co.uk",
		Address:   "320 Junctions Street, Howechester, Kentucky, 93689, Namibia"}
	err := testRepo.Create(customer)
	assert.NoError(t, err)

	// Read the customer from the DB to ensure the data is the data we wrote:
	customers, err := testRepo.FetchAll()
	assert.NoError(t, err)
	customerId := customers[0].ID
	fetchedCustomer, err := testRepo.GetByID(customerId)
	assert.NoError(t, err)
	assert.Equal(t, fetchedCustomer.FirstName, "Jackie")
	assert.Equal(t, fetchedCustomer.LastName, "Jones")
	assert.Equal(t, fetchedCustomer.BirthDate.Day(), 11)
	assert.Equal(t, fetchedCustomer.BirthDate.Month(), time.Month(10))
	assert.Equal(t, fetchedCustomer.BirthDate.Year(), 1986)
	assert.Equal(t, fetchedCustomer.Gender, "female")
	assert.Equal(t, fetchedCustomer.Email, "jackie@jones.co.uk")
	assert.Equal(t, fetchedCustomer.Address, "320 Junctions Street, Howechester, Kentucky, 93689, Namibia")
}
