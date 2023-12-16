package sessionservice

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/ditrit/badaas/orm"
)

func TestSessionCtx(t *testing.T) {
	ctx := context.Background()
	sessionClaims := &SessionClaims{orm.NilUUID, orm.UUID(uuid.New())}
	ctx = SetSessionClaimsContext(ctx, sessionClaims)
	claims := GetSessionClaimsFromContext(ctx)
	assert.Equal(t, orm.NilUUID, claims.UserID)
}

func TestSessionCtxPanic(t *testing.T) {
	ctx := context.Background()
	assert.Panics(t, func() { GetSessionClaimsFromContext(ctx) })
}
