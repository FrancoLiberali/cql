package a

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

func testSetMultipleNotJoinedInSameLine() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq("asd"),
	).SetMultiple(conditions.City.Name.Set().Eq("asd")) // want "conditions.City is not joined by the query"
}

func testSetMultipleNotJoinedInDifferentLines() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq("asd"),
	).SetMultiple(
		conditions.City.Name.Set().Eq("asd"), // want "conditions.City is not joined by the query"
	)
}

func testSetMultipleMainModel() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "conditions.City is not joined by the query"
	).SetMultiple(
		conditions.Brand.Name.Set().Eq("asd"),
	)
}

func testSetMultipleMainModelMultipleTimes() {
	cql.Update[models.Brand](
		db,
		conditions.Brand.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "conditions.City is not joined by the query"
	).SetMultiple(
		// TODO ver que no se repitan
		conditions.Brand.Name.Set().Eq("asd"),
		conditions.Brand.Name.Set().Eq("asd"),
	)
}

func testSetMultipleJoinedModel() {
	cql.Update[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "conditions.City is not joined by the query"
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
		conditions.Child.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "conditions.City is not joined by the query"
	).SetMultiple(
		conditions.Parent1.Name.Set().Eq("asd"),
		conditions.ParentParent.Name.Set().Eq("asd"),
	)
}
