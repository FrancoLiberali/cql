package sessionservice

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSessionCtx(t *testing.T) {
	ctx := context.Background()
	sessionClaims := &SessionClaims{uint(12), uuid.New()}
	ctx = SetSessionClaimsContext(ctx, sessionClaims)
	claims := GetSessionClaimsFromContext(ctx)
	assert.Equal(t, uint(12), claims.UserID)
}

func TestSessionCtxPanic(t *testing.T) {
	ctx := context.Background()
	assert.Panics(t, func() { GetSessionClaimsFromContext(ctx) })
}
