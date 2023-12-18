package cql

import (
	"github.com/FrancoLiberali/cql/orm/model"
	"github.com/FrancoLiberali/cql/orm/sql"
)

func And[T model.Model](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewConnectionCondition(sql.And, conditions...)
}
