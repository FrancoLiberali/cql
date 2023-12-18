package middlewares

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/ditrit/badaas/httperrors"
)

// transform a JSON handler into a standard [http.HandlerFunc]
// handle [github.com/ditrit/badaas/httperrors.HTTPError] and JSON marshaling
type JSONController interface {
	// Marshall the response from the JSONHandler and handle HTTPError if needed
	Wrap(handler JSONHandler) func(response http.ResponseWriter, request *http.Request)
}

// check interface compliance
var _ JSONController = (*jsonControllerImpl)(nil)

// The concrete implementation of JsonController
type jsonControllerImpl struct {
	logger *zap.Logger
}

func NewJSONController(logger *zap.Logger) JSONController {
	return &jsonControllerImpl{logger}
}

// Transforms a JSONHandler into a standard [http.HandlerFunc]
// It marshalls the response from the JSONHandler and handles HTTPError if needed
func (controller *jsonControllerImpl) Wrap(handler JSONHandler) func(response http.ResponseWriter, request *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		object, herr := handler(response, request)
		if herr != nil {
			herr.Write(response, controller.logger)
			return
		}

		if object == nil {
			return
		}

		payload, err := json.Marshal(object)
		if err != nil {
			httperrors.NewInternalServerError(
				"json marshall error",
				"Can't marshall the object returned by the JSON handler",
				nil,
			).Write(response, controller.logger)

			return
		}

		response.Header().Set("Content-Type", "application/json")

		_, err = response.Write(payload)
		if err != nil {
			controller.logger.Error(
				"Error while writing http response",
				zap.String("error", err.Error()),
			)
		}
	}
}
