package middlewares

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/services/sessionservice"
)

var NotAuthenticated = httperrors.NewUnauthorizedError("Authentication Error", "not authenticated")

// The authentication middleware
type AuthenticationMiddleware interface {
	Handle(next http.Handler) http.Handler
}

// Check interface compliance
var _ AuthenticationMiddleware = (*authenticationMiddleware)(nil)

// The AuthenticationMiddleware implementation
type authenticationMiddleware struct {
	sessionService sessionservice.SessionService
	logger         *zap.Logger
}

// The AuthenticationMiddleware constructor
func NewAuthenticationMiddleware(sessionService sessionservice.SessionService, logger *zap.Logger) AuthenticationMiddleware {
	return &authenticationMiddleware{
		sessionService: sessionService,
		logger:         logger,
	}
}

// The authentication middleware
func (authenticationMiddleware *authenticationMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		accessTokenCookie, err := request.Cookie("access_token")
		if err != nil {
			NotAuthenticated.Write(response, authenticationMiddleware.logger)
			return
		}

		extractedUUID, err := orm.ParseUUID(accessTokenCookie.Value)
		if err != nil {
			NotAuthenticated.Write(response, authenticationMiddleware.logger)
			return
		}
		ok, sessionClaims := authenticationMiddleware.sessionService.IsValid(extractedUUID)
		if !ok {
			NotAuthenticated.Write(response, authenticationMiddleware.logger)
			return
		}
		herr := authenticationMiddleware.sessionService.RollSession(extractedUUID)
		if herr != nil {
			herr.Write(response, authenticationMiddleware.logger)
			return
		}
		request = request.WithContext(sessionservice.SetSessionClaimsContext(
			request.Context(), sessionClaims))
		next.ServeHTTP(response, request)
	})
}
