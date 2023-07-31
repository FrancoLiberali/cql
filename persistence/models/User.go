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

func UserEmailCondition(email string) orm.Condition[User] {
	return orm.WhereCondition[User]{
		Field: "email",
		Value: email,
	}
}
