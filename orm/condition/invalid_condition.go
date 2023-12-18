package condition

import "github.com/ditrit/badaas/orm/query"

// Condition used to returns an error when the query is executed
type invalidCondition[T any] struct {
	Err error
}

func (condition invalidCondition[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition invalidCondition[T]) ApplyTo(_ *query.GormQuery, _ query.Table) error {
	return condition.Err
}

func (condition invalidCondition[T]) GetSQL(_ *query.GormQuery, _ query.Table) (string, []any, error) {
	return "", nil, condition.Err
}

func (condition invalidCondition[T]) AffectsDeletedAt() bool {
	return false
}

// Condition used to returns an error when the query is executed
func newInvalidCondition[T any](err error) invalidCondition[T] {
	return invalidCondition[T]{
		Err: err,
	}
}
