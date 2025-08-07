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

func (condition fieldCondition[TObject, TAtribute]) interfaceVerificationMethod(_ TObject) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

// Returns a gorm Where condition that can be used
// to filter that the Field as a value of Value
//

func (condition fieldCondition[TObject, TAtribute]) applyTo(query *GormQuery, table Table) error {
	return ApplyWhereCondition[TObject](condition, query, table)
}

func (condition fieldCondition[TObject, TAtribute]) affectsDeletedAt() bool {
	return condition.FieldIdentifier.fieldName() == deletedAtField
}

func (condition fieldCondition[TObject, TAtribute]) getSQL(query *GormQuery, table Table) (string, []any, error) {
	sqlString, values, err := condition.Operator.ToSQL(
		query,
		condition.FieldIdentifier.columnSQL(query, table),
	)
	if err != nil {
		return "", nil, conditionOperatorError[TObject](err, condition)
	}

	return sqlString, values, nil
}

func NewFieldCondition[TObject model.Model, TAttribute any](
	fieldIdentifier Field[TObject, TAttribute],
	operator Operator[TAttribute],
) WhereCondition[TObject] {
	return &fieldCondition[TObject, TAttribute]{
		FieldIdentifier: fieldIdentifier,
		Operator:        operator,
	}
}
