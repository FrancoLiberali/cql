package router

import (
	"github.com/gorilla/mux"
)

// Router to use in Badaas server
func NewRouter() *mux.Router {
	return mux.NewRouter()
}
