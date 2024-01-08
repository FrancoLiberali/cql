package condition

import (
	"github.com/FrancoLiberali/cql/model"
	"github.com/elliotchance/pie/v2"
)

type Collection[TObject model.Model, TAttribute model.Model] struct {
	name    string
	t1Field string
	t2Field string
}

// Preload collection of models
//
// nestedPreloads can be used to preload relations of the models inside the collection
func (collection Collection[TObject, TAttribute]) Preload(nestedPreloads ...JoinCondition[TAttribute]) Condition[TObject] {
	return NewCollectionPreloadCondition[TObject, TAttribute](collection.name, nestedPreloads)
}

// Any generates a condition that is true if at least one model in the collection fulfills the conditions
func (collection Collection[TObject, TAttribute]) Any(firstCondition WhereCondition[TAttribute], conditions ...WhereCondition[TAttribute]) WhereCondition[TObject] {
	return newExistsCondition[TObject, TAttribute](firstCondition, conditions, collection.name, collection.t1Field, collection.t2Field)
}

// None generates a condition that is true if no model in the collection fulfills the conditions
func (collection Collection[TObject, TAttribute]) None(firstCondition WhereCondition[TAttribute], conditions ...WhereCondition[TAttribute]) WhereCondition[TObject] {
	return Not[TObject](
		newExistsCondition[TObject, TAttribute](firstCondition, conditions, collection.name, collection.t1Field, collection.t2Field),
	)
}

// All generates a condition that is true if all models in the collection fulfill the conditions (or is empty)
func (collection Collection[TObject, TAttribute]) All(firstCondition WhereCondition[TAttribute], conditions ...WhereCondition[TAttribute]) WhereCondition[TObject] {
	return Not[TObject](
		newExistsCondition[TObject, TAttribute](
			Not[TAttribute](
				pie.Unshift(conditions, firstCondition)...,
			),
			[]WhereCondition[TAttribute]{},
			collection.name, collection.t1Field, collection.t2Field,
		),
	)
}

func NewCollection[TObject model.Model, TAttribute model.Model](name, t1Field, t2Field string) Collection[TObject, TAttribute] {
	return Collection[TObject, TAttribute]{
		name:    name,
		t1Field: t1Field,
		t2Field: t2Field,
	}
}
