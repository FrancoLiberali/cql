package dynamic

import (
	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/sql"
)

func NewValueOperator[T any](
	sqlOperator sql.Operator,
	field orm.FieldIdentifier[T],
) *orm.ValueOperator[T] {
	return orm.NewValueOperator[T](sqlOperator, field)
}
