package cql

import (
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/preload"
)

var (
	// query
	ErrFieldModelNotConcerned = condition.ErrFieldModelNotConcerned
	ErrJoinMustBeSelected     = condition.ErrJoinMustBeSelected

	// conditions
	ErrOnlyPreloadsAllowed = condition.ErrOnlyPreloadsAllowed

	// crud
	ErrMoreThanOneObjectFound = condition.ErrMoreThanOneObjectFound
	ErrObjectNotFound         = condition.ErrObjectNotFound

	// database
	ErrUnsupportedByDatabase = condition.ErrUnsupportedByDatabase
	ErrOrderByMustBeCalled   = condition.ErrOrderByMustBeCalled

	// preload
	ErrRelationNotLoaded = preload.ErrRelationNotLoaded
)
