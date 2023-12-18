package condition

import (
	"fmt"

	"github.com/ditrit/badaas/orm/errors"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/sql"
)

func conditionOperatorError[TObject model.Model, TAtribute any](operatorErr error, condition fieldCondition[TObject, TAtribute]) error {
	return fmt.Errorf(
		"%w; model: %T, field: %s",
		operatorErr,
		*new(TObject),
		condition.FieldIdentifier.Field,
	)
}

func emptyConditionsError[T model.Model](connector sql.Operator) error {
	return fmt.Errorf(
		"%w; connector: %s; model: %T",
		errors.ErrEmptyConditions,
		connector.Name(),
		*new(T),
	)
}

func onlyPreloadsAllowedError[T model.Model](fieldName string) error {
	return fmt.Errorf(
		"%w; model: %T, field: %s",
		errors.ErrOnlyPreloadsAllowed,
		*new(T),
		fieldName,
	)
}
