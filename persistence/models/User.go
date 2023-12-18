package models

import (
	"github.com/ditrit/badaas/orm/condition"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/operator"
	"github.com/ditrit/badaas/orm/query"
)

// Represents a user
type User struct {
	model.UUIDModel
	Username string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`

	// password hash
	Password []byte `gorm:"not null"`
}

func UserEmailCondition(operator operator.Operator[string]) condition.WhereCondition[User] {
	return condition.NewFieldCondition[User, string](
		query.FieldIdentifier[string]{
			Field: "Email",
		},
		operator,
	)
}
