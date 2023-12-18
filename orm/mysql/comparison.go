package mysql

import (
	"github.com/ditrit/badaas/orm/operator"
	"github.com/ditrit/badaas/orm/sql"
)

// Pattern Matching

// As an extension to standard SQL, MySQL permits LIKE on numeric expressions.
func Like[T string |
	int | int8 | int16 | int32 | int64 |
	uint | uint8 | uint16 | uint32 | uint64 |
	float32 | float64](pattern string,
) operator.Operator[T] {
	return operator.NewValueOperator[T](sql.Like, pattern)
}

// ref: https://dev.mysql.com/doc/refman/8.0/en/regexp.html#operator_regexp
func RegexP(pattern string) operator.Operator[string] {
	return operator.NewValueOperator[string](sql.MySQLRegexp, pattern)
}
