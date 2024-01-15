package not_concerned

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

func testSetDynamicNotJoinedInSameLine() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq("asd"),
	).Set(conditions.Brand.Name.Set().Dynamic(conditions.City.Name.Value())) // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
}

func testSetDynamicNotJoinedInDifferentLines() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq("asd"),
	).Set(conditions.Brand.Name.Set().Dynamic(
		conditions.City.Name.Value(), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	))
}

func testSetDynamicMainModel() {
	cql.Update[models.Product](
		db,
		conditions.Product.String.IsDynamic().Eq(conditions.City.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Set(
		conditions.Product.Int.Set().Dynamic(
			conditions.Product.IntPointer.Value(),
		),
	)
}

func testSetDynamicMainModelMultipleTimes() {
	cql.Update[models.Product](
		db,
		conditions.Product.String.IsDynamic().Eq(conditions.City.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Set(
		conditions.Product.String.Set().Dynamic(
			conditions.Product.String2.Value(),
		),
		conditions.Product.Int.Set().Dynamic(
			conditions.Product.IntPointer.Value(),
		),
	)
}

func testSetDynamicJoinedModel() {
	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Set(
		conditions.Phone.Name.Set().Dynamic(conditions.Brand.Name.Value()),
	)
}

func testSetDynamicNestedJoinedModel() {
	cql.Update[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Set(
		conditions.Child.Name.Set().Dynamic(conditions.Parent1.Name.Value()),
		conditions.Child.Number.Set().Dynamic(conditions.ParentParent.Number.Value()),
	)
}

func testSetMultipleNotJoinedInSameLine() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq("asd"),
	).SetMultiple(conditions.City.Name.Set().Eq("asd")) // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
}

func testSetMultipleNotJoinedInDifferentLines() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq("asd"),
	).SetMultiple(
		conditions.City.Name.Set().Eq("asd"), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testSetMultipleMainModel() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).SetMultiple(
		conditions.Brand.Name.Set().Eq("asd"),
	)
}

func testSetMultipleMainModelMultipleTimes() {
	cql.Update[models.Product](
		db,
		conditions.Product.String.IsDynamic().Eq(conditions.City.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).SetMultiple(
		conditions.Product.String.Set().Eq("asd"),
		conditions.Product.Int.Set().Eq(1),
	)
}

func testSetMultipleJoinedModel() {
	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).SetMultiple(
		conditions.Phone.Name.Set().Eq("asd"),
		conditions.Brand.Name.Set().Eq("asd"),
	)
}

func testSetMultipleNestedJoinedModel() {
	cql.Update[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).SetMultiple(
		conditions.Parent1.Name.Set().Eq("asd"),
		conditions.ParentParent.Name.Set().Eq("asd"),
	)
}

func testSetDynamicNotJoinedWithFunction() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq("asd"),
	).Set(conditions.Brand.Name.Set().Dynamic(
		conditions.City.Name.Value().Concat("asd"), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	))
}

func testSetDynamicNotJoinedWithTwoFunction() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq("asd"),
	).Set(conditions.Brand.Name.Set().Dynamic(
		conditions.City.Name.Value().Concat("asd").Concat("asd"), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	))
}
