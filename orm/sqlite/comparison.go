package sqlite

import (
	"github.com/ditrit/badaas/orm/cql"
	"github.com/ditrit/badaas/orm/sql"
)

// ref: https://www.sqlie.org/lang_expr.html#like
func Glob(pattern string) cql.Operator[string] {
	return cql.NewValueOperator[string](sql.SQLiteGlob, pattern)
}
