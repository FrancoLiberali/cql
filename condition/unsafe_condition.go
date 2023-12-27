package condition

import (
	"fmt"

	"github.com/FrancoLiberali/cql/model"
)

// Condition that can be used to express conditions that are not supported (yet?) by cql
// Example: table1.columnX = table2.columnY
type UnsafeCondition[T model.Model] struct {
	SQLCondition string
	Values       []any
}

func (unsafeCondition UnsafeCondition[T]) interfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (unsafeCondition UnsafeCondition[T]) applyTo(query *GormQuery, table Table) error {
	return ApplyWhereCondition[T](unsafeCondition, query, table)
}

func (unsafeCondition UnsafeCondition[T]) getSQL(_ *GormQuery, table Table) (string, []any, error) {
	return fmt.Sprintf(
		unsafeCondition.SQLCondition,
		table.Alias,
	), unsafeCondition.Values, nil
}

func (unsafeCondition UnsafeCondition[T]) affectsDeletedAt() bool {
	return false
}
