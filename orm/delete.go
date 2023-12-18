package orm

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/orm/cql"
	"github.com/FrancoLiberali/cql/orm/model"
)

// Create a Delete to which the conditions are applied inside transaction tx
func Delete[T model.Model](tx *gorm.DB, conditions ...cql.Condition[T]) *cql.Delete[T] {
	return cql.NewDelete(tx, conditions...)
}
