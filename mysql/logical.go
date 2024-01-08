package mysql

import (
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
	"github.com/elliotchance/pie/v2"
)

func Xor[T model.Model](firstCondition condition.WhereCondition[T], conditions ...condition.WhereCondition[T]) condition.WhereCondition[T] {
	return condition.NewConnectionCondition(
		sql.MySQLXor,
		pie.Unshift(conditions, firstCondition),
	)
}
