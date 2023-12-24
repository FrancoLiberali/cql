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
	ErrEmptyConditions     = condition.ErrEmptyConditions
	ErrOnlyPreloadsAllowed = condition.ErrOnlyPreloadsAllowed

	// crud
	ErrMoreThanOneObjectFound = condition.ErrMoreThanOneObjectFound
	ErrObjectNotFound         = condition.ErrObjectNotFound

	ErrUnsupportedByDatabase = condition.ErrUnsupportedByDatabase
	ErrOrderByMustBeCalled   = condition.ErrOrderByMustBeCalled

	// preload
	ErrRelationNotLoaded = preload.ErrRelationNotLoaded
)
