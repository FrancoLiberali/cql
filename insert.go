package cql

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/model"
)

// TODO docs
func Insert[T model.Model](tx *gorm.DB, models ...*T) *condition.Insert[T] {
	return condition.NewInsert(tx, models)
}
