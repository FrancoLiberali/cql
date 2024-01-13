package cql

import (
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/preload"
)

var (
	// query

	ErrFieldModelNotConcerned   = condition.ErrFieldModelNotConcerned
	ErrAppearanceMustBeSelected = condition.ErrAppearanceMustBeSelected
	ErrFieldIsRepeated          = condition.ErrFieldIsRepeated

	// crud

	ErrMoreThanOneObjectFound = condition.ErrMoreThanOneObjectFound
	ErrObjectNotFound         = condition.ErrObjectNotFound

	// database

	ErrUnsupportedByDatabase = condition.ErrUnsupportedByDatabase
	ErrOrderByMustBeCalled   = condition.ErrOrderByMustBeCalled

	// preload

	ErrOnlyPreloadsAllowed = condition.ErrOnlyPreloadsAllowed
	ErrRelationNotLoaded   = preload.ErrRelationNotLoaded
)
