package condition

import (
	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

func And[T model.Model](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewConnectionCondition(sql.And, conditions)
}

func Or[T model.Model](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewConnectionCondition(sql.Or, conditions)
}

func Not[T model.Model](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewContainerCondition(sql.Not, conditions)
}
