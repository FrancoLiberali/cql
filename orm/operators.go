package orm

// Comparison Operators
// ref: https://www.postgresql.org/docs/current/functions-comparison.html

// EqualTo
func Eq[T any](value T) Operator[T] {
	return NewValueOperator[T]("=", value)
}
