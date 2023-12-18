package psql

import (
	"github.com/ditrit/badaas/orm/operator"
	"github.com/ditrit/badaas/orm/sql"
)

// Pattern Matching

func ILike(pattern string) operator.Operator[string] {
	return operator.NewValueOperator[string](sql.PostgreSQLILike, pattern)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-SIMILARTO-REGEXP
func SimilarTo(pattern string) operator.Operator[string] {
	return operator.NewValueOperator[string](sql.PostgreSQLSimilarTo, pattern)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-POSIX-REGEXP
func POSIXMatch(pattern string) operator.Operator[string] {
	return operator.NewValueOperator[string](sql.PostgreSQLPosixMatch, pattern)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-POSIX-REGEXP
func POSIXIMatch(pattern string) operator.Operator[string] {
	return operator.NewValueOperator[string](sql.PostgreSQLPosixIMatch, pattern)
}
