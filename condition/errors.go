package condition

import (
	"errors"
	"fmt"

	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
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

	ErrUnsupportedByDatabase = errors.New("method not supported by database")
	ErrOrderByMustBeCalled   = errors.New("order by must be called before limit in an update statement")
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

func joinMustBeSelectedError(field IField) error {
	return fmt.Errorf("%w; joined multiple times model: %s",
		ErrJoinMustBeSelected,
		field.getModelType(),
	)
}

func preloadsInReturningNotAllowed(dialector sql.Dialector) error {
	return fmt.Errorf("%w; preloads in returning are not allowed for database: %s",
		ErrUnsupportedByDatabase,
		dialector,
	)
}

//nolint:unused // is used
func conditionOperatorError[TObject model.Model, TAtribute any](operatorErr error, condition fieldCondition[TObject, TAtribute]) error {
	return fmt.Errorf(
		"%w; model: %T, field: %s",
		operatorErr,
		*new(TObject),
		condition.FieldIdentifier.name,
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
