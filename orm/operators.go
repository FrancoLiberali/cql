package orm

// Comparison Operators
// ref: https://www.postgresql.org/docs/current/functions-comparison.html

// EqualTo
// IsNotDistinct must be used in cases where value can be NULL
func Eq[T any](value T) Operator[T] {
	return NewValueOperator[T]("=", value)
}

// NotEqualTo
// IsDistinct must be used in cases where value can be NULL
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

// Comparison Predicates
// refs: https://www.postgresql.org/docs/current/functions-comparison.html#FUNCTIONS-COMPARISON-PRED-TABLE

func IsNull[T any]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS NULL")
}

func IsNotNull[T any]() PredicateOperator[T] {
	return NewPredicateOperator[T]("IS NOT NULL")
}

// Boolean Comparison Predicates

func IsTrue() PredicateOperator[bool] {
	return NewPredicateOperator[bool]("IS TRUE")
}

func IsNotTrue() PredicateOperator[bool] {
	return NewPredicateOperator[bool]("IS NOT TRUE")
}

func IsFalse() PredicateOperator[bool] {
	return NewPredicateOperator[bool]("IS FALSE")
}

func IsNotFalse() PredicateOperator[bool] {
	return NewPredicateOperator[bool]("IS NOT FALSE")
}

func IsUnknown() PredicateOperator[bool] {
	return NewPredicateOperator[bool]("IS UNKNOWN")
}

func IsNotUnknown() PredicateOperator[bool] {
	return NewPredicateOperator[bool]("IS NOT UNKNOWN")
}

func IsDistinct[T any](value T) ValueOperator[T] {
	return NewValueOperator[T]("IS DISTINCT FROM", value)
}

func IsNotDistinct[T any](value T) ValueOperator[T] {
	return NewValueOperator[T]("IS NOT DISTINCT FROM", value)
}

// Row and Array Comparisons

func ArrayIn[T any](values ...T) ValueOperator[T] {
	return NewValueOperator[T]("IN", values)
}

func ArrayNotIn[T any](values ...T) ValueOperator[T] {
	return NewValueOperator[T]("NOT IN", values)
}

// Pattern Matching

type LikeOperator struct {
	ValueOperator[string]
}

func NewLikeOperator(sqlOperator string, pattern string) LikeOperator {
	return LikeOperator{
		ValueOperator: NewValueOperator[string](sqlOperator, pattern),
	}
}

func (operator LikeOperator) Escape(escape rune) ValueOperator[string] {
	return operator.AddOperation("ESCAPE", string(escape))
}

// Patterns:
//   - An underscore (_) in pattern stands for (matches) any single character.
//   - A percent sign (%) matches any sequence of zero or more characters.
//
// ref: https://www.postgresql.org/docs/current/functions-matching.html#FUNCTIONS-LIKE
func Like(pattern string) LikeOperator {
	return NewLikeOperator("LIKE", pattern)
}
