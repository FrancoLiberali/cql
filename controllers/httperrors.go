package controllers

import (
	"net/http"

	"github.com/ditrit/badaas/services/httperrors"
)

var (
	// Sent when the request is malformed
	HTTPErrRequestMalformed httperrors.HTTPError = httperrors.NewHTTPError(
		http.StatusBadRequest,
		"Request malformed",
		"The schema of the received data is not correct",
		nil,
		false)
)
