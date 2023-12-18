package orm

import (
	"github.com/FrancoLiberali/cql/orm/cql"
	"github.com/FrancoLiberali/cql/orm/model"
	"github.com/FrancoLiberali/cql/orm/sql"
)

// Logical Operators
// ref: https://www.postgresql.org/docs/current/functions-logical.html

func And[T model.Model](conditions ...cql.WhereCondition[T]) cql.WhereCondition[T] {
	return cql.And(conditions...)
}

func Or[T model.Model](conditions ...cql.WhereCondition[T]) cql.WhereCondition[T] {
	return cql.NewConnectionCondition(sql.Or, conditions...)
}

func Not[T model.Model](conditions ...cql.WhereCondition[T]) cql.WhereCondition[T] {
	return cql.NewContainerCondition(sql.Not, conditions...)
}
