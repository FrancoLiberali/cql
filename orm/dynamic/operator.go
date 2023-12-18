package dynamic

import (
	"github.com/ditrit/badaas/orm/operator"
	"github.com/ditrit/badaas/orm/query"
	"github.com/ditrit/badaas/orm/sql"
)

func newValueOperator[T any](
	sqlOperator sql.Operator,
	field query.FieldIdentifier[T],
) *operator.ValueOperator[T] {
	return operator.NewValueOperator[T](sqlOperator, field)
}
