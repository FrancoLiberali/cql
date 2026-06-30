package condition

import (
	"github.com/elliotchance/pie/v2"

	"github.com/FrancoLiberali/cql/model"
)

type Collection[TObject model.Model, TAttribute model.Model] struct {
	name    string
	t1Field string
	t2Field string

	// loader is the per-relation HasManyLoader wired by cql-gen output.
	// Forwarded to collectionPreloadCondition on .Preload() so the
	// runtime takes the fast path.
	loader *HasManyLoader[TObject, TAttribute]
	// parentScanner is the per-model Scanner for TObject (Company), wired
	// by cql-gen output. Needed so collectionPreloadCondition can surface
	// it via getScannerErased — without this, a query whose only
	// condition is a Collection.Preload() can't find a scanner and falls
	// back to gorm.
	parentScanner *Scanner[TObject]
}

// Preload collection of models
//
// nestedPreloads can be used to preload relations of the models inside the collection
func (collection Collection[TObject, TAttribute]) Preload(nestedPreloads ...JoinCondition[TAttribute]) Condition[TObject] {
	if collection.loader != nil {
		c := NewCollectionPreloadCondition[TObject, TAttribute](collection.name, nestedPreloads, collection.loader)
		// Inject the parent scanner so the resulting condition can
		// surface it via getScannerErased — necessary for the fast path
		// to engage when the Collection.Preload() is the only condition.
		if cpc, ok := c.(collectionPreloadCondition[TObject, TAttribute]); ok {
			cpc.parentScanner = collection.parentScanner
			return cpc
		}

		return c
	}

	return NewCollectionPreloadCondition[TObject, TAttribute](collection.name, nestedPreloads)
}

// Any generates a condition that is true if at least one model in the collection fulfills the conditions
func (collection Collection[TObject, TAttribute]) Any(
	firstCondition WhereCondition[TAttribute],
	conditions ...WhereCondition[TAttribute],
) WhereCondition[TObject] {
	return newExistsCondition[TObject, TAttribute](firstCondition, conditions, collection.name, collection.t1Field, collection.t2Field)
}

// None generates a condition that is true if no model in the collection fulfills the conditions
func (collection Collection[TObject, TAttribute]) None(
	firstCondition WhereCondition[TAttribute],
	conditions ...WhereCondition[TAttribute],
) WhereCondition[TObject] {
	return Not[TObject](
		newExistsCondition[TObject, TAttribute](firstCondition, conditions, collection.name, collection.t1Field, collection.t2Field),
	)
}

// All generates a condition that is true if all models in the collection fulfill the conditions (or is empty)
func (collection Collection[TObject, TAttribute]) All(
	firstCondition WhereCondition[TAttribute],
	conditions ...WhereCondition[TAttribute],
) WhereCondition[TObject] {
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

// NewCollection builds a Collection. The optional loader is the per-relation
// HasManyLoader wired by cql-gen output. When provided, .Preload() takes the
// fast path: one extra child SELECT through Query[Child] + a generated mounter,
// no gorm reflective preload.
func NewCollection[TObject model.Model, TAttribute model.Model](
	name, t1Field, t2Field string,
	loader ...*HasManyLoader[TObject, TAttribute],
) Collection[TObject, TAttribute] {
	c := Collection[TObject, TAttribute]{
		name:    name,
		t1Field: t1Field,
		t2Field: t2Field,
	}

	if len(loader) > 0 {
		c.loader = loader[0]
	}

	return c
}

// WithParentScanner attaches the parent model's Scanner so a query whose
// only condition is Collection.Preload() can still take the fast path.
// Generated init() calls this after the per-model scanner var is declared.
func (collection Collection[TObject, TAttribute]) WithParentScanner(s *Scanner[TObject]) Collection[TObject, TAttribute] {
	collection.parentScanner = s
	return collection
}
