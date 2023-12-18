package mysql

import (
	"github.com/FrancoLiberali/cql/orm/cql"
	"github.com/FrancoLiberali/cql/orm/model"
	"github.com/FrancoLiberali/cql/orm/sql"
)

func Xor[T model.Model](conditions ...cql.WhereCondition[T]) cql.WhereCondition[T] {
	return cql.NewConnectionCondition(sql.MySQLXor, conditions...)
}
