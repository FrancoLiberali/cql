package condition

import (
	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

// Condition that contains a internal condition.
// Example: NOT (internal condition)
type containerCondition[T model.Model] struct {
	ConnectionCondition WhereCondition[T]
	Prefix              sql.Operator
}

func (condition containerCondition[T]) interfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition containerCondition[T]) applyTo(query *CQLQuery, table Table) error {
	return ApplyWhereCondition[T](condition, query, table)
}

func (condition containerCondition[T]) getSQL(query *CQLQuery, table Table) (string, []any, error) {
	sqlString, values, err := condition.ConnectionCondition.getSQL(query, table)
	if err != nil {
		return "", nil, err
	}

	sqlString = condition.Prefix.String() + " (" + sqlString + ")"

	return sqlString, values, nil
}

func (condition containerCondition[T]) affectsDeletedAt() bool {
	return condition.ConnectionCondition.affectsDeletedAt()
}

// Condition that contains a internal condition.
// Example: NOT (internal condition)
func NewContainerCondition[T model.Model](prefix sql.Operator, conditions []WhereCondition[T]) WhereCondition[T] {
	if len(conditions) == 0 {
		return newInvalidCondition[T](emptyConditionsError[T](prefix))
	}

	return containerCondition[T]{
		Prefix:              prefix,
		ConnectionCondition: And(conditions...),
	}
}
