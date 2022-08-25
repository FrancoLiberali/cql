package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/resources"
)

// Info controller, return json with status and version of api.
func Info(response http.ResponseWriter, _ *http.Request) {

	infos := models.BadaasServerInfo{
		Status:  "OK",
		Version: resources.Version,
	}

	payload, err := json.Marshal(&infos)
	if err != nil {
		http.Error(response, "error while marshaling response", http.StatusInternalServerError)
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(payload)
}
