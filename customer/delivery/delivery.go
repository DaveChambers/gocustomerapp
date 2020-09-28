package delivery

import (
	"encoding/json"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/DaveChambers/gocustomerapp/domain"
	"github.com/DaveChambers/gocustomerapp/errors"
	"github.com/DaveChambers/gocustomerapp/testhelper"
	"github.com/gorilla/mux"

	"fmt"
	"html/template"
	"log"
	"net/http"
)

// The Handler struct that holds the templates and a CustomerUsecase
type Handler struct {
	templates *template.Template
	uc        domain.CustomerUsecase
}

func newRouter(h Handler) *mux.Router {
	router := mux.NewRouter()

	// API Routes:
	router.HandleFunc("/checkemail", h.CheckEmailHandler)
	router.HandleFunc("/fetchcustomers", h.FetchCustomersHandler)
	router.HandleFunc("/deletecustomer", h.DeleteCustomerHandler)

	// Web Routes:
	router.HandleFunc("/create/", h.CreateHandler)
	router.HandleFunc("/show/{id}", h.ShowHandler)
	router.HandleFunc("/edit/{id}", h.EditHandler)
	router.HandleFunc("/search/", h.SearchHandler)
	router.HandleFunc("/save", h.SaveHandler)
	router.HandleFunc("/saveedit/{id}", h.SaveEditHandler)

	staticDir := "/static/"
	router.
		PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))
	return router
}

// NewHandler will initialize the handler with injected CustomerUsecase
func NewHandler(uc domain.CustomerUsecase) {
	pathToRoot, _ := testhelper.GetRootPath() // We may or may not be testing.  If we are tmpl path will need adjusting
	createEditPath := path.Join(pathToRoot, "tmpl", "create-edit.html")
	showPath := path.Join(pathToRoot, "tmpl", "show.html")
	searchPath := path.Join(pathToRoot, "tmpl", "search.html")
	notFoundPath := path.Join(pathToRoot, "tmpl", "404.html")
	templates := template.Must(template.ParseFiles(createEditPath, showPath, searchPath, notFoundPath))

	handler := &Handler{
		templates: templates,
		uc:        uc,
	}

	router := newRouter(*handler)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}

// Helpers
func fetchCustomer(h *Handler, w http.ResponseWriter, r *http.Request) (domain.Customer, error) {
	id, err := strconv.Atoi(filepath.Base(r.URL.String()))
	if err != nil {
		return domain.Customer{}, err
	}
	customer, err := h.uc.GetByID(id)
	if err != nil {
		return domain.Customer{}, err
	}
	return customer, nil
}

func displaySavedOrEditedRecord(id int, w http.ResponseWriter, r *http.Request) {
	showURL := fmt.Sprintf("/show/%d", id)
	http.Redirect(w, r, showURL, http.StatusFound)
}

func getDOB(birthdate string) time.Time {
	ymd := strings.Split(birthdate, "-")
	var ymdAsInts = []int{}
	for _, i := range ymd {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		ymdAsInts = append(ymdAsInts, j)
	}
	dob := time.Date(ymdAsInts[0], time.Month(ymdAsInts[1]), ymdAsInts[2], 0, 0, 0, 0, time.UTC)
	return dob
}

func executeTemplate(templates *template.Template, tmpl string, customer domain.Customer, w http.ResponseWriter, r *http.Request) {
	templateErr := templates.ExecuteTemplate(w, tmpl+".html", customer)
	if templateErr != nil {
		http.Error(w, templateErr.Error(), http.StatusBadRequest)
	}
}

func handleErrors(templates *template.Template, err error, w http.ResponseWriter, r *http.Request) {
	switch err.(type) {
	// Test for the errors we care about and handle those specific cases...
	case *errors.CustomerNotFoundError:
		executeTemplate(templates, "404", domain.Customer{}, w, r)
	case *errors.EmailNotFoundError:
		// Return 404 to front end meaning that the email is available
		w.WriteHeader(http.StatusNotFound)
	default:
		// Any error types we don't specifically look out for default
		// to serving a HTTP 500
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}

// CreateHandler ...
func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(h.templates, "create-edit", domain.Customer{}, w, r)
}

// SearchHandler ...
func (h *Handler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(h.templates, "search", domain.Customer{}, w, r)
}

// ShowHandler ...
func (h *Handler) ShowHandler(w http.ResponseWriter, r *http.Request) {
	customer, err := fetchCustomer(h, w, r)
	if err != nil {
		handleErrors(h.templates, err, w, r)
	} else {
		executeTemplate(h.templates, "show", customer, w, r)
	}
}

// EditHandler ...
func (h *Handler) EditHandler(w http.ResponseWriter, r *http.Request) {
	customer, err := fetchCustomer(h, w, r)
	if err != nil {
		handleErrors(h.templates, err, w, r)
	}
	executeTemplate(h.templates, "create-edit", customer, w, r)
}

// SaveHandler ...
func (h *Handler) SaveHandler(w http.ResponseWriter, r *http.Request) {
	customer := &domain.Customer{FirstName: r.FormValue("fname"),
		LastName:  r.FormValue("lname"),
		BirthDate: getDOB(r.FormValue("birthdate")),
		Gender:    r.FormValue("gender"),
		Email:     r.FormValue("email"),
		Address:   r.FormValue("address")}
	err := h.uc.Create(customer)
	if err != nil {
		handleErrors(h.templates, err, w, r)
	}
	displaySavedOrEditedRecord(customer.ID, w, r)
}

// SaveEditHandler ...
func (h *Handler) SaveEditHandler(w http.ResponseWriter, r *http.Request) {
	customer, err := fetchCustomer(h, w, r)
	if err != nil {
		handleErrors(h.templates, err, w, r)
	}
	// Here we just set all the values in case any changed...
	customer.FirstName = r.FormValue("fname")
	customer.LastName = r.FormValue("lname")
	customer.BirthDate = getDOB(r.FormValue("birthdate"))
	customer.Gender = r.FormValue("gender")
	customer.Email = r.FormValue("email")
	customer.Address = r.FormValue("address")
	theErr := h.uc.Update(&customer)
	if theErr != nil {
		handleErrors(h.templates, theErr, w, r)
	}
	displaySavedOrEditedRecord(customer.ID, w, r)
}

// FetchCustomersHandler ...
func (h *Handler) FetchCustomersHandler(w http.ResponseWriter, r *http.Request) {
	customers, err := h.uc.FetchAll()
	if err != nil {
		handleErrors(h.templates, err, w, r)
	}
	js, err := json.Marshal(customers)
	if err != nil {
		handleErrors(h.templates, err, w, r)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// DeleteCustomerHandler ...
func (h *Handler) DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	var customer = domain.Customer{}
	// Try to decode the request body (it only has the ID present) into a Customer struct.
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		handleErrors(h.templates, err, w, r)
	}
	// Perform the deletion:
	theErr := h.uc.Delete(&customer)
	if theErr != nil {
		handleErrors(h.templates, err, w, r)
	}
}

// CheckEmailHandler ...
func (h *Handler) CheckEmailHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	email := params.Get("email")
	customer, err := h.uc.GetByEmail(email)
	if err != nil {
		handleErrors(h.templates, err, w, r)
	}
	// Return the customer so FE can check if the TAKEN email belongs to this customer or not. Necessary for EDIT screen.
	js, err := json.Marshal(customer)
	if err != nil {
		handleErrors(h.templates, err, w, r)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
