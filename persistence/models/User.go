package models

import (
	"github.com/ditrit/badaas/orm/model"
)

// Represents a user
type User struct {
	model.UUIDModel
	Username string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`

	// password hash
	Password []byte `gorm:"not null"`
}
