package sqlite

import (
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/sql"
)

// ref: https://www.sqlie.org/lang_expr.html#like
func Glob(pattern string) condition.Operator[string] {
	return condition.NewValueOperator[string](sql.SQLiteGlob, condition.String(pattern))
}
