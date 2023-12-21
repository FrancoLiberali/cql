package condition

import (
	"fmt"

	"github.com/FrancoLiberali/cql/sql"
)

// Operator that compares the value of the column against a fixed value
// If Operations has multiple entries, operations will be nested
// Example (single): value = v1
// Example (multi): value LIKE v1 ESCAPE v2
type ValueOperator[T any] struct {
	Operations []operation
	Modifier   map[sql.Dialector]string
}

type operation struct {
	SQLOperator            sql.Operator
	SQLOperatorByDialector map[sql.Dialector]sql.Operator
	Value                  any
	JoinNumber             int
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

func (operator ValueOperator[T]) ToSQL(query *GormQuery, columnName string) (string, []any, error) {
	operationString := columnName

	if modifier, isPresent := operator.Modifier[query.Dialector()]; isPresent {
		operationString = modifier + " " + columnName
	}

	values := []any{}

	// add each operation to the sql
	for _, operation := range operator.Operations {
		sqlOperator := operation.SQLOperator
		if operation.SQLOperatorByDialector != nil {
			sqlOperator = operation.SQLOperatorByDialector[query.Dialector()]
		}

		if !sqlOperator.Supports(query.Dialector()) {
			return "", nil, operatorError(ErrUnsupportedByDatabase, sqlOperator)
		}

		field, isField := operation.Value.(IField)
		if isField {
			// if the value of the operation is a field,
			// verify that this field is concerned by the query
			// (a join was performed with the model to which this field belongs)
			// and get the alias of the table of this model.
			modelTable, err := getModelTable(query, field, operation.JoinNumber, sqlOperator)
			if err != nil {
				return "", nil, err
			}

			operationString += fmt.Sprintf(
				" %s %s",
				sqlOperator,
				field.ColumnSQL(query, modelTable),
			)
		} else {
			operationString += " " + sqlOperator.String() + " ?"
			values = append(values, operation.Value)
		}
	}

	return operationString, values, nil
}

func getModelTable(query *GormQuery, field IField, joinNumber int, sqlOperator sql.Operator) (Table, error) {
	table, err := query.GetModelTable(field, joinNumber)
	if err != nil {
		return Table{}, operatorError(err, sqlOperator)
	}

	return table, nil
}

func (operator *ValueOperator[T]) AddOperation(sqlOperator any, value any) *ValueOperator[T] {
	var newOperation operation
	switch sqlOperatorTyped := sqlOperator.(type) {
	case sql.Operator:
		newOperation = operation{
			Value:       value,
			SQLOperator: sqlOperatorTyped,
			JoinNumber:  UndefinedJoinNumber,
		}
	case map[sql.Dialector]sql.Operator:
		newOperation = operation{
			Value:                  value,
			SQLOperatorByDialector: sqlOperatorTyped,
			JoinNumber:             UndefinedJoinNumber,
		}
	default:
		return operator
	}

	operator.Operations = append(operator.Operations, newOperation)

	return operator
}
