package condition

import (
	"strings"

	"github.com/elliotchance/pie/v2"

	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/query"
	"github.com/ditrit/badaas/orm/sql"
)

// Condition that connects multiple conditions.
// Example: condition1 AND condition2
type connectionCondition[T model.Model] struct {
	Connector  sql.Operator
	Conditions []WhereCondition[T]
}

func (condition connectionCondition[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition connectionCondition[T]) ApplyTo(query *query.Query, table query.Table) error {
	return ApplyWhereCondition[T](condition, query, table)
}

func (condition connectionCondition[T]) GetSQL(query *query.Query, table query.Table) (string, []any, error) {
	sqlStrings := []string{}
	values := []any{}

	for _, internalCondition := range condition.Conditions {
		internalSQLString, internalValues, err := internalCondition.GetSQL(query, table)
		if err != nil {
			return "", nil, err
		}

		sqlStrings = append(sqlStrings, internalSQLString)

		values = append(values, internalValues...)
	}

	return strings.Join(
		sqlStrings,
		" "+condition.Connector.String()+" ",
	), values, nil
}

func (condition connectionCondition[T]) AffectsDeletedAt() bool {
	return pie.Any(condition.Conditions, func(internalCondition WhereCondition[T]) bool {
		return internalCondition.AffectsDeletedAt()
	})
}

// Condition that connects multiple conditions.
// Example: condition1 AND condition2
func NewConnectionCondition[T model.Model](connector sql.Operator, conditions ...WhereCondition[T]) WhereCondition[T] {
	return connectionCondition[T]{
		Connector:  connector,
		Conditions: conditions,
	}
}
