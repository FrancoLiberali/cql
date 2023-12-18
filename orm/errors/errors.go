package errors

import (
	"errors"
)

var (
	// query
	ErrFieldModelNotConcerned = errors.New("field's model is not concerned by the query (not joined)")
	ErrJoinMustBeSelected     = errors.New("field's model is joined more than once, select which one you want to use")

	// conditions
	ErrEmptyConditions     = errors.New("condition must have at least one inner condition")
	ErrOnlyPreloadsAllowed = errors.New("only conditions that do a preload are allowed")

	// crud
	ErrMoreThanOneObjectFound = errors.New("found more that one object that meet the requested conditions")
	ErrObjectNotFound         = errors.New("no object exists that meets the requested conditions")

	// preload
	ErrRelationNotLoaded = errors.New("relation not loaded")
)
