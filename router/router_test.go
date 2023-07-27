package router

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	controllersMocks "github.com/ditrit/badaas/mocks/controllers"
	middlewaresMocks "github.com/ditrit/badaas/mocks/router/middlewares"
)

func TestSetupRouter(t *testing.T) {
	jsonController := middlewaresMocks.NewJSONController(t)
	middlewareLogger := middlewaresMocks.NewMiddlewareLogger(t)
	authenticationMiddleware := middlewaresMocks.NewAuthenticationMiddleware(t)

	basicController := controllersMocks.NewBasicAuthenticationController(t)
	informationController := controllersMocks.NewInformationController(t)
	jsonController.On("Wrap", mock.Anything).Return(func(response http.ResponseWriter, request *http.Request) {})
	router := SetupRouter(jsonController, middlewareLogger, authenticationMiddleware, basicController, informationController)
	assert.NotNil(t, router)
}
