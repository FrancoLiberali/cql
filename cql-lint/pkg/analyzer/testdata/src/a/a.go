package a

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

var db *gorm.DB

func testError() {
	cql.Query[models.Phone]( // want "conditions.City is not joined by the query"
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.IsDynamic().Eq(conditions.City.Name.Value()),
		),
	).Find()
}
