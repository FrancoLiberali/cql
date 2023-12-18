package mysql

import (
	"github.com/ditrit/badaas/orm/cql"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/sql"
)

func Xor[T model.Model](conditions ...cql.WhereCondition[T]) cql.WhereCondition[T] {
	return cql.NewConnectionCondition(sql.MySQLXor, conditions...)
}
