package condition

import (
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/sql"
)

func And[T model.Model](conditions ...WhereCondition[T]) WhereCondition[T] {
	return NewConnectionCondition(sql.And, conditions...)
}
