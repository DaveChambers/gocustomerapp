package delivery

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/DaveChambers/gocustomerapp/domain"
	"github.com/DaveChambers/gocustomerapp/domain/mocks"
	"github.com/DaveChambers/gocustomerapp/testhelper"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
)

func TestShowHandler(t *testing.T) {
	var mockCustomer domain.Customer
	err := faker.FakeData(&mockCustomer)
	assert.NoError(t, err)
	mockUCase := new(mocks.CustomerUsecase)
	num := int(mockCustomer.ID)

	// Mock the delivery layer's call to usecase layer
	mockUCase.On("GetByID", num).Return(mockCustomer, nil)

	req, err := http.NewRequest("GET", "/show/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	// Refactor this out when adding more Unit Tests...

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

	handler.ShowHandler(rec, req) // The actual function we are testing

	mockUCase.AssertExpectations(t)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))

}
