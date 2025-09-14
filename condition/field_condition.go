package condition

import (
	"github.com/FrancoLiberali/cql/model"
)

const deletedAtField = "DeletedAt"

// Condition that verifies the value of a field,
// using the Operator
type fieldCondition[TObject model.Model, TAtribute any] struct {
	FieldIdentifier IField
	Operator        Operator[TAtribute]
}

func (condition fieldCondition[TObject, TAtribute]) interfaceVerificationMethod(_ TObject) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

// Returns a gorm Where condition that can be used
// to filter that the Field as a value of Value
//

func (condition fieldCondition[TObject, TAtribute]) applyTo(query *CQLQuery, table Table) error {
	return ApplyWhereCondition[TObject](condition, query, table)
}

func (condition fieldCondition[TObject, TAtribute]) affectsDeletedAt() bool {
	return condition.FieldIdentifier.fieldName() == deletedAtField
}

func (condition fieldCondition[TObject, TAtribute]) getSQL(query *CQLQuery, table Table) (string, []any, error) {
	fieldSQL, fieldValues, err := condition.FieldIdentifier.ToSQLForTable(query, table)
	if err != nil {
		return "", nil, err
	}

	sqlString, values, err := condition.Operator.ToSQL(
		query,
		fieldSQL,
	)
	if err != nil {
		return "", nil, conditionOperatorError[TObject](err, condition)
	}

	return sqlString, append(fieldValues, values...), nil
}

func NewFieldCondition[TObject model.Model, TAttribute any](
	fieldIdentifier IField,
	operator Operator[TAttribute],
) WhereCondition[TObject] {
	return &fieldCondition[TObject, TAttribute]{
		FieldIdentifier: fieldIdentifier,
		Operator:        operator,
	}
}
