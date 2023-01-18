package sessionservice

import (
	"context"

	"github.com/ditrit/badaas/persistence/models"
	"github.com/google/uuid"
)

// The session claims passed in the request context
type SessionClaims struct {
	UserID      uuid.UUID
	SessionUUID uuid.UUID
}

// Unique claim key type
type sessionClaimsKeyT int

// Unique claim key
var sessionClaimsKey sessionClaimsKeyT

// Set KV pair in request context
func SetSessionClaimsContext(ctx context.Context, sessionClaims *SessionClaims) context.Context {
	return context.WithValue(ctx, sessionClaimsKey,
		sessionClaims,
	)
}

// Create a SessionClaims for a Session
func makeSessionClaims(session *models.Session) *SessionClaims {
	return &SessionClaims{
		UserID:      session.UserID,
		SessionUUID: session.ID,
	}
}

// Extract SessionClaims from request context
// Panics if the claims are not in the context
func GetSessionClaimsFromContext(ctx context.Context) *SessionClaims {
	claims, ok := ctx.Value(sessionClaimsKey).(*SessionClaims)
	if !ok {
		panic("could not extract claims from context")
	}
	return claims
}
