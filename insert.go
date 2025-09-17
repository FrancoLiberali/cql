package cql

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/model"
)

// Insert creates an INSERT statement that will allow to create
// the models received by parameter in the db
// and apply on conflict clauses
func Insert[T model.Model](tx *gorm.DB, models ...*T) *condition.Insert[T] {
	return condition.NewInsert(tx, models)
}
