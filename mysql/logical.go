package mysql

import (
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

func Xor[T model.Model](conditions ...condition.WhereCondition[T]) condition.WhereCondition[T] {
	return condition.NewConnectionCondition(sql.MySQLXor, conditions...)
}
