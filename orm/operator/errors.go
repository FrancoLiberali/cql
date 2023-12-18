package operator

import (
	"fmt"

	"github.com/ditrit/badaas/orm/sql"
)

func operatorError(err error, sqlOperator sql.Operator) error {
	return fmt.Errorf("%w; operator: %s", err, sqlOperator.Name())
}
