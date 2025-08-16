package condition

// Condition used to returns an error when the query is executed
type invalidCondition[T any] struct {
	Err error
}

//nolint:unused // This method is necessary to get the compiler to verify that an object is of type Condition[T]
func (condition invalidCondition[T]) interfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

//nolint:unused // This method is necessary to get the compiler to verify that an object is of type Condition[T]
func (condition invalidCondition[T]) applyTo(_ *CQLQuery, _ Table) error {
	return condition.Err
}

//nolint:unused // This method is necessary to get the compiler to verify that an object is of type Condition[T]
func (condition invalidCondition[T]) getSQL(_ *CQLQuery, _ Table) (string, []any, error) {
	return "", nil, condition.Err
}

//nolint:unused // This method is necessary to get the compiler to verify that an object is of type Condition[T]
func (condition invalidCondition[T]) affectsDeletedAt() bool {
	return false
}

// Condition used to returns an error when the query is executed
func newInvalidCondition[T any](err error) invalidCondition[T] {
	return invalidCondition[T]{
		Err: err,
	}
}
