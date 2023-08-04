package orm

import (
	"errors"
	"fmt"

	"github.com/ditrit/badaas/orm/sql"
)

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
