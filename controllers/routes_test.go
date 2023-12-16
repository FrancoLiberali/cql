package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	mockControllers "github.com/ditrit/badaas/mocks/controllers"
	mockMiddlewares "github.com/ditrit/badaas/mocks/router/middlewares"
	"github.com/ditrit/badaas/router"
	"github.com/ditrit/badaas/router/middlewares"
)

var logger, _ = zap.NewDevelopment()

func TestAddInfoRoutes(t *testing.T) {
	jsonController := middlewares.NewJSONController(logger)
	informationController := NewInfoController(semver.MustParse("1.0.1"))

	router := router.NewRouter()
	AddInfoRoutes(
		router,
		jsonController,
		informationController,
	)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"GET",
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

	router := router.NewRouter()
	AddAuthRoutes(
		router,
		authenticationMiddleware,
		basicAuthenticationController,
		jsonController,
	)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(
		"POST",
		"/login",
		nil,
	)

	router.ServeHTTP(response, request)
	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, response.Body.String(), "{\"login\":\"called\"}")
}
