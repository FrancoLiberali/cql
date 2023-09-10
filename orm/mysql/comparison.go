package mysql

import (
	"github.com/ditrit/badaas/orm/cql"
	"github.com/ditrit/badaas/orm/sql"
)

// Pattern Matching

// As an extension to standard SQL, MySQL permits LIKE on numeric expressions.
func Like[T string |
	int | int8 | int16 | int32 | int64 |
	uint | uint8 | uint16 | uint32 | uint64 |
	float32 | float64](pattern string,
) cql.Operator[T] {
	return cql.NewValueOperator[T](sql.Like, pattern)
}

// ref: https://dev.mysql.com/doc/refman/8.0/en/regexp.html#operator_regexp
func RegexP(pattern string) cql.Operator[string] {
	return cql.NewValueOperator[string](sql.MySQLRegexp, pattern)
}
