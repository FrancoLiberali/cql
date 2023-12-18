package models

import (
	"time"

	"github.com/ditrit/badaas/orm/model"
)

// Represent a user session
type Session struct {
	model.UUIDModel
	UserID    model.UUID `gorm:"not null"`
	ExpiresAt time.Time  `gorm:"not null"`
}

// Create a new session
func NewSession(userID model.UUID, sessionDuration time.Duration) *Session {
	return &Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(sessionDuration),
	}
}

// Return true is expired
func (session *Session) IsExpired() bool {
	return time.Now().After(session.ExpiresAt)
}

// Return true if the session is expired in less than an hour
func (session *Session) CanBeRolled(rollInterval time.Duration) bool {
	return time.Now().After(session.ExpiresAt.Add(-rollInterval))
}
