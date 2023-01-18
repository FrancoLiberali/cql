package models

import (
	"time"

	"github.com/google/uuid"
)

// Represent a user session
type Session struct {
	BaseModel
	UserID    uuid.UUID `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

// Return true is expired
func (session *Session) IsExpired() bool {
	return time.Now().After(session.ExpiresAt)
}

// Return true if the session is expired in less than an hour
func (session *Session) CanBeRolled(rollInterval time.Duration) bool {
	return time.Now().After(session.ExpiresAt.Add(-rollInterval))
}

// Return the pluralized table name
//
// Satisfie the Tabler interface
func (Session) TableName() string {
	return "sessions"
}
