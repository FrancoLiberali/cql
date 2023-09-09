package orm

import (
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/sql"
)

// Logical Operators
// ref: https://www.postgresql.org/docs/current/functions-logical.html

func And[T model.Model](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewConnectionCondition(sql.And, conditions...)
}

func Or[T model.Model](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewConnectionCondition(sql.Or, conditions...)
}

func Not[T model.Model](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewContainerCondition(sql.Not, conditions...)
}
