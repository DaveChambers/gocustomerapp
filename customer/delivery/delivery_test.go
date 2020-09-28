package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/DaveChambers/gocustomerapp/domain"
	"github.com/DaveChambers/gocustomerapp/domain/mocks"
	"github.com/DaveChambers/gocustomerapp/errors"
	"github.com/DaveChambers/gocustomerapp/testhelper"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
)

func getHandlerWithMockUsecase(mockUCase *mocks.CustomerUsecase) *Handler {
	pathToRoot, _ := testhelper.GetRootPath() // tmpl path will need adjusting
	createEditPath := path.Join(pathToRoot, "tmpl", "create-edit.html")
	showPath := path.Join(pathToRoot, "tmpl", "show.html")
	searchPath := path.Join(pathToRoot, "tmpl", "search.html")
	notFoundPath := path.Join(pathToRoot, "tmpl", "404.html")
	templates := template.Must(template.ParseFiles(createEditPath, showPath, searchPath, notFoundPath))

	handler := &Handler{
		templates: templates,
		uc:        mockUCase,
	}
	return handler
}

func TestCheckEmailHandler(t *testing.T) {
	var mockCustomer domain.Customer
	err := faker.FakeData(&mockCustomer)
	assert.NoError(t, err)
	mockUCase := new(mocks.CustomerUsecase)
	email := mockCustomer.Email

	t.Run("Success", func(t *testing.T) {

		// Mock the delivery layer's call to usecase layer
		mockUCase.On("GetByEmail", email).Return(mockCustomer, nil).Once()

		req, err := http.NewRequest("GET", fmt.Sprintf("/checkemail?email=%s", mockCustomer.Email), strings.NewReader(""))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()

		handler := getHandlerWithMockUsecase(mockUCase)

		handler.CheckEmailHandler(rec, req) // The actual function we are testing

		mockUCase.AssertExpectations(t)

		assert.Equal(t, http.StatusOK, rec.Code)

		b, err := json.Marshal(mockCustomer)
		if err != nil {
			fmt.Println(err)
		}

		//Confirm the returned json is what we expected
		assert.JSONEq(t, string(b), rec.Body.String(), "Response body differs")

	})
	t.Run("EmailNotFoundError", func(t *testing.T) {

		// Mock the delivery layer's call to usecase layer
		mockUCase.On("GetByEmail", email).Return(domain.Customer{}, &errors.EmailNotFoundError{}).Once()

		req, err := http.NewRequest("GET", fmt.Sprintf("/checkemail?email=%s", mockCustomer.Email), strings.NewReader(""))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()

		handler := getHandlerWithMockUsecase(mockUCase)

		handler.CheckEmailHandler(rec, req) // The actual function we are testing

		mockUCase.AssertExpectations(t)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		b, err := json.Marshal(domain.Customer{})
		if err != nil {
			fmt.Println(err)
		}

		//Confirm the returned json is what we expected
		assert.JSONEq(t, string(b), rec.Body.String(), "Response body differs")
	})
}

func TestFetchCustomersHandler(t *testing.T) {
	mockCustomer1 := &domain.Customer{
		FirstName: "Dave",
		LastName:  "Chambers",
		BirthDate: time.Date(1976, time.Month(2), 14, 0, 0, 0, 0, time.UTC),
		Gender:    "male",
		Email:     "dave@davechambers.co.uk",
		Address:   "Metsavälu, Põllküla, Lääne-Harju Vald, Harjumaa, 76712, Estonia"}

	mockCustomer2 := &domain.Customer{
		FirstName: "Oda",
		LastName:  "Senger",
		BirthDate: time.Date(1997, time.Month(2), 27, 0, 0, 0, 0, time.UTC),
		Gender:    "male",
		Email:     "sylvesterbergnaum@stark.biz",
		Address:   "8245 Skyway Street, Vilmamouth, Delaware, 89520, Timor-Leste"}

	mockCustomer3 := &domain.Customer{
		FirstName: "Marcelina",
		LastName:  "Denesik",
		BirthDate: time.Date(1991, time.Month(12), 25, 0, 0, 0, 0, time.UTC),
		Gender:    "female",
		Email:     "dave@davechambers.co.uk",
		Address:   "97159 Valley Street, Hillaryside, Arizona, 10744, Nauru"}

	mockUCase := new(mocks.CustomerUsecase)

	customers := []domain.Customer{*mockCustomer1, *mockCustomer2, *mockCustomer3}

	// Mock the delivery layer's call to usecase layer
	mockUCase.On("FetchAll").Return(customers, nil).Once()

	req, err := http.NewRequest("GET", "/fetchcustomers", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	handler := getHandlerWithMockUsecase(mockUCase)

	handler.FetchCustomersHandler(rec, req) // The actual function we are testing

	mockUCase.AssertExpectations(t)

	assert.Equal(t, http.StatusOK, rec.Code)

	expected := string(`[{"ID":0,"FirstName":"Dave","LastName":"Chambers","BirthDate":"1976-02-14T00:00:00Z","Gender":"male","Email":"dave@davechambers.co.uk","Address":"Metsavälu, Põllküla, Lääne-Harju Vald, Harjumaa, 76712, Estonia"},{"ID":0,"FirstName":"Oda","LastName":"Senger","BirthDate":"1997-02-27T00:00:00Z","Gender":"male","Email":"sylvesterbergnaum@stark.biz","Address":"8245 Skyway Street, Vilmamouth, Delaware, 89520, Timor-Leste"},{"ID":0,"FirstName":"Marcelina","LastName":"Denesik","BirthDate":"1991-12-25T00:00:00Z","Gender":"female","Email":"dave@davechambers.co.uk","Address":"97159 Valley Street, Hillaryside, Arizona, 10744, Nauru"}]`)
	//Confirm the returned json is what we expected
	assert.JSONEq(t, expected, rec.Body.String(), "Response body differs")
}

func TestDeleteCustomerHandler(t *testing.T) {
	mockCustomer := &domain.Customer{ID: 1}

	mockUCase := new(mocks.CustomerUsecase)

	// Mock the delivery layer's call to usecase layer
	mockUCase.On("Delete", mockCustomer).Return(nil).Once()

	req, err := http.NewRequest("POST", "/deletecustomer", bytes.NewBuffer([]byte(`{"id" : 1}`)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	handler := getHandlerWithMockUsecase(mockUCase)

	handler.DeleteCustomerHandler(rec, req) // The actual function we are testing

	mockUCase.AssertExpectations(t)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCreateHandler(t *testing.T) {
	mockUCase := new(mocks.CustomerUsecase)
	req, err := http.NewRequest("GET", "/create/", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	handler := getHandlerWithMockUsecase(mockUCase)

	handler.CreateHandler(rec, req) // The actual function we are testing

	mockUCase.AssertExpectations(t)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), "First name:")
}

func TestShowHandler(t *testing.T) {
	var mockCustomer domain.Customer
	err := faker.FakeData(&mockCustomer)
	assert.NoError(t, err)
	mockUCase := new(mocks.CustomerUsecase)
	num := int(mockCustomer.ID)

	t.Run("Success", func(t *testing.T) {

		// Mock the delivery layer's call to usecase layer
		mockUCase.On("GetByID", num).Return(mockCustomer, nil).Once()

		req, err := http.NewRequest("GET", "/show/"+strconv.Itoa(num), strings.NewReader(""))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()

		handler := getHandlerWithMockUsecase(mockUCase)

		handler.ShowHandler(rec, req) // The actual function we are testing

		mockUCase.AssertExpectations(t)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))

	})
	t.Run("404", func(t *testing.T) {

		// Mock the delivery layer's call to usecase layer
		mockUCase.On("GetByID", num).Return(domain.Customer{}, &errors.CustomerNotFoundError{}).Once()

		req, err := http.NewRequest("GET", "/show/"+strconv.Itoa(num), strings.NewReader(""))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()

		handler := getHandlerWithMockUsecase(mockUCase)

		handler.ShowHandler(rec, req) // The actual function we are testing

		mockUCase.AssertExpectations(t)

		assert.Contains(t, rec.Body.String(), "404 - Customer Not Found!")
		assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))

	})
}

func TestEditHandler(t *testing.T) {
	var mockCustomer domain.Customer
	err := faker.FakeData(&mockCustomer)
	assert.NoError(t, err)
	mockUCase := new(mocks.CustomerUsecase)
	num := int(mockCustomer.ID)

	t.Run("Success", func(t *testing.T) {

		// Mock the delivery layer's call to usecase layer
		mockUCase.On("GetByID", num).Return(mockCustomer, nil).Once()

		req, err := http.NewRequest("GET", "/edit/"+strconv.Itoa(num), strings.NewReader(""))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()

		handler := getHandlerWithMockUsecase(mockUCase)

		handler.EditHandler(rec, req) // The actual function we are testing

		mockUCase.AssertExpectations(t)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))

	})
	t.Run("404", func(t *testing.T) {

		// Mock the delivery layer's call to usecase layer
		mockUCase.On("GetByID", num).Return(domain.Customer{}, &errors.CustomerNotFoundError{}).Once()

		req, err := http.NewRequest("GET", "/edit/"+strconv.Itoa(num), strings.NewReader(""))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()

		handler := getHandlerWithMockUsecase(mockUCase)

		handler.EditHandler(rec, req) // The actual function we are testing

		mockUCase.AssertExpectations(t)

		assert.Contains(t, rec.Body.String(), "404 - Customer Not Found!")
		assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))

	})
}

func TestSearchHandler(t *testing.T) {
	mockUCase := new(mocks.CustomerUsecase)
	req, err := http.NewRequest("GET", "/search/", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	handler := getHandlerWithMockUsecase(mockUCase)

	handler.SearchHandler(rec, req) // The actual function we are testing

	mockUCase.AssertExpectations(t)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), "Search Customers")
}

func TestSaveHandler(t *testing.T) {
	mockCustomer := &domain.Customer{
		FirstName: "Dave",
		LastName:  "Chambers",
		BirthDate: time.Date(1976, time.Month(2), 14, 0, 0, 0, 0, time.UTC),
		Gender:    "male",
		Email:     "dave@davechambers.co.uk",
		Address:   "Metsavälu, Põllküla, Lääne-Harju Vald, Harjumaa, 76712, Estonia"}

	mockUCase := new(mocks.CustomerUsecase)

	// Mock the delivery layer's call to usecase layer
	mockUCase.On("Create", mockCustomer).Return(nil).Once()

	readerParams := fmt.Sprintf("fname=%v&lname=%v&birthdate=1976-02-14&gender=%v&email=%v&address=%v", mockCustomer.FirstName, mockCustomer.LastName, mockCustomer.Gender, mockCustomer.Email, mockCustomer.Address)
	req, err := http.NewRequest("POST", "/save", strings.NewReader(readerParams))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	handler := getHandlerWithMockUsecase(mockUCase)

	handler.SaveHandler(rec, req) // The actual function we are testing

	mockUCase.AssertExpectations(t)

	assert.Equal(t, 302, rec.Code)
}

func TestSaveEditHandler(t *testing.T) {
	mockCustomer := &domain.Customer{
		ID:        1,
		FirstName: "Dave",
		LastName:  "Chambers",
		BirthDate: time.Date(1976, time.Month(2), 14, 0, 0, 0, 0, time.UTC),
		Gender:    "male",
		Email:     "dave@davechambers.com",
		Address:   "Some cool place"}

	mockUCase := new(mocks.CustomerUsecase)

	// Mock the delivery layer's call to usecase layer
	mockUCase.On("GetByID", mockCustomer.ID).Return(*mockCustomer, nil).Once()

	// Mock the delivery layer's call to usecase layer
	mockUCase.On("Update", mockCustomer).Return(nil).Once()

	url := fmt.Sprintf("/saveedit/%v", mockCustomer.ID)
	readerParams := fmt.Sprintf("fname=%v&lname=%v&birthdate=1976-02-14&gender=%v&email=%v&address=%v", mockCustomer.FirstName, mockCustomer.LastName, mockCustomer.Gender, mockCustomer.Email, mockCustomer.Address)
	req, err := http.NewRequest("POST", url, strings.NewReader(readerParams))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	handler := getHandlerWithMockUsecase(mockUCase)

	handler.SaveEditHandler(rec, req) // The actual function we are testing

	mockUCase.AssertExpectations(t)

	assert.Equal(t, 302, rec.Code)
}
