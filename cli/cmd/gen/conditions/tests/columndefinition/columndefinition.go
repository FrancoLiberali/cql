package columndefinition

import "github.com/ditrit/badaas/orm"

type ColumnDefinition struct {
	orm.UUIDModel

	String string `gorm:"column:string_something_else"`
}
