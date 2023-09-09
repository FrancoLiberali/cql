package sqlite

import (
	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/sql"
)

// ref: https://www.sqlie.org/lang_expr.html#like
func Glob(pattern string) orm.Operator[string] {
	return orm.NewValueOperator[string](sql.SQLiteGlob, pattern)
}
