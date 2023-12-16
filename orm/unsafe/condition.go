package unsafe

import (
	"fmt"

	"github.com/ditrit/badaas/orm"
)

// Condition that can be used to express conditions that are not supported (yet?) by badaas-orm
// Example: table1.columnX = table2.columnY
type Condition[T orm.Model] struct {
	SQLCondition string
	Values       []any
}

func (condition Condition[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition Condition[T]) ApplyTo(query *orm.Query, table orm.Table) error {
	return orm.ApplyWhereCondition[T](condition, query, table)
}

func (condition Condition[T]) GetSQL(_ *orm.Query, table orm.Table) (string, []any, error) {
	return fmt.Sprintf(
		condition.SQLCondition,
		table.Alias,
	), condition.Values, nil
}

func (condition Condition[T]) AffectsDeletedAt() bool {
	return false
}

// Condition that can be used to express conditions that are not supported (yet?) by badaas-orm
// Example: table1.columnX = table2.columnY
func NewCondition[T orm.Model](condition string, values ...any) Condition[T] {
	return Condition[T]{
		SQLCondition: condition,
		Values:       values,
	}
}
