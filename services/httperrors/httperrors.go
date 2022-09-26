package httperrors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ditrit/badaas/persistence/models/dto"
	"go.uber.org/zap"
)

type HTTPError interface {
	error

	// Convert the Error to a valid json string
	ToJSON() string

	// Write the error to the http response
	Write(httpResponse http.ResponseWriter)

	// do we log the error
	Log() bool
}

// Describe an HTTP error
type HTTPErrorImpl struct {
	Status      int
	Err         string
	Message     string
	GolangError error
	toLog       bool
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

// Convert an HTTPError to a json string
func (httpError *HTTPErrorImpl) ToJSON() string {
	dto := &dto.DTOHTTPError{
		Error:   httpError.Err,
		Message: httpError.Message,
		Status:  http.StatusText(httpError.Status),
	}
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
func (httpError *HTTPErrorImpl) Write(httpResponse http.ResponseWriter) {
	if httpError.toLog {
		logHTTPError(httpError)
	}
	http.Error(httpResponse, httpError.ToJSON(), httpError.Status)
}

func logHTTPError(httpError *HTTPErrorImpl) {
	zap.L().Info(
		"http error",
		zap.String("error", httpError.Err),
		zap.String("msg", httpError.Message),
		zap.Int("status", httpError.Status),
	)
}

// A contructor for an HttpError "Not Found"
func NewErrorNotFound(ressourceName string, msg string) HTTPError {
	return NewHTTPError(
		http.StatusNotFound,
		fmt.Sprintf("%s not found", ressourceName),
		msg,
		nil,
		false,
	)
}

// A contructor for an HttpError "Internal Server Error"
func NewInternalServerError(errorName string, msg string, err error) HTTPError {
	return NewHTTPError(
		http.StatusInternalServerError,
		errorName,
		msg,
		err,
		true,
	)
}

// A contructor for an HttpError "Unauthorized Error"
func NewUnauthorizedError(errorName string, msg string) HTTPError {
	return NewHTTPError(
		http.StatusUnauthorized,
		errorName,
		msg,
		nil,
		true,
	)
}
