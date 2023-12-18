package mysql

import (
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/sql"
)

// Pattern Matching

// As an extension to standard SQL, MySQL permits LIKE on numeric expressions.
func Like[T string |
	int | int8 | int16 | int32 | int64 |
	uint | uint8 | uint16 | uint32 | uint64 |
	float32 | float64](pattern string,
) condition.Operator[T] {
	return condition.NewValueOperator[T](sql.Like, pattern)
}

// ref: https://dev.mysql.com/doc/refman/8.0/en/regexp.html#operator_regexp
func RegexP(pattern string) condition.Operator[string] {
	return condition.NewValueOperator[string](sql.MySQLRegexp, pattern)
}
