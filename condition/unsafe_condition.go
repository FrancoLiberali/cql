package condition

import (
	"fmt"
	"strings"

	"github.com/FrancoLiberali/cql/model"
)

// Condition that can be used to express conditions that are not supported (yet?) by cql
// Example: table1.columnX = table2.columnY
type UnsafeCondition[T model.Model] struct {
	SQLCondition string
	Values       []any
}

//nolint:unused // is used
func (unsafeCondition UnsafeCondition[T]) interfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (unsafeCondition UnsafeCondition[T]) applyTo(query *GormQuery, table Table) error {
	return ApplyWhereCondition[T](unsafeCondition, query, table)
}

//nolint:unused // is used
func (unsafeCondition UnsafeCondition[T]) getSQL(_ *GormQuery, table Table) (string, []any, error) {
	if strings.Contains(unsafeCondition.SQLCondition, "%s") {
		return fmt.Sprintf(
			unsafeCondition.SQLCondition,
			table.Alias,
		), unsafeCondition.Values, nil
	}

	return unsafeCondition.SQLCondition, unsafeCondition.Values, nil
}

//nolint:unused // is used
func (unsafeCondition UnsafeCondition[T]) affectsDeletedAt() bool {
	return false
}
