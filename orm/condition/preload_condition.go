package condition

import (
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/query"
)

// Condition used to the preload the attributes of a model
type preloadCondition[T model.Model] struct {
	Fields []query.IFieldIdentifier
}

func (condition preloadCondition[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition preloadCondition[T]) ApplyTo(query *query.GormQuery, table query.Table) error {
	for _, fieldID := range condition.Fields {
		query.AddSelect(table, fieldID)
	}

	return nil
}

// Condition used to the preload the attributes of a model
func NewPreloadCondition[T model.Model](fields ...query.IFieldIdentifier) Condition[T] {
	return preloadCondition[T]{
		Fields: fields,
	}
}
