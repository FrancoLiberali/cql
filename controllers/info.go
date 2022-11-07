package controllers

import (
	"net/http"

	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models/dto"
	"github.com/ditrit/badaas/resources"
)

// The information controller
type InformationController interface {
	// Return the badaas server informations
	Info(response http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
}

// check interface compliance
var _ InformationController = (*infoControllerImpl)(nil)

// The InformationController constructor
func NewInfoController() InformationController {
	return &infoControllerImpl{}
}

// The concrete implementation of the InformationController
type infoControllerImpl struct{}

// Return the badaas server informations
func (*infoControllerImpl) Info(response http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {

	infos := &dto.DTOBadaasServerInfo{
		Status:  "OK",
		Version: resources.Version,
	}
	return infos, nil
}
