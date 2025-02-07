// Code generated by mockery v2.2.2. DO NOT EDIT.

package mocks

import (
	domain "github.com/DaveChambers/gocustomerapp/domain"
	mock "github.com/stretchr/testify/mock"
)

// CustomerUsecase is an autogenerated mock type for the CustomerUsecase type
type CustomerUsecase struct {
	mock.Mock
}

// Create provides a mock function with given fields: customer
func (_m *CustomerUsecase) Create(customer *domain.Customer) error {
	ret := _m.Called(customer)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Customer) error); ok {
		r0 = rf(customer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: customer
func (_m *CustomerUsecase) Delete(customer *domain.Customer) error {
	ret := _m.Called(customer)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Customer) error); ok {
		r0 = rf(customer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FetchAll provides a mock function with given fields:
func (_m *CustomerUsecase) FetchAll() ([]domain.Customer, error) {
	ret := _m.Called()

	var r0 []domain.Customer
	if rf, ok := ret.Get(0).(func() []domain.Customer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Customer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByEmail provides a mock function with given fields: email
func (_m *CustomerUsecase) GetByEmail(email string) (domain.Customer, error) {
	ret := _m.Called(email)

	var r0 domain.Customer
	if rf, ok := ret.Get(0).(func(string) domain.Customer); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Get(0).(domain.Customer)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: id
func (_m *CustomerUsecase) GetByID(id int) (domain.Customer, error) {
	ret := _m.Called(id)

	var r0 domain.Customer
	if rf, ok := ret.Get(0).(func(int) domain.Customer); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(domain.Customer)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: customer
func (_m *CustomerUsecase) Update(customer *domain.Customer) error {
	ret := _m.Called(customer)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Customer) error); ok {
		r0 = rf(customer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
