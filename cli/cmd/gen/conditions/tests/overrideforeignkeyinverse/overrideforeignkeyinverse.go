package overrideforeignkeyinverse

import "github.com/ditrit/badaas/orm"

type User struct {
	orm.UUIDModel
	CreditCard CreditCard `gorm:"foreignKey:UserReference"`
}

type CreditCard struct {
	orm.UUIDModel
	UserReference orm.UUID
}
