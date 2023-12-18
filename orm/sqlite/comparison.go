package sqlite

import (
	"github.com/FrancoLiberali/cql/orm/cql"
	"github.com/FrancoLiberali/cql/orm/sql"
)

// ref: https://www.sqlie.org/lang_expr.html#like
func Glob(pattern string) cql.Operator[string] {
	return cql.NewValueOperator[string](sql.SQLiteGlob, pattern)
}
