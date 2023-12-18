package condition

import (
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/query"
)

// Conditions that can be used in a where clause
// (or in a on of a join)
type WhereCondition[T model.Model] interface {
	Condition[T]

	// Get the sql string and values to use in the query
	GetSQL(query *query.GormQuery, table query.Table) (string, []any, error)

	// Returns true if the DeletedAt column if affected by the condition
	// If no condition affects the DeletedAt, the verification that it's null will be added automatically
	AffectsDeletedAt() bool
}

// apply WhereCondition of any type on the query
func ApplyWhereCondition[T model.Model](condition WhereCondition[T], query *query.GormQuery, table query.Table) error {
	sql, values, err := condition.GetSQL(query, table)
	if err != nil {
		return err
	}

	if condition.AffectsDeletedAt() {
		query.Unscoped()
	}

	query.Where(
		sql,
		values...,
	)

	return nil
}
