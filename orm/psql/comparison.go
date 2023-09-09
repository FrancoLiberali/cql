package psql

import (
	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/sql"
)

// Pattern Matching

func ILike(pattern string) orm.Operator[string] {
	return orm.NewValueOperator[string](sql.PostgreSQLILike, pattern)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-SIMILARTO-REGEXP
func SimilarTo(pattern string) orm.Operator[string] {
	return orm.NewValueOperator[string](sql.PostgreSQLSimilarTo, pattern)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-POSIX-REGEXP
func POSIXMatch(pattern string) orm.Operator[string] {
	return orm.NewValueOperator[string](sql.PostgreSQLPosixMatch, pattern)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-POSIX-REGEXP
func POSIXIMatch(pattern string) orm.Operator[string] {
	return orm.NewValueOperator[string](sql.PostgreSQLPosixIMatch, pattern)
}
