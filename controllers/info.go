package controllers

import (
	"net/http"

	"github.com/ditrit/badaas/persistence/models/dto"
	"github.com/ditrit/badaas/resources"
	"github.com/ditrit/badaas/services/httperrors"
)

// Return the badaas server informations
func Info(response http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {

	infos := &dto.DTOBadaasServerInfo{
		Status:  "OK",
		Version: resources.Version,
	}
	return infos, nil
}
