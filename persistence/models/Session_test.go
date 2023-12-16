package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/persistence/models"
)

func TestNewSession(t *testing.T) {
	sessionInstance := models.NewSession(orm.NilUUID, time.Second)
	assert.NotNil(t, sessionInstance)
	assert.Equal(t, orm.NilUUID, sessionInstance.UserID)
}

func TestExpired(t *testing.T) {
	sessionInstance := &models.Session{
		ExpiresAt: time.Now().Add(time.Second),
	}
	assert.False(t, sessionInstance.IsExpired())
	sessionInstance.ExpiresAt = time.Now().Add(-5 * time.Second)
	assert.True(t, sessionInstance.IsExpired())
}

func TestCanBeRolled(t *testing.T) {
	sessionDuration := 500 * time.Millisecond
	sessionInstance := &models.Session{
		ExpiresAt: time.Now().Add(sessionDuration),
	}
	assert.False(t, sessionInstance.CanBeRolled(sessionDuration/4))
	time.Sleep(400 * time.Millisecond)
	assert.True(t, sessionInstance.CanBeRolled(sessionDuration))
}
