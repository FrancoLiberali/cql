package orm

import (
	"errors"
	"fmt"

	"github.com/ditrit/badaas/orm/sql"
)

// operators

var (
	ErrFieldModelNotConcerned = errors.New("field's model is not concerned by the query (not joined)")
	ErrJoinMustBeSelected     = errors.New("field's model is joined more than once, select which one you want to use with SelectJoin")
)

func OperatorError(err error, sqlOperator sql.Operator) error {
	return fmt.Errorf("%w; operator: %s", err, sqlOperator.Name())
}

func fieldModelNotConcernedError(field iFieldIdentifier, sqlOperator sql.Operator) error {
	return OperatorError(fmt.Errorf("%w; not concerned model: %s",
		ErrFieldModelNotConcerned,
		field.GetModelType(),
	), sqlOperator)
}

func joinMustBeSelectedError(field iFieldIdentifier, sqlOperator sql.Operator) error {
	return OperatorError(fmt.Errorf("%w; joined multiple times model: %s",
		ErrJoinMustBeSelected,
		field.GetModelType(),
	), sqlOperator)
}

// conditions

var (
	ErrEmptyConditions     = errors.New("condition must have at least one inner condition")
	ErrOnlyPreloadsAllowed = errors.New("only conditions that do a preload are allowed")
)

func conditionOperatorError[TObject Model, TAtribute any](operatorErr error, condition FieldCondition[TObject, TAtribute]) error {
	return fmt.Errorf(
		"%w; model: %T, field: %s",
		operatorErr,
		*new(TObject),
		condition.FieldIdentifier.Field,
	)
}

func emptyConditionsError[T Model](connector sql.Operator) error {
	return fmt.Errorf(
		"%w; connector: %s; model: %T",
		ErrEmptyConditions,
		connector.Name(),
		*new(T),
	)
}

func onlyPreloadsAllowedError[T Model](fieldName string) error {
	return fmt.Errorf(
		"%w; model: %T, field: %s",
		ErrOnlyPreloadsAllowed,
		*new(T),
		fieldName,
	)
}
