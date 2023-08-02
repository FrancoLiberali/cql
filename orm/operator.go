package orm

type Operator[T any] interface {
	// Transform the Operator to a SQL string and a list of values to use in the query
	// columnName is used by the operator to determine which is the objective column.
	ToSQL(columnName string) (string, []any, error)

	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T],
	// since if no method receives by parameter a type T,
	// any other Operator[T2] would also be considered a Operator[T].
	InterfaceVerificationMethod(T)
}

// Operator that compares the value of the column against a fixed value
// If Operations has multiple entries, operations will be nested
// Example (single): value = v1
// Example (multi): value LIKE v1 ESCAPE v2
type ValueOperator[T any] struct {
	Operations []operation
}

type operation struct {
	SQLOperator string
	Value       any
}

func (operator ValueOperator[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T]
}

func (operator ValueOperator[T]) ToSQL(columnName string) (string, []any, error) {
	operatorString := columnName
	values := []any{}

	for _, operation := range operator.Operations {
		operatorString += " " + operation.SQLOperator + " ?"
		values = append(values, operation.Value)
	}

	return operatorString, values, nil
}

func NewValueOperator[T any](sqlOperator string, value any) ValueOperator[T] {
	operator := ValueOperator[T]{}

	return operator.AddOperation(sqlOperator, value)
}

func (operator *ValueOperator[T]) AddOperation(sqlOperator string, value any) ValueOperator[T] {
	operator.Operations = append(
		operator.Operations,
		operation{
			Value:       value,
			SQLOperator: sqlOperator,
		},
	)

	return *operator
}
