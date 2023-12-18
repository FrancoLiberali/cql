package controllers

import (
	"net/http"

	"github.com/Masterminds/semver/v3"

	"github.com/ditrit/badaas/httperrors"
)

// The information controller
type InformationController interface {
	// Return the badaas server information
	Info(response http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError)
}

// check interface compliance
var _ InformationController = (*infoControllerImpl)(nil)

// The concrete implementation of the InformationController
type infoControllerImpl struct {
	Version *semver.Version
}

// The InformationController constructor
func NewInfoController(version *semver.Version) InformationController {
	return &infoControllerImpl{
		Version: version,
	}
}

// Return the badaas server information
func (c *infoControllerImpl) Info(_ http.ResponseWriter, _ *http.Request) (any, httperrors.HTTPError) {
	return &BadaasServerInfo{
		Status:  "OK",
		Version: c.Version.String(),
	}, nil
}

type BadaasServerInfo struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}
