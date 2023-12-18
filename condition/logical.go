package condition

import (
	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

func And[T model.Model](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewConnectionCondition(sql.And, conditions...)
}
