package cql

import (
	"github.com/ditrit/badaas/orm/model"
)

// Condition used to the preload the attributes of a model
type preloadCondition[T model.Model] struct {
	Fields []IField
}

func (condition preloadCondition[T]) InterfaceVerificationMethod(_ T) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T]
}

func (condition preloadCondition[T]) ApplyTo(query *GormQuery, table Table) error {
	for _, fieldID := range condition.Fields {
		query.AddSelect(table, fieldID)
	}

	return nil
}

// Condition used to the preload the attributes of a model
func NewPreloadCondition[T model.Model](fields ...IField) Condition[T] {
	return preloadCondition[T]{
		Fields: fields,
	}
}
