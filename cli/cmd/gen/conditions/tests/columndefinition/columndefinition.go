package columndefinition

import "github.com/ditrit/badaas/orm/model"

type ColumnDefinition struct {
	model.UUIDModel

	String string `gorm:"column:string_something_else"`
}
