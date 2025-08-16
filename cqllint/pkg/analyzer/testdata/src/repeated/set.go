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
		conditions.Product.Int.Is().Eq(cql.Int(0)),
	).Set(
		conditions.Product.Int.Set().Eq(cql.Int(1)), // want "conditions.Product.Int is repeated"
		conditions.Product.Int.Set().Eq(cql.Int(2)), // want "conditions.Product.Int is repeated"
	)
}

func testSetNotRepeated() {
	cql.Update[models.Product](
		db,
		conditions.Product.Int.Is().Eq(cql.Int(0)),
	).Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	)
}

func testSetDynamicRepeated() {
	cql.Update[models.Product](
		db,
		conditions.Product.Int.Is().Eq(cql.Int(0)),
	).Set(
		conditions.Product.Int.Set().Eq(conditions.Product.IntPointer.Value()), // want "conditions.Product.Int is repeated"
		conditions.Product.Int.Set().Eq(conditions.Product.IntPointer.Value()), // want "conditions.Product.Int is repeated"
	)
}

func testSetDynamicNotRepeated() {
	cql.Update[models.Product](
		db,
		conditions.Product.Int.Is().Eq(cql.Int(0)),
	).Set(
		conditions.Product.Int.Set().Eq(conditions.Product.IntPointer.Value()),
	)
}

func testSetMultipleRepeated() {
	cql.Update[models.Product](
		db,
		conditions.Product.Int.Is().Eq(cql.Int(0)),
	).SetMultiple(
		conditions.Product.Int.Set().Eq(cql.Int(1)), // want "conditions.Product.Int is repeated"
		conditions.Product.Int.Set().Eq(cql.Int(2)), // want "conditions.Product.Int is repeated"
	)
}

func testSetMultipleNotRepeated() {
	cql.Update[models.Product](
		db,
		conditions.Product.Int.Is().Eq(cql.Int(0)),
	).SetMultiple(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	)
}

func testSetDynamicSameValue() {
	cql.Update[models.Product](
		db,
		conditions.Product.Int.Is().Eq(cql.Int(0)),
	).Set(
		conditions.Product.Int.Set().Eq(conditions.Product.Int.Value()), // want "conditions.Product.Int is set to itself"
	)
}

func testSetDynamicSameValueWithFunction() {
	cql.Update[models.Product](
		db,
		conditions.Product.Int.Is().Eq(cql.Int(0)),
	).Set(
		conditions.Product.Int.Set().Eq(conditions.Product.Int.Value().Plus(1)),
	)
}
