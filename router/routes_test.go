package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/ditrit/badaas/controllers"
	mockControllers "github.com/ditrit/badaas/mocks/controllers"
	mockMiddlewares "github.com/ditrit/badaas/mocks/router/middlewares"
	"github.com/ditrit/badaas/router/middlewares"
)

var logger, _ = zap.NewDevelopment()

func TestAddInfoRoutes(t *testing.T) {
	jsonController := middlewares.NewJSONController(logger)
	informationController := controllers.NewInfoController(semver.MustParse("1.0.1"))

	router := NewRouter()
	AddInfoRoutes(
		router,
		jsonController,
		informationController,
	)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodGet,
		"/info",
		nil,
	)

	router.ServeHTTP(response, request)
	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, response.Body.String(), "{\"status\":\"OK\",\"version\":\"1.0.1\"}")
}

func TestAddLoginRoutes(t *testing.T) {
	jsonController := middlewares.NewJSONController(logger)

	basicAuthenticationController := mockControllers.NewBasicAuthenticationController(t)
	basicAuthenticationController.
		On("BasicLoginHandler", mock.Anything, mock.Anything).
		Return(map[string]string{"login": "called"}, nil)

	authenticationMiddleware := mockMiddlewares.NewAuthenticationMiddleware(t)

	router := NewRouter()
	AddAuthRoutes(
		router,
		authenticationMiddleware,
		basicAuthenticationController,
		jsonController,
	)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/login",
		nil,
	)

	router.ServeHTTP(response, request)
	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, response.Body.String(), "{\"login\":\"called\"}")
}
