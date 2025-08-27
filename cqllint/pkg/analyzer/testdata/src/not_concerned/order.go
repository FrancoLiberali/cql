package not_concerned

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

func testOrderNotJoinedInSameLine() {
	cql.Query[models.Brand](
		db,
	).Descending(conditions.City.Name).Find() // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
}

func testOrderNotJoinedInDifferentLines() {
	cql.Query[models.Brand](
		db,
	).Descending(
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testOrderMainModel() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Descending(
		conditions.Brand.Name,
	).Find()
}

func testOrderJoinedInVariable() {
	value := conditions.Brand.Name

	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
	).Descending(
		value,
	).Find()
}

func testOrderNotJoinedInVariable() {
	value := conditions.City.Name // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"

	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
	).Descending(
		value,
	).Find()
}

func testOrderJoinedModelWithoutConditions() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Descending(
		conditions.Brand.Name,
	).Find()
}

func testOrderJoinedModelWithConditions() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("asd")),
		),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Descending(
		conditions.Brand.Name,
	).Find()
}

func testOrderJoinedModelWithoutConditionsWithPreload() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand().Preload(),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Descending(
		conditions.Brand.Name,
	).Find()
}

func testOrderJoinedModelWithConditionsWithPreload() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("asd")),
		).Preload(),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
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
		conditions.Child.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
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
		conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Descending(
		conditions.Brand.Name,
	).Limit(1).Find()
}

func testOrderMainModelWithLimitBefore() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Limit(1).Descending(
		conditions.Brand.Name,
	).Find()
}

func testOrderNotJoinedWithLimitAfter() {
	cql.Query[models.Brand](
		db,
	).Descending(
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Limit(1).Find()
}

func testOrderNotJoinedWithLimitBefore() {
	cql.Query[models.Brand](
		db,
	).Limit(1).Descending(
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}
