package orm

// Comparison Operators
// ref: https://www.postgresql.org/docs/current/functions-comparison.html

// EqualTo
func Eq[T any](value T) Operator[T] {
	return NewValueOperator[T]("=", value)
}

// NotEqualTo
func NotEq[T any](value T) Operator[T] {
	return NewValueOperator[T]("<>", value)
}

// LessThan
func Lt[T any](value T) Operator[T] {
	return NewValueOperator[T]("<", value)
}

// LessThanOrEqualTo
func LtOrEq[T any](value T) Operator[T] {
	return NewValueOperator[T]("<=", value)
}

// GreaterThan
func Gt[T any](value T) Operator[T] {
	return NewValueOperator[T](">", value)
}

// GreaterThanOrEqualTo
func GtOrEq[T any](value T) Operator[T] {
	return NewValueOperator[T](">=", value)
}
