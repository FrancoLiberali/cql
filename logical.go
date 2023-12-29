package cql

import (
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

// Logical Operators
// ref: https://www.postgresql.org/docs/current/functions-logical.html

func And[T model.Model](conditions ...condition.WhereCondition[T]) condition.WhereCondition[T] {
	return condition.And(conditions...)
}

func Or[T model.Model](conditions ...condition.WhereCondition[T]) condition.WhereCondition[T] {
	return condition.NewConnectionCondition(sql.Or, conditions...)
}

func Not[T model.Model](conditions ...condition.WhereCondition[T]) condition.WhereCondition[T] {
	return condition.Not(conditions...)
}
