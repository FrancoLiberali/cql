package condition

import (
	"strings"

	"github.com/elliotchance/pie/v2"

	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

// Condition that connects multiple conditions.
// Example: condition1 AND condition2
type connectionCondition[T model.Model] struct {
	Connector  sql.Operator
	Conditions []WhereCondition[T]
}

func (condition connectionCondition[T]) interfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition connectionCondition[T]) applyTo(query *CQLQuery, table Table) error {
	return ApplyWhereCondition[T](condition, query, table)
}

func (condition connectionCondition[T]) getSQL(query *CQLQuery, table Table) (string, []any, error) {
	sqlStrings := []string{}
	values := []any{}

	for _, internalCondition := range condition.Conditions {
		internalSQLString, internalValues, err := internalCondition.getSQL(query, table)
		if err != nil {
			return "", nil, err
		}

		if internalSQLString != "" {
			sqlStrings = append(sqlStrings, internalSQLString)

			values = append(values, internalValues...)
		}
	}

	return connectSQLs(sqlStrings, condition.Connector.String()), values, nil
}

func connectSQLs(sqlStrings []string, connector string) string {
	if len(sqlStrings) > 0 {
		return "(" + strings.Join(
			sqlStrings,
			" "+connector+" ",
		) + ")"
	}

	return ""
}

func (condition connectionCondition[T]) affectsDeletedAt() bool {
	return pie.Any(condition.Conditions, func(internalCondition WhereCondition[T]) bool {
		return internalCondition.affectsDeletedAt()
	})
}

// Condition that connects multiple conditions.
// Example: condition1 AND condition2
func NewConnectionCondition[T model.Model](connector sql.Operator, conditions []WhereCondition[T]) WhereCondition[T] {
	return connectionCondition[T]{
		Connector:  connector,
		Conditions: conditions,
	}
}
