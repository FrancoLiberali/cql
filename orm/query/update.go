package query

import (
	"github.com/ditrit/badaas/orm/model"
)

type ISet interface {
	Field() IFieldIdentifier
	Value() any
	JoinNumber() int
}

// TODO ver donde pongo esto
type Set[T model.Model] struct {
	fieldID    IFieldIdentifier
	value      any
	joinNumber int
}

func (set Set[T]) Field() IFieldIdentifier {
	return set.fieldID
}

func (set Set[T]) Value() any {
	return set.value
}

func (set Set[T]) JoinNumber() int {
	return set.joinNumber
}

// TODO ver donde pongo esto
// TODO nombre muy parecido
type FieldSet[TObject model.Model, TAttribute any] struct {
	FieldID FieldIdentifier[TAttribute]
}

func (set FieldSet[TObject, TAttribute]) Eq(value TAttribute) *Set[TObject] {
	return &Set[TObject]{
		fieldID: set.FieldID,
		value:   value,
	}
}

// joinNumber can be used to select the join in case the field is joined more than once
func (set FieldSet[TObject, TAttribute]) Dynamic(field FieldIdentifier[TAttribute], joinNumber ...uint) *Set[TObject] {
	return &Set[TObject]{
		fieldID:    set.FieldID,
		value:      field,
		joinNumber: GetJoinNumber(joinNumber),
	}
}

func (set FieldSet[TObject, TAttribute]) Unsafe(value any) *Set[TObject] {
	return &Set[TObject]{
		fieldID: set.FieldID,
		value:   value,
	}
}
