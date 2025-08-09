package condition

// Condition used to returns an error when the query is executed
type invalidCondition[T any] struct {
	Err error
}

func (condition invalidCondition[T]) interfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition invalidCondition[T]) applyTo(_ *GormQuery, _ Table) error {
	return condition.Err
}

func (condition invalidCondition[T]) getSQL(_ *GormQuery, _ Table) (string, []any, error) {
	return "", nil, condition.Err
}

func (condition invalidCondition[T]) affectsDeletedAt() bool {
	return false
}

// Condition used to returns an error when the query is executed
func newInvalidCondition[T any](err error) invalidCondition[T] {
	return invalidCondition[T]{
		Err: err,
	}
}
