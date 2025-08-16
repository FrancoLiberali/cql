package not_concerned

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

func testSetDynamicNotJoinedInSameLine() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(cql.String("asd")),
	).Set(conditions.Brand.Name.Set().Eq(conditions.City.Name.Value())) // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
}

func testSetDynamicNotJoinedInDifferentLines() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(cql.String("asd")),
	).Set(conditions.Brand.Name.Set().Eq(
		conditions.City.Name.Value(), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	))
}

func testSetDynamicMainModel() {
	cql.Update[models.Product](
		db,
		conditions.Product.String.Is().Eq(conditions.City.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Set(
		conditions.Product.Int.Set().Eq(
			conditions.Product.IntPointer.Value(),
		),
	)
}

func testSetDynamicMainModelMultipleTimes() {
	cql.Update[models.Product](
		db,
		conditions.Product.String.Is().Eq(conditions.City.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Set(
		conditions.Product.String.Set().Eq(
			conditions.Product.String2.Value(),
		),
		conditions.Product.Int.Set().Eq(
			conditions.Product.IntPointer.Value(),
		),
	)
}

func testSetDynamicJoinedModel() {
	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Name.Is().Eq(conditions.City.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Set(
		conditions.Phone.Name.Set().Eq(conditions.Brand.Name.Value()),
	)
}

func testSetDynamicNestedJoinedModel() {
	cql.Update[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Name.Is().Eq(conditions.City.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Set(
		conditions.Child.Name.Set().Eq(conditions.Parent1.Name.Value()),
		conditions.Child.Number.Set().Eq(conditions.ParentParent.Number.Value()),
	)
}

func testSetMultipleNotJoinedInSameLine() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(cql.String("asd")),
	).SetMultiple(conditions.City.Name.Set().Eq(cql.String("asd"))) // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
}

func testSetMultipleNotJoinedInDifferentLines() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(cql.String("asd")),
	).SetMultiple(
		conditions.City.Name.Set().Eq(cql.String("asd")), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testSetMultipleMainModel() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.City.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).SetMultiple(
		conditions.Brand.Name.Set().Eq(cql.String("asd")),
	)
}

func testSetMultipleMainModelMultipleTimes() {
	cql.Update[models.Product](
		db,
		conditions.Product.String.Is().Eq(conditions.City.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).SetMultiple(
		conditions.Product.String.Set().Eq(cql.String("asd")),
		conditions.Product.Int.Set().Eq(cql.Int(1)),
	)
}

func testSetMultipleJoinedModel() {
	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Name.Is().Eq(conditions.City.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).SetMultiple(
		conditions.Phone.Name.Set().Eq(cql.String("asd")),
		conditions.Brand.Name.Set().Eq(cql.String("asd")),
	)
}

func testSetMultipleNestedJoinedModel() {
	cql.Update[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Name.Is().Eq(conditions.City.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).SetMultiple(
		conditions.Parent1.Name.Set().Eq(cql.String("asd")),
		conditions.ParentParent.Name.Set().Eq(cql.String("asd")),
	)
}

func testSetDynamicNotJoinedWithFunction() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(cql.String("asd")),
	).Set(conditions.Brand.Name.Set().Eq(
		conditions.City.Name.Value().Concat("asd"), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	))
}

func testSetDynamicNotJoinedWithTwoFunction() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(cql.String("asd")),
	).Set(conditions.Brand.Name.Set().Eq(
		conditions.City.Name.Value().Concat("asd").Concat("asd"), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	))
}
