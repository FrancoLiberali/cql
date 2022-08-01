package router

import (
	"github.com/ditrit/badaas/controllers"
	"github.com/gorilla/mux"
)

// Default router of badaas, initialize all routes.
func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/info", controllers.Info).Methods("GET")

	return router
}
