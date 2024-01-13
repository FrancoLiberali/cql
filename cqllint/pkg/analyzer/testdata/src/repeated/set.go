package not_concerned

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
	"gorm.io/gorm"
)

var db *gorm.DB

func testSetRepeated() {
	cql.Update[models.Product](
		db,
		conditions.Product.Int.Is().Eq(0),
	).Set(
		conditions.Product.Int.Set().Eq(1), // want "conditions.Product.Int is repeated"
		conditions.Product.Int.Set().Eq(2), // want "conditions.Product.Int is repeated"
	)
}

func testSetNotRepeated() {
	cql.Update[models.Product](
		db,
		conditions.Product.Int.Is().Eq(0),
	).Set(
		conditions.Product.Int.Set().Eq(2),
	)
}

func testSetDynamicRepeated() {
	cql.Update[models.Product](
		db,
		conditions.Product.Int.Is().Eq(0),
	).Set(
		conditions.Product.Int.Set().Dynamic(conditions.Product.IntPointer.Value()), // want "conditions.Product.Int is repeated"
		conditions.Product.Int.Set().Dynamic(conditions.Product.IntPointer.Value()), // want "conditions.Product.Int is repeated"
	)
}

func testSetDynamicNotRepeated() {
	cql.Update[models.Product](
		db,
		conditions.Product.Int.Is().Eq(0),
	).Set(
		conditions.Product.Int.Set().Dynamic(conditions.Product.IntPointer.Value()),
	)
}

func testSetMultipleRepeated() {
	cql.Update[models.Product](
		db,
		conditions.Product.Int.Is().Eq(0),
	).SetMultiple(
		conditions.Product.Int.Set().Eq(1), // want "conditions.Product.Int is repeated"
		conditions.Product.Int.Set().Eq(2), // want "conditions.Product.Int is repeated"
	)
}

func testSetMultipleNotRepeated() {
	cql.Update[models.Product](
		db,
		conditions.Product.Int.Is().Eq(0),
	).SetMultiple(
		conditions.Product.Int.Set().Eq(2),
	)
}
