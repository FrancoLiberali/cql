package cql

import (
	"context"

	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/model"
)

// Create a Update to which the conditions are applied inside transaction tx.
//
// At least one condition is required to avoid updating all values in a table.
// In case this is the desired behavior, use cql.True.
//
// For details see https://compiledquerylenguage.readthedocs.io/en/latest/cql/update.html
func Update[T model.Model](ctx context.Context, tx *DB, conditions ...condition.Condition[T]) *condition.Update[T] {
	return condition.NewUpdate(
		tx.GormDB.WithContext(ctx),
		conditions,
	)
}
