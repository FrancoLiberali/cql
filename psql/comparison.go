package psql

import (
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/sql"
)

// Pattern Matching

func ILike(pattern string) condition.Operator[string] {
	return condition.NewValueOperator[string](sql.PostgreSQLILike, pattern)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-SIMILARTO-REGEXP
func SimilarTo(pattern string) condition.Operator[string] {
	return condition.NewValueOperator[string](sql.PostgreSQLSimilarTo, pattern)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-POSIX-REGEXP
func POSIXMatch(pattern string) condition.Operator[string] {
	return condition.NewValueOperator[string](sql.PostgreSQLPosixMatch, pattern)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-POSIX-REGEXP
func POSIXIMatch(pattern string) condition.Operator[string] {
	return condition.NewValueOperator[string](sql.PostgreSQLPosixIMatch, pattern)
}
