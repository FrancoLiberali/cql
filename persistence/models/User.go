package models

import "github.com/ditrit/badaas/orm"

// Represents a user
type User struct {
	orm.UUIDModel
	Username string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`

	// password hash
	Password []byte `gorm:"not null"`
}

// Return the pluralized table name
//
// Satisfie the Tabler interface
func (User) TableName() string {
	return "users"
}
