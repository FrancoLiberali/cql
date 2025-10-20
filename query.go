package cql

import (
	"context"

	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/model"
)

// Create a Query to which the conditions are applied inside transaction tx
//
// For details see https://compiledquerylenguage.readthedocs.io/en/latest/cql/query.html
func Query[T model.Model](ctx context.Context, tx *DB, conditions ...condition.Condition[T]) *condition.Query[T] {
	return condition.NewQuery(tx.gormDBWithContext(ctx), conditions...)
}
