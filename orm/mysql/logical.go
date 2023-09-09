package mysql

import (
	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/sql"
)

func Xor[T model.Model](conditions ...orm.WhereCondition[T]) orm.WhereCondition[T] {
	return orm.NewConnectionCondition(sql.MySQLXor, conditions...)
}
