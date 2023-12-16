package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models/dto"
	"github.com/ditrit/badaas/services/sessionservice"
	"github.com/ditrit/badaas/services/userservice"
	"go.uber.org/zap"
)

var (
	// Sent when the request is malformed
	HTTPErrRequestMalformed httperrors.HTTPError = httperrors.NewHTTPError(
		http.StatusBadRequest,
		"Request malformed",
		"The schema of the received data is not correct",
		nil,
		false,
	)
)

type BasicAuthenticationController interface {
	BasicLoginHandler(http.ResponseWriter, *http.Request) (any, httperrors.HTTPError)
	Logout(http.ResponseWriter, *http.Request) (any, httperrors.HTTPError)
}

// Check interface compliance
var _ BasicAuthenticationController = (*basicAuthenticationController)(nil)

type basicAuthenticationController struct {
	logger         *zap.Logger
	userService    userservice.UserService
	sessionService sessionservice.SessionService
}

// BasicAuthenticationController constructor
func NewBasicAuthenticationController(
	logger *zap.Logger,
	userService userservice.UserService,
	sessionService sessionservice.SessionService,
) BasicAuthenticationController {
	return &basicAuthenticationController{
		logger:         logger,
		userService:    userService,
		sessionService: sessionService,
	}
}

// Log In with username and password
func (basicAuthController *basicAuthenticationController) BasicLoginHandler(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	var loginJSONStruct dto.UserLoginDTO
	err := json.NewDecoder(r.Body).Decode(&loginJSONStruct)
	if err != nil {
		return nil, HTTPErrRequestMalformed
	}
	user, herr := basicAuthController.userService.GetUser(loginJSONStruct)
	if herr != nil {
		return nil, herr
	}

	// On valid password, generate a session and return it's uuid to the client
	herr = basicAuthController.sessionService.LogUserIn(user, w)
	if herr != nil {
		return nil, herr

	}

	return dto.DTOLoginSuccess{
		Email:    user.Email,
		ID:       user.ID.String(),
		Username: user.Username,
	}, nil
}

// Log Out the user
func (basicAuthController *basicAuthenticationController) Logout(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	return nil, basicAuthController.sessionService.LogUserOut(sessionservice.GetSessionClaimsFromContext(r.Context()), w)
}
