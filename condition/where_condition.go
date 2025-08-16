package condition

import (
	"github.com/FrancoLiberali/cql/model"
)

// Conditions that can be used in a where clause
// (or in a on of a join)
type WhereCondition[T model.Model] interface {
	Condition[T]

	// Get the sql string and values to use in the query
	getSQL(query *CQLQuery, table Table) (string, []any, error)

	// Returns true if the DeletedAt column if affected by the condition
	// If no condition affects the DeletedAt, the verification that it's null will be added automatically
	affectsDeletedAt() bool
}

// apply WhereCondition of any type on the query
func ApplyWhereCondition[T model.Model](condition WhereCondition[T], query *CQLQuery, table Table) error {
	sql, values, err := condition.getSQL(query, table)
	if err != nil {
		return err
	}

	if condition.affectsDeletedAt() {
		query.Unscoped()
	}

	query.Where(
		sql,
		values...,
	)

	return nil
}
