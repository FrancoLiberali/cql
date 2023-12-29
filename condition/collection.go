package condition

import (
	"github.com/FrancoLiberali/cql/model"
)

type Collection[TObject model.Model, TAttribute model.Model] struct {
	name string
}

// Preload collection of models
//
// nestedPreloads can be used to preload relations of the models inside the collection
func (collection Collection[TObject, TAttribute]) Preload(nestedPreloads ...JoinCondition[TAttribute]) Condition[TObject] {
	return NewCollectionPreloadCondition[TObject, TAttribute](collection.name, nestedPreloads)
}

func NewCollection[TObject model.Model, TAttribute model.Model](name string) Collection[TObject, TAttribute] {
	return Collection[TObject, TAttribute]{
		name: name,
	}
}
