package orm

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/orm/cql"
	"github.com/FrancoLiberali/cql/orm/model"
)

// Create a Query to which the conditions are applied inside transaction tx
func Query[T model.Model](tx *gorm.DB, conditions ...cql.Condition[T]) *cql.Query[T] {
	return cql.NewQuery(tx, conditions...)
}
