package cql

import (
	"github.com/FrancoLiberali/cql/condition"
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
)
