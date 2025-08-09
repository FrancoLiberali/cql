package condition

import (
	"github.com/FrancoLiberali/cql/model"
)

// Condition used to the preload the attributes of a model
type preloadCondition[T model.Model] struct {
	Fields []IField
}

func (condition preloadCondition[T]) interfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition preloadCondition[T]) applyTo(query *GormQuery, table Table) error {
	for _, fieldID := range condition.Fields {
		query.AddSelectField(table, fieldID, true)
	}

	return nil
}

// Condition used to the preload the attributes of a model
func NewPreloadCondition[T model.Model](fields ...IField) Condition[T] {
	return preloadCondition[T]{
		Fields: fields,
	}
}
