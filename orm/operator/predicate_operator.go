package operator

import (
	"fmt"

	"github.com/ditrit/badaas/orm/query"
)

// Operator that verifies a predicate
// Example: value IS TRUE
type PredicateOperator[T any] struct {
	SQLOperator string
}

func (operator PredicateOperator[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T]
}

func (operator PredicateOperator[T]) ToSQL(_ *query.Query, columnName string) (string, []any, error) {
	return fmt.Sprintf("%s %s", columnName, operator.SQLOperator), []any{}, nil
}

func NewPredicateOperator[T any](sqlOperator string) PredicateOperator[T] {
	return PredicateOperator[T]{
		SQLOperator: sqlOperator,
	}
}
