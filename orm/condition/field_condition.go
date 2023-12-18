package condition

import (
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/operator"
	"github.com/ditrit/badaas/orm/query"
)

const deletedAtField = "DeletedAt"

// Condition that verifies the value of a field,
// using the Operator
type fieldCondition[TObject model.Model, TAtribute any] struct {
	FieldIdentifier query.FieldIdentifier[TAtribute]
	Operator        operator.Operator[TAtribute]
}

func (condition fieldCondition[TObject, TAtribute]) InterfaceVerificationMethod(_ TObject) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

// Returns a gorm Where condition that can be used
// to filter that the Field as a value of Value
func (condition fieldCondition[TObject, TAtribute]) ApplyTo(query *query.Query, table query.Table) error {
	return ApplyWhereCondition[TObject](condition, query, table)
}

func (condition fieldCondition[TObject, TAtribute]) AffectsDeletedAt() bool {
	return condition.FieldIdentifier.Field == deletedAtField
}

func (condition fieldCondition[TObject, TAtribute]) GetSQL(query *query.Query, table query.Table) (string, []any, error) {
	sqlString, values, err := condition.Operator.ToSQL(
		query,
		condition.FieldIdentifier.ColumnSQL(query, table),
	)
	if err != nil {
		return "", nil, conditionOperatorError[TObject](err, condition)
	}

	return sqlString, values, nil
}

func NewFieldCondition[TObject model.Model, TAtribute any](
	fieldIdentifier query.FieldIdentifier[TAtribute],
	operator operator.Operator[TAtribute],
) WhereCondition[TObject] {
	return fieldCondition[TObject, TAtribute]{
		FieldIdentifier: fieldIdentifier,
		Operator:        operator,
	}
}
