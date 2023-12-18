package orm

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm/cql"
	"github.com/ditrit/badaas/orm/model"
)

// TODO null y zero para update

// Create a Update to which the conditions are applied inside transaction tx
func Update[T model.Model](tx *gorm.DB, conditions ...cql.Condition[T]) *cql.Update[T] {
	return cql.NewUpdate(tx, conditions...)
}
