package unsafe

import (
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/model"
)

// Condition that can be used to express conditions that are not supported (yet?) by cql
// Example: table1.columnX = table2.columnY
func NewCondition[T model.Model](sqlCondition string, values ...any) condition.Condition[T] {
	return condition.UnsafeCondition[T]{
		SQLCondition: sqlCondition,
		Values:       values,
	}
}
