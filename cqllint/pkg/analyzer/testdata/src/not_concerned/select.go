package not_concerned

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

type Result struct{}

func testQueryMainModelInsideSelect() {
	cql.Select[models.Brand, any](
		cql.Query[models.Brand](
			db,
			conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
		),
	)
}

func testQueryJoinedInsideSelect() {
	cql.Select[models.Phone, any](
		cql.Query[models.Phone](
			db,
			conditions.Phone.Brand(),
			conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		),
	)
}

func testQueryNotJoinedInsideSelect() {
	cql.Select[models.Brand, any](
		cql.Query[models.Brand](
			db,
			conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
			conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	)
}

func testSelectMainModel() {
	cql.Select(
		cql.Query[models.Brand](
			db,
			conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
		),
		cql.ValueInto(conditions.Brand.Name, func(_ string, _ *Result) {}),
	)
}

func testSelectJoinedModel() {
	cql.Select(
		cql.Query[models.Phone](
			db,
			conditions.Phone.Brand(),
			conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		),
		cql.ValueInto(conditions.Brand.Name, func(_ string, _ *Result) {}),
	)
}

func testSelectNotJoinedModel() {
	cql.Select[models.Brand, Result](
		cql.Query[models.Brand](
			db,
			conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
		),
		cql.ValueInto(conditions.City.Name, func(_ string, _ *Result) {}), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}
