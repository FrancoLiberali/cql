package sessionservice

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ditrit/badaas/orm/model"
)

func TestSessionCtx(t *testing.T) {
	ctx := context.Background()
	sessionClaims := &SessionClaims{model.NilUUID, model.UUID(uuid.New())}
	ctx = SetSessionClaimsContext(ctx, sessionClaims)
	claims := GetSessionClaimsFromContext(ctx)
	assert.Equal(t, model.NilUUID, claims.UserID)
}

func TestSessionCtxPanic(t *testing.T) {
	ctx := context.Background()

	assert.Panics(t, func() { GetSessionClaimsFromContext(ctx) })
}
