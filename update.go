package cql

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/model"
)

// TODO ver que pasa con el update de las relaciones y de los arreglos de relaciones

// Create a Update to which the conditions are applied inside transaction tx
func Update[T model.Model](tx *gorm.DB, conditions ...condition.Condition[T]) *condition.Update[T] {
	return condition.NewUpdate(tx, conditions...)
}
