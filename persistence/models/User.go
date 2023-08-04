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

func UserEmailCondition(operator orm.Operator[string]) orm.WhereCondition[User] {
	return orm.FieldCondition[User, string]{
		FieldIdentifier: orm.FieldIdentifier[string]{
			Field: "Email",
		},
		Operator: operator,
	}
}
