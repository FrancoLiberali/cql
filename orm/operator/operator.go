package operator

import "github.com/ditrit/badaas/orm/query"

type Operator[T any] interface {
	// Transform the Operator to a SQL string and a list of values to use in the query
	// columnName is used by the operator to determine which is the objective column.
	ToSQL(query *query.Query, columnName string) (string, []any, error)

	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T],
	// since if no method receives by parameter a type T,
	// any other Operator[T2] would also be considered a Operator[T].
	InterfaceVerificationMethod(T)
}

type DynamicOperator[T any] interface {
	Operator[T]

	// Allows to choose which number of join use
	// for the value in position "valueNumber"
	// when the value is a field and its model is joined more than once.
	// Does nothing if the valueNumber is bigger than the amount of values.
	SelectJoin(valueNumber, joinNumber uint) DynamicOperator[T]
}
