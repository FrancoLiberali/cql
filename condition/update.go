package condition

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

type Update[T model.Model] struct {
	OrderLimitReturning[T]
}

// Set allows updating multiple attributes of the same table.
func (update *Update[T]) Set(sets ...*Set[T]) (int64, error) {
	setsAsInterface := []ISet{}
	for _, set := range sets {
		setsAsInterface = append(setsAsInterface, set)
	}

	return update.unsafeSet(setsAsInterface, "Set")
}

// SetMultiple allows updating multiple tables in the same query.
//
// available for: mysql
func (update *Update[T]) SetMultiple(sets ...ISet) (int64, error) {
	methodName := "SetMultiple"

	if update.query.gormQuery.Dialector() != sql.MySQL {
		update.query.addError(methodError(ErrUnsupportedByDatabase, methodName))
	}

	return update.unsafeSet(sets, methodName)
}

func (update *Update[T]) unsafeSet(sets []ISet, methodName string) (int64, error) {
	if update.query.err != nil {
		return 0, update.query.err
	}

	updated, err := update.query.gormQuery.Update(sets)
	if err != nil {
		return 0, methodError(err, methodName)
	}

	return updated, nil
}

// Ascending specify an ascending order when updating models
//
// available for: mysql
func (update *Update[T]) Ascending(field IField) *Update[T] {
	update.OrderLimitReturning.Ascending(field)

	return update
}

// Descending specify a descending order when updating models
//
// available for: mysql
func (update *Update[T]) Descending(field IField) *Update[T] {
	update.OrderLimitReturning.Descending(field)

	return update
}

// Limit specify the number of models to be updated
//
// Limit conditions can be cancelled by using `Limit(-1)`
//
// available for: mysql
func (update *Update[T]) Limit(limit int) *Update[T] {
	update.OrderLimitReturning.Limit(limit)

	return update
}

// available for: postgres, sqlite, sqlserver
//
// warning: in sqlite preloads are not allowed
func (update *Update[T]) Returning(dest *[]T) *Update[T] {
	update.OrderLimitReturning.Returning(dest)

	return update
}

// Create a Update to which the conditions are applied inside transaction tx
func NewUpdate[T model.Model](tx *gorm.DB, conditions []Condition[T]) *Update[T] {
	var err error

	if len(conditions) == 0 {
		err = methodError(ErrEmptyConditions, "Update")
	}

	query := NewQuery(tx, conditions...)
	if err != nil {
		query.err = err
	}

	return &Update[T]{
		OrderLimitReturning: OrderLimitReturning[T]{
			query:         query,
			orderByCalled: false,
		},
	}
}

type ISet interface {
	getField() IField
	getValue() IValue
}

type Set[T model.Model] struct {
	field IField
	value IValue
}

func (set Set[T]) getField() IField {
	return set.field
}

func (set Set[T]) getValue() IValue {
	return set.value
}

type FieldSet[TModel model.Model, TAttribute any] struct {
	field UpdatableField[TModel, TAttribute]
}

func (set FieldSet[TModel, TAttribute]) Eq(value ValueOfType[TAttribute]) *Set[TModel] {
	return &Set[TModel]{
		field: set.field,
		value: value,
	}
}

func (set FieldSet[TModel, TAttribute]) Unsafe(value IValue) *Set[TModel] {
	return &Set[TModel]{
		field: set.field,
		value: unsafeValue{Value: value},
	}
}

type NullableFieldSet[TModel model.Model, TAttribute any] struct {
	FieldSet[TModel, TAttribute]
}

func (set NullableFieldSet[TModel, TAttribute]) Null() *Set[TModel] {
	return &Set[TModel]{
		field: set.field,
		value: nil,
	}
}

type NumericFieldSet[TModel model.Model, TAttribute int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64] struct {
	field NumericField[TModel, TAttribute]
}

func (set NumericFieldSet[TModel, TAttribute]) Eq(value ValueOfType[float64]) *Set[TModel] {
	return &Set[TModel]{
		field: set.field,
		value: value,
	}
}

func (set NumericFieldSet[TModel, TAttribute]) Unsafe(value IValue) *Set[TModel] {
	return &Set[TModel]{
		field: set.field,
		value: unsafeValue{Value: value},
	}
}
