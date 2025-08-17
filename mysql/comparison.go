package mysql

import (
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/sql"
)

// Pattern Matching

// As an extension to standard SQL, MySQL permits LIKE on numeric expressions.
func Like[T string | condition.Numeric](pattern string,
) condition.Operator[T] {
	return condition.NewValueOperator[T](sql.Like, condition.String(pattern))
}

// ref: https://dev.mysql.com/doc/refman/8.0/en/regexp.html#operator_regexp
func Regexp(pattern string) condition.Operator[string] {
	return condition.NewValueOperator[string](sql.MySQLRegexp, condition.String(pattern))
}
