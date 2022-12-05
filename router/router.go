package router

import (
	"net/http"

	"github.com/ditrit/badaas/controllers"
	"github.com/ditrit/badaas/router/middlewares"
	"github.com/gorilla/mux"
)

// Default router of badaas, initialize all routes.
func SetupRouter(
	//middlewares
	jsonController middlewares.JSONController,
	middlewareLogger middlewares.MiddlewareLogger,
	authenticationMiddleware middlewares.AuthenticationMiddleware,

	// controllers
	basicAuthentificationController controllers.BasicAuthentificationController,
	informationController controllers.InformationController,
) http.Handler {
	router := mux.NewRouter()
	router.Use(middlewareLogger.Handle)

	router.HandleFunc(
		"/info",
		jsonController.Wrap(informationController.Info),
	).Methods("GET")
	router.HandleFunc(
		"/login",
		jsonController.Wrap(
			basicAuthentificationController.BasicLoginHandler,
		),
	).Methods("POST")

	protected := router.PathPrefix("").Subrouter()
	protected.Use(authenticationMiddleware.Handle)

	protected.HandleFunc("/logout", jsonController.Wrap(basicAuthentificationController.Logout)).Methods("GET")

	return router
}
