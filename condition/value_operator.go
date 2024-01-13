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
}

func NewValueOperator[T any](sqlOperator sql.Operator, value any) *ValueOperator[T] {
	return new(ValueOperator[T]).AddOperation(sqlOperator, value)
}

func (operator ValueOperator[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Operator[T]
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

		iValue, isIValue := operation.Value.(IValue)
		if isIValue {
			// if the value of the operation is a field,
			// verify that this field is concerned by the query
			// (a join was performed with the model to which this field belongs)
			// and get the alias of the table of this model.
			field := iValue.getField()

			modelTable, err := getModelTable(query, field, field.getAppearance(), sqlOperator)
			if err != nil {
				return "", nil, err
			}

			valueSQL, valueValues, err := iValue.toSQL(query, modelTable)
			if err != nil {
				return "", nil, err
			}

			operationString += fmt.Sprintf(
				" %s %s",
				sqlOperator,
				valueSQL,
			)

			values = append(values, valueValues...)
		} else {
			operationString += " " + sqlOperator.String() + " ?"
			values = append(values, operation.Value)
		}
	}

	return operationString, values, nil
}

func getModelTable(query *GormQuery, field IField, appearance int, sqlOperator sql.Operator) (Table, error) {
	table, err := query.GetModelTable(field, appearance)
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
		}
	case map[sql.Dialector]sql.Operator:
		newOperation = operation{
			Value:                  value,
			SQLOperatorByDialector: sqlOperatorTyped,
		}
	default:
		return operator
	}

	operator.Operations = append(operator.Operations, newOperation)

	return operator
}
