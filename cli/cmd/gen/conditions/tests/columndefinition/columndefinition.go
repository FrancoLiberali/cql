package columndefinition

import "github.com/FrancoLiberali/cql/orm/model"

type ColumnDefinition struct {
	model.UUIDModel

	String string `gorm:"column:string_something_else"`
}
