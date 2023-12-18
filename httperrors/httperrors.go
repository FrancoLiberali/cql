package httperrors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/ditrit/badaas/persistence/models/dto"
)

// ErrForTests is an HTTPError instance useful for testing. If the code does not care
// about HTTPError specifics, and only needs to return the HTTPError for example, this
// HTTPError should be used to make the test code more readable.
var ErrForTests HTTPError = &HTTPErrorImpl{
	Status:      -1,
	Err:         "TESTING ERROR",
	Message:     "USE ONLY FOR TESTING",
	GolangError: nil,
	toLog:       true,
}

type HTTPError interface {
	error

	// Convert the Error to a valid json string
	ToJSON() string

	// Write the error to the http response
	Write(httpResponse http.ResponseWriter, logger *zap.Logger)

	// do we log the error
	Log() bool
}

// Describe an HTTP error
//
//nolint:errname // this name is correct for a type name
type HTTPErrorImpl struct {
	Status      int
	Err         string
	Message     string
	GolangError error
	toLog       bool
}

// Convert an HTTPError to a json string
func (httpError *HTTPErrorImpl) ToJSON() string {
	dto := &dto.HTTPError{
		Error:   httpError.Err,
		Message: httpError.Message,
		Status:  http.StatusText(httpError.Status),
	}

	//nolint:errchkjson // TODO fix it
	payload, _ := json.Marshal(dto)

	return string(payload)
}

// Implement the Error interface
func (httpError *HTTPErrorImpl) Error() string {
	return fmt.Sprintf(`HTTPError: %s`, httpError.ToJSON())
}

// Return true is the error is logged
func (httpError *HTTPErrorImpl) Log() bool {
	return httpError.toLog
}

// Write the HTTPError to the [http.ResponseWriter] passed as argument.
func (httpError *HTTPErrorImpl) Write(httpResponse http.ResponseWriter, logger *zap.Logger) {
	if httpError.toLog && logger != nil {
		logHTTPError(httpError, logger)
	}

	http.Error(httpResponse, httpError.ToJSON(), httpError.Status)
}

func logHTTPError(httpError *HTTPErrorImpl, logger *zap.Logger) {
	logger.Info(
		"http error",
		zap.String("error", httpError.Err),
		zap.String("msg", httpError.Message),
		zap.Int("status", httpError.Status),
	)
}

// HTTPError constructor
func NewHTTPError(status int, err string, message string, golangError error, toLog bool) HTTPError {
	return &HTTPErrorImpl{
		Status:      status,
		Err:         err,
		Message:     message,
		GolangError: golangError,
		toLog:       toLog,
	}
}

// A constructor for an HttpError "Not Found"
func NewErrorNotFound(resourceName string, msg string) HTTPError {
	return NewHTTPError(
		http.StatusNotFound,
		fmt.Sprintf("%s not found", resourceName),
		msg,
		nil,
		false,
	)
}

// A constructor for an HttpError "Internal Server Error"
func NewInternalServerError(errorName string, msg string, err error) HTTPError {
	return NewHTTPError(
		http.StatusInternalServerError,
		errorName,
		msg,
		err,
		true,
	)
}

// Constructor for an HttpError "DB Error", a internal server error produced by a query
func NewDBError(err error) HTTPError {
	return NewInternalServerError("db error", "database query failed", err)
}

// A constructor for an HttpError "Unauthorized Error"
func NewUnauthorizedError(errorName string, msg string) HTTPError {
	return NewHTTPError(
		http.StatusUnauthorized,
		errorName,
		msg,
		nil,
		true,
	)
}
