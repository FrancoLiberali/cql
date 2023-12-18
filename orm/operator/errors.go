package operator

import (
	"fmt"

	"github.com/ditrit/badaas/orm/errors"
	"github.com/ditrit/badaas/orm/query"
	"github.com/ditrit/badaas/orm/sql"
)

func operatorError(err error, sqlOperator sql.Operator) error {
	return fmt.Errorf("%w; operator: %s", err, sqlOperator.Name())
}

func fieldModelNotConcernedError(field query.IFieldIdentifier, sqlOperator sql.Operator) error {
	return operatorError(fmt.Errorf("%w; not concerned model: %s",
		errors.ErrFieldModelNotConcerned,
		field.GetModelType(),
	), sqlOperator)
}

func joinMustBeSelectedError(field query.IFieldIdentifier, sqlOperator sql.Operator) error {
	return operatorError(fmt.Errorf("%w; joined multiple times model: %s",
		errors.ErrJoinMustBeSelected,
		field.GetModelType(),
	), sqlOperator)
}
