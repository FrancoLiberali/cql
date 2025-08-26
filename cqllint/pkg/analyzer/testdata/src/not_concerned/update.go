package not_concerned

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

func testSetDynamicJoined() {
	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
	).Set(conditions.Phone.Name.Set().Eq(conditions.Brand.Name))
}

func testSetDynamicNotJoinedInSameLine() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(cql.String("asd")),
	).Set(conditions.Brand.Name.Set().Eq(conditions.City.Name)) // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
}

func testSetDynamicNotJoinedInDifferentLines() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(cql.String("asd")),
	).Set(conditions.Brand.Name.Set().Eq(
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	))
}

func testSetDynamicNotJoinedInMultiple() {
	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
	).Set(
		conditions.Phone.Name.Set().Eq(conditions.Brand.Name),
		conditions.Phone.Name.Set().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testSetDynamicJoinedFromVariable() {
	set := conditions.Brand.Name

	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
	).Set(conditions.Phone.Name.Set().Eq(set))
}

func testSetDynamicNotJoinedFromVariable() {
	set := conditions.City.Name

	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
	).Set(conditions.Phone.Name.Set().Eq(set)) // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
}

func testSetDynamicJoinedInVariable() {
	set := conditions.Phone.Name.Set().Eq(conditions.Brand.Name)

	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
	).Set(set)
}

func testSetDynamicNotJoinedInVariable() {
	set := conditions.Phone.Name.Set().Eq(conditions.City.Name) // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"

	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
	).Set(set)
}

func testSetDynamicJoinedInList() {
	sets := []*condition.Set[models.Phone]{
		conditions.Phone.Name.Set().Eq(conditions.Brand.Name),
	}

	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
	).Set(sets...)
}

func testSetDynamicNotJoinedInList() {
	sets := []*condition.Set[models.Phone]{
		conditions.Phone.Name.Set().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	}

	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
	).Set(sets...)
}

func testSetDynamicNotJoinedInListMultiple() {
	sets := []*condition.Set[models.Phone]{
		conditions.Phone.Name.Set().Eq(conditions.Brand.Name),
		conditions.Phone.Name.Set().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	}

	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
	).Set(sets...)
}

func testUpdateJoined() {
	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
	).Set(conditions.Phone.Name.Set().Eq(conditions.Brand.Name))
}

func testUpdateNotJoined() {
	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Set(conditions.Phone.Name.Set().Eq(conditions.Brand.Name))
}

func testSetDynamicSameModel() {
	cql.Update[models.Product](
		db,
		conditions.Product.String.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Set(
		conditions.Product.Int.Set().Eq(
			conditions.Product.IntPointer,
		),
	)
}

func testSetDynamicSameModelMultipleTimes() {
	cql.Update[models.Product](
		db,
		conditions.Product.String.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Set(
		conditions.Product.String.Set().Eq(
			conditions.Product.String2,
		),
		conditions.Product.Int.Set().Eq(
			conditions.Product.IntPointer,
		),
	)
}

func testSetDynamicJoinedModel() {
	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Set(
		conditions.Phone.Name.Set().Eq(conditions.Brand.Name),
	)
}

func testSetDynamicNestedJoinedModel() {
	cql.Update[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Set(
		conditions.Child.Name.Set().Eq(conditions.Parent1.Name),
		conditions.Child.Number.Set().Eq(conditions.ParentParent.Number),
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

func testSetMultipleDynamicNotJoined() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(cql.String("asd")),
	).SetMultiple(
		conditions.Brand.Name.Set().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testSetMultipleMainModel() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).SetMultiple(
		conditions.Brand.Name.Set().Eq(cql.String("asd")),
	)
}

func testSetMultipleMainModelMultipleTimes() {
	cql.Update[models.Product](
		db,
		conditions.Product.String.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).SetMultiple(
		conditions.Product.String.Set().Eq(cql.String("asd")),
		conditions.Product.Int.Set().Eq(cql.Int(1)),
	)
}

func testSetMultipleJoinedModel() {
	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
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
		conditions.Child.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
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
		conditions.City.Name.Concat("asd"), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	))
}

func testSetDynamicNotJoinedWithTwoFunction() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(cql.String("asd")),
	).Set(conditions.Brand.Name.Set().Eq(
		conditions.City.Name.Concat("asd").Concat("asd"), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	))
}
