package errors

// CustomerNotFoundError - a custom error to wrap the GORM 404
type CustomerNotFoundError struct{}

func (m *CustomerNotFoundError) Error() string {
	return "Customer Not Found"
}

// EmailNotFoundError - a custom error to notify the FE that email is available!
type EmailNotFoundError struct{}

func (m *EmailNotFoundError) Error() string {
	return "Email Not Found"
}
