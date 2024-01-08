package cql

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/model"
	"github.com/elliotchance/pie/v2"
)

// Create a Delete to which the conditions are applied inside transaction tx.
//
// At least one condition is required to avoid deleting all values in a table.
// In case this is the desired behavior, use cql.True.
//
// For details see https://compiledquerylenguage.readthedocs.io/en/latest/cql/delete.html
func Delete[T model.Model](tx *gorm.DB, firstCondition condition.Condition[T], conditions ...condition.Condition[T]) *condition.Delete[T] {
	return condition.NewDelete(
		tx,
		pie.Unshift(conditions, firstCondition),
	)
}
