package psql

import (
	"github.com/FrancoLiberali/cql/orm/cql"
	"github.com/FrancoLiberali/cql/orm/sql"
)

// Pattern Matching

func ILike(pattern string) cql.Operator[string] {
	return cql.NewValueOperator[string](sql.PostgreSQLILike, pattern)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-SIMILARTO-REGEXP
func SimilarTo(pattern string) cql.Operator[string] {
	return cql.NewValueOperator[string](sql.PostgreSQLSimilarTo, pattern)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-POSIX-REGEXP
func POSIXMatch(pattern string) cql.Operator[string] {
	return cql.NewValueOperator[string](sql.PostgreSQLPosixMatch, pattern)
}

// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-POSIX-REGEXP
func POSIXIMatch(pattern string) cql.Operator[string] {
	return cql.NewValueOperator[string](sql.PostgreSQLPosixIMatch, pattern)
}
