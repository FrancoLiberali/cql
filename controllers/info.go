package controllers

import (
	"github.com/bitly/go-simplejson"
	"github.com/ditrit/badaas/resources"
	"net/http"
)

// Info controller, return json with status and version of api.
func Info(response http.ResponseWriter, _ *http.Request) {
	json := simplejson.New()
	json.Set("status", "OK")
	json.Set("version", resources.Version)

	response.WriteHeader(http.StatusOK)
	payload, _ := json.MarshalJSON()

	response.WriteHeader(http.StatusOK)
	response.Header().Set("Content-Type", "application/json")
	response.Write(payload)
}
