package middlewares

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	configurationmocks "github.com/ditrit/badaas/mocks/configuration"
)

func TestMiddlewareLogger(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "http",
			Host:   "localhost",
			Path:   "/whatever",
		},
	}
	res := httptest.NewRecorder()

	actuallyRunned := false
	// create a handler to use as "next" which will verify the request
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actuallyRunned = true
	})
	loggerConfiguration := configurationmocks.NewLoggerConfiguration(t)
	loggerConfiguration.On("GetRequestTemplate").Return("Receive {{method}} request on {{url}}")
	loggerMiddleware, err := NewMiddlewareLogger(observedLogger, loggerConfiguration)
	assert.NoError(t, err)

	loggerMiddleware.Handle(nextHandler).ServeHTTP(res, req)
	assert.True(t, actuallyRunned)
	require.Equal(t, 1, observedLogs.Len())
	assert.Equal(t, "Receive GET request on /whatever", observedLogs.All()[0].Message)
}
