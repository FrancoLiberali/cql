package not_concerned

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

// Having
// TODO Verificar que tenga los campos seleccionados? seria otro tipo de test, que esta bueno pero complicado
// TODO groupby repetidos

var results = []struct{}{}

func testGroupBySameModel() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
	).GroupBy(
		conditions.Brand.Name,
	)
}

func testGroupByJoinedModel() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Phone.Name,
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupByJoinedWithJoinedWithCondition() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("asd")),
		),
	).GroupBy(
		conditions.Phone.Name,
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupByJoinedWithJoinedWithPreload() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand().Preload(),
	).GroupBy(
		conditions.Phone.Name,
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupByJoinedWithJoinedWithConditionsWithPreload() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("asd")),
		).Preload(),
	).GroupBy(
		conditions.Phone.Name,
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupByJoinedModelInVariable() {
	value := conditions.Brand.Name

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		value,
	)
}

func testGroupByNotJoined() {
	cql.Query[models.Brand](
		db,
	).GroupBy(
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupByNotJoinedInVariable() {
	value := conditions.City.Name // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		value,
	)
}

func testGroupByJoinedConditionInList() {
	values := []condition.IField{
		conditions.Phone.Name,
	}

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		values...,
	)
}

func testGroupByNotJoinedConditionInList() {
	values := []condition.IField{
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	}

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		values...,
	)
}

func testGroupByJoinedConditionInListWithAppend() {
	values := []condition.IField{}

	values = append(
		values,
		conditions.Phone.Name,
	)

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		values...,
	)
}

func testGroupByNotJoinedConditionInListWithAppend() {
	values := []condition.IField{}

	values = append(
		values,
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		values...,
	)
}

func testGroupByNotJoinedConditionInListWithAppendSecond() {
	values := []condition.IField{}

	values = append(
		values,
		conditions.Phone.Name,
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		values...,
	)
}

func testGroupByNotJoinedConditionInListWithAppendMultiple() {
	values := []condition.IField{}

	values = append(
		values,
		conditions.Phone.Name,
	)

	values = append(
		values,
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		values...,
	)
}

func testSelectSameModel() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
	).GroupBy(
		conditions.Brand.Name,
	).Select(
		conditions.Brand.Name.Aggregate().Max(), "aggregation1",
	)
}

func testSelectJoinedModel() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Phone.Name,
	).Select(
		conditions.Brand.Name.Aggregate().Max(), "aggregation1",
	)
}

func testSelectNotJoined() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
	).GroupBy(
		conditions.Brand.Name,
	).Select(
		conditions.City.Name.Aggregate().Max(), "aggregation1", // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testSelectJoinedModelInVariable() {
	value := conditions.Brand.Name.Aggregate().Max()

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Select(
		value, "aggregation1",
	)
}

func testSelectNotJoinedInVariable() {
	value := conditions.City.Name.Aggregate().Max() // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Select(
		value, "aggregation1",
	)
}

func testSelectJoinedWithFunction() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Select(
		conditions.Brand.Name.Concat("asd").Aggregate().Max(), "aggregation1",
	)
}

func testSelectJoinedWithFunctionVariable() {
	value := conditions.Brand.Name.Concat("asd").Aggregate().Max()

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Select(
		value, "aggregation1",
	)
}

func testSelectJoinedWithFunctionOverVariable() {
	value := conditions.Brand.Name

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Select(
		value.Concat("asd").Aggregate().Max(), "aggregation1",
	)
}

func testSelectNotJoinedWithFunction() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Select(
		conditions.City.Name.Concat("asd").Aggregate().Max(), "aggregation1", // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testSelectNotJoinedWithFunctionVariable() {
	value := conditions.City.Name.Concat("asd").Aggregate().Max() // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Select(
		value, "aggregation1",
	)
}

func testSelectNotJoinedWithFunctionOverVariable() {
	value := conditions.City.Name

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Select(
		value.Concat("asd").Aggregate().Max(), "aggregation1", // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testSelectNotJoinedWithTwoFunctions() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Select(
		conditions.City.Name.Concat("asd").Concat("asd").Aggregate().Max(), "aggregation1", // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testHavingSameModel() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.Brand.Name.Aggregate().Max().Eq(cql.String("asd")),
	)
}

func testHavingJoinedModel() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Phone.Name,
	).Having(
		conditions.Brand.Name.Aggregate().Max().Eq(cql.String("asd")),
	)
}

func testHavingNotJoined() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.City.Name.Aggregate().Max().Eq(cql.String("asd")), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testHavingJoinedModelInVariable() {
	value := conditions.Brand.Name.Aggregate().Max().Eq(cql.String("asd"))

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		value,
	)
}

func testHavingNotJoinedInVariable() {
	value := conditions.City.Name.Aggregate().Max().Eq(cql.String("asd")) // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		value,
	)
}

func testHavingJoinedWithFunction() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.Brand.Name.Concat("asd").Aggregate().Max().Eq(cql.String("asd")),
	)
}

func testHavingJoinedWithFunctionVariable() {
	value := conditions.Brand.Name.Concat("asd").Aggregate().Max().Eq(cql.String("asd"))

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		value,
	)
}

func testHavingJoinedWithFunctionOverVariable() {
	value := conditions.Brand.Name

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		value.Concat("asd").Aggregate().Max().Eq(cql.String("asd")),
	)
}

func testHavingNotJoinedWithFunction() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.City.Name.Concat("asd").Aggregate().Max().Eq(cql.String("asd")), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testHavingNotJoinedWithFunctionVariable() {
	value := conditions.City.Name.Concat("asd").Aggregate().Max().Eq(cql.String("asd")) // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		value,
	)
}

func testHavingNotJoinedWithFunctionOverVariable() {
	value := conditions.City.Name

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		value.Concat("asd").Aggregate().Max().Eq(cql.String("asd")), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testHavingNotJoinedWithTwoFunctions() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.City.Name.Concat("asd").Concat("asd").Aggregate().Max().Eq(cql.String("asd")), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testHavingJoinedConditionInList() {
	values := []condition.AggregationCondition{
		conditions.Brand.Name.Aggregate().Max().Eq(cql.String("asd")),
	}

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		values...,
	)
}

func testHavingNotJoinedConditionInList() {
	values := []condition.AggregationCondition{
		conditions.City.Name.Aggregate().Max().Eq(cql.String("asd")), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	}

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		values...,
	)
}

func testHavingJoinedConditionInListWithAppend() {
	values := []condition.AggregationCondition{}

	values = append(
		values,
		conditions.Brand.Name.Aggregate().Max().Eq(cql.String("asd")),
	)

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		values...,
	)
}

func testHavingNotJoinedConditionInListWithAppend() {
	values := []condition.AggregationCondition{}

	values = append(
		values,
		conditions.City.Name.Aggregate().Max().Eq(cql.String("asd")), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		values...,
	)
}

func testHavingNotJoinedConditionInListWithAppendSecond() {
	values := []condition.AggregationCondition{}

	values = append(
		values,
		conditions.Phone.Name.Aggregate().Max().Eq(cql.String("asd")),
		conditions.City.Name.Aggregate().Max().Eq(cql.String("asd")), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		values...,
	)
}

func testHavingNotJoinedConditionInListWithAppendMultiple() {
	values := []condition.AggregationCondition{}

	values = append(
		values,
		conditions.Phone.Name.Aggregate().Max().Eq(cql.String("asd")),
	)

	values = append(
		values,
		conditions.City.Name.Aggregate().Max().Eq(cql.String("asd")), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		values...,
	)
}

// TODO having:
// del otro lado del eq (y con funcion y en list y en variable (de ambos lados))
