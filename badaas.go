//Package main :
package main

import (
	"github.com/ditrit/badaas/router"
	"log"
	"net/http"
	"time"
)

// Badaas application, run a http-server on 8000.
func main() {
	router := router.SetupRouter()

	srv := &http.Server{
		Handler: router,
		Addr:    ":8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
