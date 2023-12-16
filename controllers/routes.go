package controllers

import (
	"github.com/gorilla/mux"

	"github.com/ditrit/badaas/router/middlewares"
)

func AddInfoRoutes(
	router *mux.Router,
	jsonController middlewares.JSONController,
	infoController InformationController,
) {
	router.HandleFunc(
		"/info",
		jsonController.Wrap(infoController.Info),
	).Methods("GET")
}

// Adds to the "router" the routes for handling authentication:
// /login
// /logout
// And creates a very first user
func AddAuthRoutes(
	router *mux.Router,
	authenticationMiddleware middlewares.AuthenticationMiddleware,
	basicAuthenticationController BasicAuthenticationController,
	jsonController middlewares.JSONController,
) {
	router.HandleFunc(
		"/login",
		jsonController.Wrap(basicAuthenticationController.BasicLoginHandler),
	).Methods("POST")

	protected := router.PathPrefix("").Subrouter()
	protected.Use(authenticationMiddleware.Handle)

	protected.HandleFunc(
		"/logout",
		jsonController.Wrap(basicAuthenticationController.Logout),
	).Methods("GET")
}
