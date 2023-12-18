package cql

// Condition used to returns an error when the query is executed
type invalidCondition[T any] struct {
	Err error
}

func (condition invalidCondition[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition invalidCondition[T]) ApplyTo(_ *GormQuery, _ Table) error {
	return condition.Err
}

func (condition invalidCondition[T]) GetSQL(_ *GormQuery, _ Table) (string, []any, error) {
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
