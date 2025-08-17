package condition

import (
	"errors"
	"fmt"

	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

var (
	// query

	ErrFieldModelNotConcerned   = errors.New("field's model is not concerned by the query (not joined)")
	ErrAppearanceMustBeSelected = errors.New("field's model appears more than once, select which one you want to use with Appearance")
	ErrAppearanceOutOfRange     = errors.New("selected appearance is bigger than field's model number of appearances")

	// conditions

	ErrEmptyConditions = errors.New("at least one condition is required")

	// crud

	ErrMoreThanOneObjectFound = errors.New("found more that one object that meet the requested conditions")
	ErrObjectNotFound         = errors.New("no object exists that meets the requested conditions")

	// database

	ErrUnsupportedByDatabase = errors.New("method not supported by database")
	ErrOrderByMustBeCalled   = errors.New("order by must be called before limit in an update statement")

	// preload

	ErrOnlyPreloadsAllowed = errors.New("only conditions that do a preload are allowed")
)

func methodError(err error, method string) error {
	return fmt.Errorf("%w; method: %s", err, method)
}

func fieldModelNotConcernedError(field IField) error {
	return fmt.Errorf("%w; not concerned model: %s",
		ErrFieldModelNotConcerned,
		field.getModelType(),
	)
}

func fieldModelError(err error, field IField) error {
	return fmt.Errorf("%w; model: %s", err, field.getModelType())
}

func appearanceMustBeSelectedError(field IField) error {
	return fieldModelError(ErrAppearanceMustBeSelected, field)
}

func appearanceOutOfRangeError(field IField) error {
	return fieldModelError(ErrAppearanceOutOfRange, field)
}

func preloadsInReturningNotAllowed(dialector sql.Dialector) error {
	return fmt.Errorf("%w; preloads in returning are not allowed for database: %s",
		ErrUnsupportedByDatabase,
		dialector,
	)
}

func conditionOperatorError[TObject model.Model, TAtribute any](operatorErr error, condition fieldCondition[TObject, TAtribute]) error {
	return fmt.Errorf(
		"%w; model: %T, field: %s",
		operatorErr,
		*new(TObject),
		condition.FieldIdentifier.fieldName(),
	)
}

func emptyConditionsError[T model.Model](connector sql.Operator) error {
	return fmt.Errorf(
		"%w; connector: %s; model: %T",
		ErrEmptyConditions,
		connector.Name(),
		*new(T),
	)
}

func onlyPreloadsAllowedError[T model.Model](fieldName string) error {
	return fmt.Errorf(
		"%w; model: %T, field: %s",
		ErrOnlyPreloadsAllowed,
		*new(T),
		fieldName,
	)
}

func operatorError(err error, sqlOperator sql.Operator) error {
	return fmt.Errorf("%w; operator: %s", err, sqlOperator.Name())
}

func functionError(err error, function sql.FunctionByDialector) error {
	return fmt.Errorf("%w; function: %s", err, function.Name)
}
