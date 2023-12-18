package condition

import (
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/query"
	"github.com/ditrit/badaas/orm/sql"
)

// Condition that contains a internal condition.
// Example: NOT (internal condition)
type containerCondition[T model.Model] struct {
	ConnectionCondition WhereCondition[T]
	Prefix              sql.Operator
}

func (condition containerCondition[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition containerCondition[T]) ApplyTo(query *query.Query, table query.Table) error {
	return ApplyWhereCondition[T](condition, query, table)
}

func (condition containerCondition[T]) GetSQL(query *query.Query, table query.Table) (string, []any, error) {
	sqlString, values, err := condition.ConnectionCondition.GetSQL(query, table)
	if err != nil {
		return "", nil, err
	}

	sqlString = condition.Prefix.String() + " (" + sqlString + ")"

	return sqlString, values, nil
}

func (condition containerCondition[T]) AffectsDeletedAt() bool {
	return condition.ConnectionCondition.AffectsDeletedAt()
}

// Condition that contains a internal condition.
// Example: NOT (internal condition)
func NewContainerCondition[T model.Model](prefix sql.Operator, conditions ...WhereCondition[T]) WhereCondition[T] {
	if len(conditions) == 0 {
		return newInvalidCondition[T](emptyConditionsError[T](prefix))
	}

	return containerCondition[T]{
		Prefix:              prefix,
		ConnectionCondition: And(conditions...),
	}
}
