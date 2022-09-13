package httperrors

import (
	"fmt"
	"net/http"
)

// Describe an HTTP error
type HTTPError struct {
	Status      int
	Err         string
	Message     string
	GolangError error
}

// HTTPError constructor
func NewHTTPError(status int, err string, message string, golangError error) *HTTPError {
	return &HTTPError{
		Status:      status,
		Err:         err,
		Message:     message,
		GolangError: golangError,
	}
}

// Convert an HTTPError to a json string
func (httpError *HTTPError) ToJSON() string {
	return fmt.Sprintf(`{"error": %q, "msg":%q, "status": %q}`, httpError.Err, httpError.Message, http.StatusText(httpError.Status))
}

// Write the HTTPError to the [http.ResponseWriter] passed as argument.
func (httpError *HTTPError) Write(httpResponse http.ResponseWriter) {
	http.Error(httpResponse, httpError.ToJSON(), httpError.Status)
}
