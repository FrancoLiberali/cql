package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/persistence/models/dto"
	"github.com/ditrit/badaas/services/sessionservice"
	"github.com/ditrit/badaas/services/userservice"
)

const accessTokenCookieExpirationTime = 48 * time.Hour

// HTTPErrRequestMalformed is sent when the request is malformed
var HTTPErrRequestMalformed httperrors.HTTPError = httperrors.NewHTTPError(
	http.StatusBadRequest,
	"Request malformed",
	"The schema of the received data is not correct",
	nil,
	false,
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

	user, err := basicAuthController.userService.GetUser(loginJSONStruct)
	if err != nil {
		if errors.Is(err, orm.ErrObjectNotFound) {
			return nil, httperrors.NewErrorNotFound(
				"user",
				fmt.Sprintf("no user found with email %q", loginJSONStruct.Email),
			)
		} else if errors.Is(err, userservice.ErrWrongPassword) {
			return nil, httperrors.NewUnauthorizedError(
				"wrong password", "the provided password is incorrect",
			)
		}

		return nil, httperrors.NewDBError(err)
	}

	// On valid password, generate a session and return it's uuid to the client
	session, err := basicAuthController.sessionService.LogUserIn(user)
	if err != nil {
		return nil, httperrors.NewDBError(err)
	}

	herr := createAndSetAccessTokenCookie(w, session.ID.String())
	if herr != nil {
		return nil, herr
	}

	return dto.LoginSuccess{
		Email:    user.Email,
		ID:       user.ID.String(),
		Username: user.Username,
	}, nil
}

// Log Out the user
func (basicAuthController *basicAuthenticationController) Logout(w http.ResponseWriter, r *http.Request) (any, httperrors.HTTPError) {
	herr := basicAuthController.sessionService.LogUserOut(sessionservice.GetSessionClaimsFromContext(r.Context()))
	if herr != nil {
		return nil, herr
	}

	herr = createAndSetAccessTokenCookie(w, "")
	if herr != nil {
		return nil, herr
	}

	return nil, nil
}

func createAndSetAccessTokenCookie(w http.ResponseWriter, sessionUUID string) httperrors.HTTPError {
	accessToken := &http.Cookie{
		Name:     "access_token",
		Path:     "/",
		Value:    sessionUUID,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode, // TODO change to http.SameSiteStrictMode in prod
		Secure:   false,                 // TODO change to true in prod
		Expires:  time.Now().Add(accessTokenCookieExpirationTime),
	}

	err := accessToken.Valid()
	if err != nil {
		return httperrors.NewInternalServerError("access token error", "unable to create access token", err)
	}

	http.SetCookie(w, accessToken)

	return nil
}
