package orm

import (
	"fmt"

	"github.com/ditrit/badaas/orm/sql"
)

type Operator[T any] interface {
	// Transform the Operator to a SQL string and a list of values to use in the query
	// columnName is used by the operator to determine which is the objective column.
	ToSQL(query *Query, columnName string) (string, []any, error)

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

const undefinedJoinNumber = -1

// Operator that compares the value of the column against a fixed value
// If Operations has multiple entries, operations will be nested
// Example (single): value = v1
// Example (multi): value LIKE v1 ESCAPE v2
type ValueOperator[T any] struct {
	Operations []operation
}

type operation struct {
	SQLOperator sql.Operator
	Value       any
	JoinNumber  int
}

func (operator ValueOperator[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T]
}

// Allows to choose which number of join use
// for the operation in position "operationNumber"
// when the value is a field and its model is joined more than once.
// Does nothing if the operationNumber is bigger than the amount of operations.
func (operator *ValueOperator[T]) SelectJoin(operationNumber, joinNumber uint) DynamicOperator[T] {
	if operationNumber >= uint(len(operator.Operations)) {
		return operator
	}

	operationSaved := operator.Operations[operationNumber]
	operationSaved.JoinNumber = int(joinNumber)
	operator.Operations[operationNumber] = operationSaved

	return operator
}

func (operator ValueOperator[T]) ToSQL(query *Query, columnName string) (string, []any, error) {
	operationString := columnName
	values := []any{}

	// add each operation to the sql
	for _, operation := range operator.Operations {
		field, isField := operation.Value.(iFieldIdentifier)
		if isField {
			// if the value of the operation is a field,
			// verify that this field is concerned by the query
			// (a join was performed with the model to which this field belongs)
			// and get the alias of the table of this model.
			modelTable, err := getModelTable(query, field, operation.JoinNumber, operation.SQLOperator)
			if err != nil {
				return "", nil, err
			}

			operationString += fmt.Sprintf(
				" %s %s",
				operation.SQLOperator,
				field.ColumnSQL(query, modelTable),
			)
		} else {
			operationString += " " + operation.SQLOperator.String() + " ?"
			values = append(values, operation.Value)
		}
	}

	return operationString, values, nil
}

func getModelTable(query *Query, field iFieldIdentifier, joinNumber int, sqlOperator sql.Operator) (Table, error) {
	modelTables := query.GetTables(field.GetModelType())
	if modelTables == nil {
		return Table{}, fieldModelNotConcernedError(field, sqlOperator)
	}

	if len(modelTables) == 1 {
		return modelTables[0], nil
	}

	if joinNumber == undefinedJoinNumber {
		return Table{}, joinMustBeSelectedError(field, sqlOperator)
	}

	return modelTables[joinNumber], nil
}

func NewValueOperator[T any](sqlOperator sql.Operator, value any) *ValueOperator[T] {
	return new(ValueOperator[T]).AddOperation(sqlOperator, value)
}

func (operator *ValueOperator[T]) AddOperation(sqlOperator sql.Operator, value any) *ValueOperator[T] {
	operator.Operations = append(
		operator.Operations,
		operation{
			Value:       value,
			SQLOperator: sqlOperator,
			JoinNumber:  undefinedJoinNumber,
		},
	)

	return operator
}

// Operator that verifies a predicate
// Example: value IS TRUE
type PredicateOperator[T any] struct {
	SQLOperator string
}

func (operator PredicateOperator[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T]
}

func (operator PredicateOperator[T]) ToSQL(_ *Query, columnName string) (string, []any, error) {
	return fmt.Sprintf("%s %s", columnName, operator.SQLOperator), []any{}, nil
}

func NewPredicateOperator[T any](sqlOperator string) PredicateOperator[T] {
	return PredicateOperator[T]{
		SQLOperator: sqlOperator,
	}
}
