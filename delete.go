package cql

import (
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/model"
)

// Create a Delete to which the conditions are applied inside transaction tx.
//
// At least one condition is required to avoid deleting all values in a table.
// In case this is the desired behavior, use cql.True.
//
// For details see https://compiledquerylenguage.readthedocs.io/en/latest/cql/delete.html
func Delete[T model.Model](tx *DB, conditions ...condition.Condition[T]) *condition.Delete[T] {
	return condition.NewDelete(
		tx.GormDB,
		conditions,
	)
}
