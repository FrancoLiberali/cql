package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/ditrit/badaas/services/httperrors"
)

var HErrorCantMarshall httperrors.HTTPError = httperrors.NewInternalServerError(
	"json marshall error",
	"Can't marshall the object returned by the JSON handler",
	nil,
)

// Marshall the response from the JSONHandler and handle HTTPError if needed
func JSONController(handler JSONHandler) func(response http.ResponseWriter, request *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		object, herr := handler(response, request)
		if herr != nil {
			herr.Write(response)
			return
		}
		if object == nil {
			return
		}
		payload, err := json.Marshal(object)
		if err != nil {
			HErrorCantMarshall.Write(response)
			return
		}
		response.Header().Set("Content-Type", "application/json")
		response.Write(payload)
	}
}
