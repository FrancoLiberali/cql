package orm

import (
	"github.com/ditrit/badaas/orm/condition"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/sql"
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
	return condition.NewContainerCondition(sql.Not, conditions...)
}
