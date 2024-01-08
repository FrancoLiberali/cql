package a

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

func testOrderNotJoinedInSameLine() {
	cql.Query[models.Brand](
		db,
	).Descending(conditions.City.Name).Find() // want "conditions.City is not joined by the query"
}

func testOrderNotJoinedInDifferentLines() {
	cql.Query[models.Brand](
		db,
	).Descending(
		conditions.City.Name, // want "conditions.City is not joined by the query"
	).Find()
}

func testOrderMainModel() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "conditions.City is not joined by the query"
	).Descending(
		conditions.Brand.Name,
	).Find()
}

func testOrderJoinedModelWithoutConditions() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "conditions.City is not joined by the query"
	).Descending(
		conditions.Brand.Name,
	).Find()
}

func testOrderJoinedModelWithConditions() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq("asd"),
		),
		conditions.Phone.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "conditions.City is not joined by the query"
	).Descending(
		conditions.Brand.Name,
	).Find()
}

func testOrderJoinedModelWithoutConditionsWithPreload() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand().Preload(),
		conditions.Phone.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "conditions.City is not joined by the query"
	).Descending(
		conditions.Brand.Name,
	).Find()
}

func testOrderJoinedModelWithConditionsWithPreload() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq("asd"),
		).Preload(),
		conditions.Phone.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "conditions.City is not joined by the query"
	).Descending(
		conditions.Brand.Name,
	).Find()
}

func testOrderNestedJoinedModel() {
	cql.Query[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "conditions.City is not joined by the query"
	).Descending(
		conditions.ParentParent.Name,
	).Find()
}

func testOrderMainModelWithoutConditions() {
	cql.Query[models.Brand](
		db,
	).Descending(
		conditions.Brand.Name,
	).Limit(1).Find()
}

func testOrderMainModelWithLimitAfter() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "conditions.City is not joined by the query"
	).Descending(
		conditions.Brand.Name,
	).Limit(1).Find()
}

func testOrderMainModelWithLimitBefore() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "conditions.City is not joined by the query"
	).Limit(1).Descending(
		conditions.Brand.Name,
	).Find()
}

func testOrderNotJoinedWithLimitAfter() {
	cql.Query[models.Brand](
		db,
	).Descending(
		conditions.City.Name, // want "conditions.City is not joined by the query"
	).Limit(1).Find()
}

func testOrderNotJoinedWithLimitBefore() {
	cql.Query[models.Brand](
		db,
	).Limit(1).Descending(
		conditions.City.Name, // want "conditions.City is not joined by the query"
	).Find()
}
