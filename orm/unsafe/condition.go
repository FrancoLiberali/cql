package unsafe

import (
	"fmt"

	"github.com/ditrit/badaas/orm/condition"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/query"
)

// Condition that can be used to express conditions that are not supported (yet?) by badaas-orm
// Example: table1.columnX = table2.columnY
type unsafeCondition[T model.Model] struct {
	SQLCondition string
	Values       []any
}

func (unsafeCondition unsafeCondition[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (unsafeCondition unsafeCondition[T]) ApplyTo(queryV *query.Query, table query.Table) error {
	return condition.ApplyWhereCondition[T](unsafeCondition, queryV, table)
}

func (unsafeCondition unsafeCondition[T]) GetSQL(_ *query.Query, table query.Table) (string, []any, error) {
	return fmt.Sprintf(
		unsafeCondition.SQLCondition,
		table.Alias,
	), unsafeCondition.Values, nil
}

func (unsafeCondition unsafeCondition[T]) AffectsDeletedAt() bool {
	return false
}

// Condition that can be used to express conditions that are not supported (yet?) by badaas-orm
// Example: table1.columnX = table2.columnY
func NewCondition[T model.Model](sqlCondition string, values ...any) condition.Condition[T] {
	return unsafeCondition[T]{
		SQLCondition: sqlCondition,
		Values:       values,
	}
}
