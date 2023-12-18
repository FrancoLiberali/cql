package query

import (
	"fmt"

	"github.com/ditrit/badaas/orm/errors"
)

func fieldModelNotConcernedError(field IFieldIdentifier) error {
	return fmt.Errorf("%w; not concerned model: %s",
		errors.ErrFieldModelNotConcerned,
		field.GetModelType(),
	)
}

func joinMustBeSelectedError(field IFieldIdentifier) error {
	return fmt.Errorf("%w; joined multiple times model: %s",
		errors.ErrJoinMustBeSelected,
		field.GetModelType(),
	)
}
