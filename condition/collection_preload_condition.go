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
}

//nolint:unused // is used
func (condition collectionPreloadCondition[T1, T2]) interfaceVerificationMethod(_ T1) {
	// This method is necessary to get the compiler to verify
	// that an object is of type Condition[T1]
}

//nolint:unused // is used
func (condition collectionPreloadCondition[T1, T2]) applyTo(query *GormQuery, _ Table) error {
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

			return preloadQuery.GormDB
		},
	)

	return nil
}

// Condition used to the preload a collection of models of a model
func NewCollectionPreloadCondition[T1, T2 model.Model](
	collectionField string,
	nestedPreloads []JoinCondition[T2],
) Condition[T1] {
	if pie.Any(nestedPreloads, func(nestedPreload JoinCondition[T2]) bool {
		return !nestedPreload.makesPreload() || nestedPreload.makesFilter()
	}) {
		return newInvalidCondition[T1](onlyPreloadsAllowedError[T1](collectionField))
	}

	return collectionPreloadCondition[T1, T2]{
		CollectionField: collectionField,
		NestedPreloads:  nestedPreloads,
	}
}
