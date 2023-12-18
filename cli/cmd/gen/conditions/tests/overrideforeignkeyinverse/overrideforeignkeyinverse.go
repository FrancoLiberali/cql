package overrideforeignkeyinverse

import (
	"github.com/ditrit/badaas/orm/model"
)

type User struct {
	model.UUIDModel
	CreditCard CreditCard `gorm:"foreignKey:UserReference"`
}

type CreditCard struct {
	model.UUIDModel
	UserReference model.UUID
}
