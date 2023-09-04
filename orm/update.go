package orm

import (
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/query"
)

// TODO ver donde pongo esto
type Set[T model.Model] struct {
	fieldID query.IFieldIdentifier
	value   any
}

// TODO ver donde pongo esto
type FieldSet[TObject model.Model, TAttribute any] struct {
	FieldID query.FieldIdentifier[TAttribute]
}

func (set FieldSet[TObject, TAttribute]) Eq(value TAttribute) *Set[TObject] {
	return &Set[TObject]{
		fieldID: set.FieldID,
		value:   value,
	}
}

func (set FieldSet[TObject, TAttribute]) Dynamic(field query.FieldIdentifier[TAttribute]) *Set[TObject] {
	// TODO falta ver el join
	return &Set[TObject]{
		fieldID: set.FieldID,
		value:   field,
	}
}

func (set FieldSet[TObject, TAttribute]) Unsafe(value any) *Set[TObject] {
	// TODO falta ver el join
	return &Set[TObject]{
		fieldID: set.FieldID,
		value:   value,
	}
}
