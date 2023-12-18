package overrideforeignkeyinverse

import (
	"github.com/FrancoLiberali/cql/orm/model"
)

type User struct {
	model.UUIDModel
	CreditCard CreditCard `gorm:"foreignKey:UserReference"`
}

type CreditCard struct {
	model.UUIDModel
	UserReference model.UUID
}
