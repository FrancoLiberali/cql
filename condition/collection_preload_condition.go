package condition

import (
	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/model"
)

// Condition used to the preload a collection of models of a model
type collectionPreloadCondition[T1, T2 model.Model] struct {
	CollectionField string
	NestedPreloads  []JoinCondition[T2]

	// hasManyLoader is the per-relation fast-scan post-mount loader, wired
	// by cql-gen output. When non-nil, applyTo registers it on the query
	// and skips the gorm.Preload call — so the runtime materializes the
	// children via Query[Child] (fast path) instead of gorm's reflective
	// preload pipeline.
	hasManyLoader *HasManyLoader[T1, T2]
	// parentScanner is the per-model Scanner[T1] (Company). Surfaced via
	// getScannerErased so Query[T1].resolveScanner can find it when the
	// Collection.Preload() is the only condition in the query.
	parentScanner *Scanner[T1]
}

// getScannerErased lets resolveScanner pick up the parent's per-model
// scanner even when no Field-based condition is in the query.
func (condition collectionPreloadCondition[T1, T2]) getScannerErased() any {
	if condition.parentScanner == nil {
		return nil
	}

	return condition.parentScanner
}

func (condition collectionPreloadCondition[T1, T2]) interfaceVerificationMethod(_ T1) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T1]
}

func (condition collectionPreloadCondition[T1, T2]) applyTo(query *CQLQuery, _ Table) error {
	// Fast path: generated loader available. Register the loader for
	// findWith/scanOne to invoke post-scan. We ALSO call gorm.Preload so
	// non-Find paths (UPDATE/DELETE RETURNING + Preload) still get
	// children via gorm's reflective scan — findWith strips Preloads
	// before executing the main query so there's no double-load.
	if condition.hasManyLoader != nil {
		nested := pie.Map(
			condition.NestedPreloads,
			func(j JoinCondition[T2]) Condition[T2] { return j },
		)

		registerHasManyLoader[T1, T2](query, condition.hasManyLoader, nested)

		if len(condition.NestedPreloads) == 0 {
			query.Preload(condition.CollectionField)
		} else {
			condition.applyGormPreloadWithNested(query)
		}

		return nil
	}

	// Fallback: gorm reflective preload.
	if len(condition.NestedPreloads) == 0 {
		query.Preload(condition.CollectionField)
		return nil
	}

	query.Preload(
		condition.CollectionField,
		func(db *gorm.DB) *gorm.DB {
			preloadsAsCondition := pie.Map(
				condition.NestedPreloads,
				func(joinCondition JoinCondition[T2]) Condition[T2] {
					return joinCondition
				},
			)

			preloadQuery, err := ApplyConditions[T2](db, preloadsAsCondition)
			if err != nil {
				_ = db.AddError(err)
				return db
			}

			return preloadQuery.gormDB
		},
	)

	return nil
}

// applyGormPreloadWithNested duplicates the fallback path's nested-preload
// gorm.Preload registration so non-Find call sites (Returning, etc.) still
// see the children.
func (condition collectionPreloadCondition[T1, T2]) applyGormPreloadWithNested(query *CQLQuery) {
	query.Preload(
		condition.CollectionField,
		func(db *gorm.DB) *gorm.DB {
			preloadsAsCondition := pie.Map(
				condition.NestedPreloads,
				func(joinCondition JoinCondition[T2]) Condition[T2] {
					return joinCondition
				},
			)

			preloadQuery, err := ApplyConditions[T2](db, preloadsAsCondition)
			if err != nil {
				_ = db.AddError(err)
				return db
			}

			return preloadQuery.gormDB
		},
	)
}

// Condition used to the preload a collection of models of a model.
//
// The optional hasManyLoader (variadic for back-compat with hand-written
// callers) is the generated fast-scan mounter. When supplied AND the user
// preloaded this collection, the runtime runs ONE child query through
// Query[T2] and groups results onto parents — no gorm reflective preload.
func NewCollectionPreloadCondition[T1, T2 model.Model](
	collectionField string,
	nestedPreloads []JoinCondition[T2],
	hasManyLoader ...*HasManyLoader[T1, T2],
) Condition[T1] {
	if pie.Any(nestedPreloads, func(nestedPreload JoinCondition[T2]) bool {
		return !nestedPreload.makesPreload() || nestedPreload.makesFilter()
	}) {
		return newInvalidCondition[T1](onlyPreloadsAllowedError[T1](collectionField))
	}

	c := collectionPreloadCondition[T1, T2]{
		CollectionField: collectionField,
		NestedPreloads:  nestedPreloads,
	}

	if len(hasManyLoader) > 0 {
		c.hasManyLoader = hasManyLoader[0]
	}

	return c
}
