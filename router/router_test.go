package router

import (
	"net/http"
	"testing"

	controllersMocks "github.com/ditrit/badaas/mocks/controllers"
	middlewaresMocks "github.com/ditrit/badaas/mocks/router/middlewares"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSetupRouter(t *testing.T) {
	jsonController := middlewaresMocks.NewJSONController(t)
	middlewareLogger := middlewaresMocks.NewMiddlewareLogger(t)
	informationController := controllersMocks.NewInformationController(t)
	jsonController.On("Wrap", mock.Anything).Return(func(response http.ResponseWriter, request *http.Request) {})
	router := SetupRouter(jsonController, middlewareLogger, informationController)
	assert.NotNil(t, router)
}
