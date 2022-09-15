package httperrors_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ditrit/badaas/services/httperrors"
)

func TestNewHTTPError(t *testing.T) {
	error := httperrors.NewHTTPError(http.StatusBadRequest, "Error while parsing json", "The request body was malformed", nil)
	if error == nil {
		t.Errorf("the HTTPError returned by the contructor should not be nil")
	}
}

func TestTojson(t *testing.T) {
	error := httperrors.NewHTTPError(http.StatusBadRequest, "Error while parsing json", "The request body was malformed", nil)
	if error.ToJSON() == "" {
		t.Errorf("the json string returned by the ToJSON method should not return an empty string")
	}
	if !json.Valid([]byte(error.ToJSON())) {
		t.Errorf("the method ToJSON should return a valid json string")
	}
}

func TestWrite(t *testing.T) {
	res := httptest.NewRecorder()
	error := httperrors.NewHTTPError(http.StatusBadRequest, "Error while parsing json", "The request body was malformed", nil)
	error.Write(res)
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("the body of the response have to be readable without errors, got =%s", err.Error())
	}
	if len(bodyBytes) == 0 {
		t.Error("the body of the response shouldn't be nul")
	}
	originalBytes := []byte(error.ToJSON())
	if !bytes.Contains(bodyBytes, originalBytes) {
		t.Error("the body should contains the jsonified HTTPError")
	}
}
