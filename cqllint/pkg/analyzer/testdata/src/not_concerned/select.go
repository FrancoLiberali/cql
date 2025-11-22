package not_concerned

import (
	"context"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

type Result struct {
	AggregationString string
}

func testQueryMainModelInsideSelect() {
	cql.Select[any](
		cql.Query[models.Brand](
			context.Background(),
			db,
			conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
		),
	)
}

func testQueryJoinedInsideSelect() {
	cql.Select[any](
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
			conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		),
	)
}

func testQueryNotJoinedInsideSelect() {
	cql.Select[any](
		cql.Query[models.Brand](
			context.Background(),
			db,
			conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
			conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	)
}

func testSelectMainModel() {
	cql.Select(
		cql.Query[models.Brand](
			context.Background(),
			db,
			conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
		),
		cql.ValueInto(conditions.Brand.Name, func(_ string, _ *Result) {}),
	)
}

func testSelectJoinedModel() {
	cql.Select(
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
			conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		),
		cql.ValueInto(conditions.Brand.Name, func(_ string, _ *Result) {}),
	)
}

func testSelectNotJoinedModel() {
	cql.Select(
		cql.Query[models.Brand](
			context.Background(),
			db,
			conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
		),
		cql.ValueInto(conditions.City.Name, func(_ string, _ *Result) {}), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testSelectNotJoinedModelSecond() {
	cql.Select[Result](
		cql.Query[models.Brand](
			context.Background(),
			db,
			conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
		),
		cql.ValueInto(conditions.Brand.Name, func(_ string, _ *Result) {}),
		cql.ValueInto(conditions.City.Name, func(_ string, _ *Result) {}), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testSelectJoinedModelWithFunction() {
	cql.Select(
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
			conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		),
		cql.ValueInto(conditions.Brand.Name.Concat(cql.String("asd")), func(_ string, _ *Result) {}),
	)
}

func testSelectNotJoinedModelWithFunction() {
	cql.Select[Result](
		cql.Query[models.Brand](
			context.Background(),
			db,
			conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
		),
		cql.ValueInto(conditions.City.Name.Concat(cql.String("asd")), func(_ string, _ *Result) {}), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testSelectJoinedModelInVar() {
	value := conditions.Brand.Name

	cql.Select(
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
			conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		),
		cql.ValueInto(value, func(_ string, _ *Result) {}),
	)
}

func testSelectNotJoinedModelInVar() {
	value := conditions.City.Name

	cql.Select[Result](
		cql.Query[models.Brand](
			context.Background(),
			db,
			conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
		),
		cql.ValueInto(value, func(_ string, _ *Result) {}), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testSelectJoinedModelInListInVar() {
	selects := []condition.Selection[Result]{
		cql.ValueInto(conditions.Brand.Name, func(_ string, _ *Result) {}),
	}

	cql.Select(
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
			conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		),
		selects...,
	)
}

func testSelectNotJoinedModelInListInVar() {
	selects := []condition.Selection[Result]{
		cql.ValueInto(conditions.City.Name, func(_ string, _ *Result) {}), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	}

	cql.Select(
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
			conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		),
		selects...,
	)
}

func testSelectNotJoinedModelInListWithAppend() {
	selects := []condition.Selection[Result]{}

	selects = append(
		selects,
		cql.ValueInto(conditions.City.Name, func(_ string, _ *Result) {}), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Select(
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
			conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		),
		selects...,
	)
}
