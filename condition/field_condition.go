package condition

import (
	"github.com/FrancoLiberali/cql/model"
)

const deletedAtField = "DeletedAt"

// Condition that verifies the value of a field,
// using the Operator
type fieldCondition[TObject model.Model, TAtribute any] struct {
	FieldIdentifier Field[TObject, TAtribute]
	Operator        Operator[TAtribute]
}

func (condition fieldCondition[TObject, TAtribute]) InterfaceVerificationMethod(_ TObject) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

// Returns a gorm Where condition that can be used
// to filter that the Field as a value of Value
func (condition fieldCondition[TObject, TAtribute]) ApplyTo(query *GormQuery, table Table) error {
	return ApplyWhereCondition[TObject](condition, query, table)
}

func (condition fieldCondition[TObject, TAtribute]) AffectsDeletedAt() bool {
	return condition.FieldIdentifier.FieldName() == deletedAtField
}

func (condition fieldCondition[TObject, TAtribute]) GetSQL(query *GormQuery, table Table) (string, []any, error) {
	sqlString, values, err := condition.Operator.ToSQL(
		query,
		condition.FieldIdentifier.ColumnSQL(query, table),
	)
	if err != nil {
		return "", nil, conditionOperatorError[TObject](err, condition)
	}

	return sqlString, values, nil
}

func (condition *fieldCondition[TObject, TAtribute]) SelectJoin(operationNumber, joinNumber uint) DynamicCondition[TObject] {
	dynamicOperator, isDynamic := condition.Operator.(DynamicOperator[TAtribute])
	if isDynamic {
		condition.Operator = dynamicOperator.SelectJoin(operationNumber, joinNumber)
	}

	return condition
}

func NewFieldCondition[TObject model.Model, TAttribute any](
	fieldIdentifier Field[TObject, TAttribute],
	operator Operator[TAttribute],
) DynamicCondition[TObject] {
	return &fieldCondition[TObject, TAttribute]{
		FieldIdentifier: fieldIdentifier,
		Operator:        operator,
	}
}
