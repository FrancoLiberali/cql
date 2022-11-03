package httperrors_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestTojson(t *testing.T) {
	err := "Error while parsing json"
	message := "The request body was malformed"
	error := httperrors.NewHTTPError(http.StatusBadRequest, err, message, nil, true)
	assert.NotEmpty(t, error.ToJSON())
	assert.True(t, json.Valid([]byte(error.ToJSON())), "output json is not valid")

	// check if is is correctly deserialized
	var content map[string]any
	json.Unmarshal([]byte(error.ToJSON()), &content)
	_, ok := content["err"]
	assert.True(t, ok, "\"err\" field should be in the json string")
	_, ok = content["msg"]
	assert.True(t, ok, "\"msg\" field should be in the json string")
	_, ok = content["status"]
	assert.True(t, ok, "\"status\" field should be in the json string")

	assert.Equal(t, err, content["err"].(string))
	assert.Equal(t, message, content["msg"].(string))
	assert.Equal(t, http.StatusText(http.StatusBadRequest), content["status"].(string))
	assert.True(t, error.Log())
}

func TestLog(t *testing.T) {
	error := httperrors.NewHTTPError(http.StatusBadRequest, "err", "message", nil, true)
	assert.True(t, error.Log())
	error = httperrors.NewHTTPError(http.StatusBadRequest, "err", "message", nil, false)
	assert.False(t, error.Log())
}

func TestError(t *testing.T) {
	error := httperrors.NewHTTPError(http.StatusBadRequest, "Error while parsing json", "The request body was malformed", nil, true)
	assert.Contains(t, error.Error(), error.ToJSON())
}

func TestWrite(t *testing.T) {
	res := httptest.NewRecorder()
	error := httperrors.NewHTTPError(http.StatusBadRequest, "Error while parsing json", "The request body was malformed", nil, true)
	error.Write(res, zap.L())
	bodyBytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.NotEmpty(t, bodyBytes)
	originalBytes := []byte(error.ToJSON())
	// can't use assert.Contains because it only support strings
	assert.True(t,
		bytes.Contains(bodyBytes, originalBytes))
}

func TestLogger(t *testing.T) {
	// creating logger
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	res := httptest.NewRecorder()
	error := httperrors.NewHTTPError(http.StatusBadRequest, "Error while parsing json", "The request body was malformed", nil, true)
	error.Write(res, observedLogger)

	require.Equal(t, 1, observedLogs.Len())
	log := observedLogs.All()[0]
	assert.Equal(t, "http error", log.Message)
	require.Len(t, log.Context, 3)
	assert.ElementsMatch(t, []zap.Field{
		{Key: "error", Type: zapcore.StringType, String: "Error while parsing json"},
		{Key: "msg", Type: zapcore.StringType, String: "The request body was malformed"},
		{Key: "status", Type: zapcore.Int64Type, Integer: http.StatusBadRequest},
	}, log.Context)
}

func TestNewErrorNotFound(t *testing.T) {
	ressourceName := "file"
	error := httperrors.NewErrorNotFound(ressourceName, "main.css is not accessible")
	assert.NotNil(t, error)
	assert.False(t, error.Log())
	dto := new(dto.DTOHTTPError)
	err := json.Unmarshal([]byte(error.ToJSON()), &dto)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusText(http.StatusNotFound), dto.Status)
	assert.Equal(t, fmt.Sprintf("%s not found", ressourceName), dto.Error)
}

func TestNewInternalServerError(t *testing.T) {
	error := httperrors.NewInternalServerError("casbin error", "the ressource is not accessible", nil)
	assert.NotNil(t, error)
	assert.True(t, error.Log())
	dto := new(dto.DTOHTTPError)
	err := json.Unmarshal([]byte(error.ToJSON()), &dto)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError), dto.Status)
}

func TestNewUnauthorizedError(t *testing.T) {
	error := httperrors.NewUnauthorizedError("json unmarshalling", "nil value whatever")
	assert.NotNil(t, error)
	assert.True(t, error.Log())
	dto := new(dto.DTOHTTPError)
	err := json.Unmarshal([]byte(error.ToJSON()), &dto)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusText(http.StatusUnauthorized), dto.Status)
}
