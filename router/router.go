package router

import (
	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/router/middlewares"
	"github.com/gorilla/mux"
)

// Default router of badaas, initialize all routes.
func SetupRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(middlewares.CreateLoggerMiddleware())

	router.HandleFunc(
		"/info",
		middlewares.JSONController(controllers.Info),
	).Methods("GET")

	return router
}
