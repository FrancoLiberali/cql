package sqlite

import (
	"github.com/ditrit/badaas/orm/operator"
	"github.com/ditrit/badaas/orm/sql"
)

// ref: https://www.sqlie.org/lang_expr.html#like
func Glob(pattern string) operator.Operator[string] {
	return operator.NewValueOperator[string](sql.SQLiteGlob, pattern)
}
