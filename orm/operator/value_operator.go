package operator

import (
	"fmt"

	"github.com/ditrit/badaas/orm/query"
	"github.com/ditrit/badaas/orm/sql"
)

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

func NewValueOperator[T any](sqlOperator sql.Operator, value any) *ValueOperator[T] {
	return new(ValueOperator[T]).AddOperation(sqlOperator, value)
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

func (operator ValueOperator[T]) ToSQL(queryV *query.GormQuery, columnName string) (string, []any, error) {
	operationString := columnName
	values := []any{}

	// add each operation to the sql
	for _, operation := range operator.Operations {
		field, isField := operation.Value.(query.IFieldIdentifier)
		if isField {
			// if the value of the operation is a field,
			// verify that this field is concerned by the query
			// (a join was performed with the model to which this field belongs)
			// and get the alias of the table of this model.
			modelTable, err := getModelTable(queryV, field, operation.JoinNumber, operation.SQLOperator)
			if err != nil {
				return "", nil, err
			}

			operationString += fmt.Sprintf(
				" %s %s",
				operation.SQLOperator,
				field.ColumnSQL(queryV, modelTable),
			)
		} else {
			operationString += " " + operation.SQLOperator.String() + " ?"
			values = append(values, operation.Value)
		}
	}

	return operationString, values, nil
}

func getModelTable(queryV *query.GormQuery, field query.IFieldIdentifier, joinNumber int, sqlOperator sql.Operator) (query.Table, error) {
	modelTables := queryV.GetTables(field.GetModelType())
	if modelTables == nil {
		return query.Table{}, fieldModelNotConcernedError(field, sqlOperator)
	}

	if len(modelTables) == 1 {
		return modelTables[0], nil
	}

	if joinNumber == undefinedJoinNumber {
		return query.Table{}, joinMustBeSelectedError(field, sqlOperator)
	}

	return modelTables[joinNumber], nil
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
