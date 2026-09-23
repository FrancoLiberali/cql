package not_concerned

import (
	"context"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

var results = []struct{}{}

func testGroupBySameModel() {
	cql.Query[models.Brand](
		context.Background(),
		db,
		conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
	).GroupBy(
		conditions.Brand.Name,
	)
}

func testGroupByJoinedModel() {
	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Phone.Name,
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupByJoinedWithJoinedWithCondition() {
	cql.Query[models.Phone](
		context.Background(),
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
		context.Background(),
		db,
		conditions.Phone.Brand().Preload(),
	).GroupBy(
		conditions.Phone.Name,
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupByJoinedWithJoinedWithConditionsWithPreload() {
	cql.Query[models.Phone](
		context.Background(),
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
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		value,
	)
}

func testGroupByNotJoined() {
	cql.Query[models.Brand](
		context.Background(),
		db,
	).GroupBy(
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupByNotJoinedInVariable() {
	value := conditions.City.Name // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"

	cql.Query[models.Phone](
		context.Background(),
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
		context.Background(),
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
		context.Background(),
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
		context.Background(),
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
		context.Background(),
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
		context.Background(),
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
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		values...,
	)
}

func testGroupBySelectSameModel() {
	cql.Select(
		cql.Query[models.Brand](
			context.Background(),
			db,
			conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
		).GroupBy(
			conditions.Brand.Name,
		),
		conditions.Brand.Name.Aggregate().Max().Into(func(result *Result) *string { return &result.AggregationString }),
	)
}

func testGroupBySelectJoinedModel() {
	cql.Select(
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
		).GroupBy(
			conditions.Phone.Name,
		),
		conditions.Brand.Name.Aggregate().Max().Into(func(result *Result) *string { return &result.AggregationString }),
	)
}

func testGroupBySelectNotJoined() {
	cql.Select(
		cql.Query[models.Brand](
			context.Background(),
			db,
			conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
		).GroupBy(
			conditions.Brand.Name,
		),
		conditions.City.Name.Aggregate().Max().Into(func(result *Result) *string { return &result.AggregationString }), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupBySelectJoinedModelInVariable() {
	value := conditions.Brand.Name.Aggregate().Max()

	cql.Select(
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
		).GroupBy(
			conditions.Brand.Name,
		),
		value.Into(func(result *Result) *string { return &result.AggregationString }),
	)
}

func testGroupBySelectNotJoinedInVariable() {
	field := conditions.City.Name

	cql.Select(
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
		).GroupBy(
			conditions.Brand.Name,
		),
		field.Aggregate().Max().Into(func(result *Result) *string { return &result.AggregationString }), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupBySelectJoinedWithFunction() {
	cql.Select(
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
		).GroupBy(
			conditions.Brand.Name,
		),
		conditions.Brand.Name.Concat(cql.String("asd")).Aggregate().Max().Into(func(result *Result) *string { return &result.AggregationString }),
	)
}

func testGroupBySelectJoinedWithFunctionVariable() {
	value := conditions.Brand.Name.Concat(cql.String("asd")).Aggregate().Max()

	cql.Select(
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
		).GroupBy(
			conditions.Brand.Name,
		),
		value.Into(func(result *Result) *string { return &result.AggregationString }),
	)
}

func testGroupBySelectJoinedWithFunctionOverVariable() {
	value := conditions.Brand.Name

	cql.Select(
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
		).GroupBy(
			conditions.Brand.Name,
		),
		value.Concat(cql.String("asd")).Aggregate().Max().Into(func(result *Result) *string { return &result.AggregationString }),
	)
}

func testGroupBySelectNotJoinedWithFunction() {
	cql.Select(
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
		).GroupBy(
			conditions.Brand.Name,
		),
		conditions.City.Name.Concat(cql.String("asd")).Aggregate().Max().Into(func(result *Result) *string { return &result.AggregationString }), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupBySelectNotJoinedWithFunctionVariable() {
	value := conditions.City.Name

	cql.Select(
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
		).GroupBy(
			conditions.Brand.Name,
		),
		value.Concat(cql.String("asd")).Aggregate().Max().Into(func(result *Result) *string { return &result.AggregationString }), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupBySelectNotJoinedWithFunctionOverVariable() {
	value := conditions.City.Name

	cql.Select(
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
		).GroupBy(
			conditions.Brand.Name,
		),
		value.Concat(cql.String("asd")).Aggregate().Max().Into(func(result *Result) *string { return &result.AggregationString }), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupBySelectNotJoinedWithTwoFunctions() {
	cql.Select(
		cql.Query[models.Phone](
			context.Background(),
			db,
			conditions.Phone.Brand(),
		).GroupBy(
			conditions.Brand.Name,
		),
		conditions.City.Name.Concat(cql.String("asd")).Concat(cql.String("asd")).Aggregate().Max().Into(func(result *Result) *string { return &result.AggregationString }), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testHavingSameModel() {
	cql.Query[models.Brand](
		context.Background(),
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
		context.Background(),
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
		context.Background(),
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
		context.Background(),
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
		context.Background(),
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
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.Brand.Name.Concat(cql.String("asd")).Aggregate().Max().Eq(cql.String("asd")),
	)
}

func testHavingJoinedWithFunctionVariable() {
	value := conditions.Brand.Name.Concat(cql.String("asd")).Aggregate().Max().Eq(cql.String("asd"))

	cql.Query[models.Phone](
		context.Background(),
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
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		value.Concat(cql.String("asd")).Aggregate().Max().Eq(cql.String("asd")),
	)
}

func testHavingNotJoinedWithFunction() {
	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.City.Name.Concat(cql.String("asd")).Aggregate().Max().Eq(cql.String("asd")), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testHavingNotJoinedWithFunctionVariable() {
	value := conditions.City.Name.Concat(cql.String("asd")).Aggregate().Max().Eq(cql.String("asd")) // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"

	cql.Query[models.Phone](
		context.Background(),
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
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		value.Concat(cql.String("asd")).Aggregate().Max().Eq(cql.String("asd")), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testHavingNotJoinedWithTwoFunctions() {
	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.City.Name.Concat(cql.String("asd")).Concat(cql.String("asd")).Aggregate().Max().Eq(cql.String("asd")), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testHavingJoinedConditionInList() {
	values := []condition.AggregationCondition{
		conditions.Brand.Name.Aggregate().Max().Eq(cql.String("asd")),
	}

	cql.Query[models.Phone](
		context.Background(),
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
		context.Background(),
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
		context.Background(),
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
		context.Background(),
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
		context.Background(),
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
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		values...,
	)
}

func testHavingDynamicSameModel() {
	cql.Query[models.Brand](
		context.Background(),
		db,
		conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.Brand.Name.Aggregate().Max().Eq(conditions.Brand.Name.Aggregate().Max()),
	)
}

func testHavingDynamicJoinedModel() {
	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Phone.Name,
	).Having(
		conditions.Brand.Name.Aggregate().Max().Eq(conditions.Phone.Name.Aggregate().Max()),
	)
}

func testHavingDynamicNotJoined() {
	cql.Query[models.Brand](
		context.Background(),
		db,
		conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.Brand.Name.Aggregate().Max().Eq(conditions.City.Name.Aggregate().Max()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testHavingDynamicNotJoinedOnLeft() {
	cql.Query[models.Brand](
		context.Background(),
		db,
		conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.City.Name.Aggregate().Max().Eq(conditions.Brand.Name.Aggregate().Max()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testHavingDynamicNotJoinedOnBoth() {
	cql.Query[models.Brand](
		context.Background(),
		db,
		conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.City.Name.Aggregate().Max().Eq( // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
			conditions.Product.String.Aggregate().Max(), // want "github.com/FrancoLiberali/cql/test/models.Product is not joined by the query"
		),
	)
}

func testHavingDynamicJoinedModelInVariable() {
	value := conditions.Brand.Name.Aggregate().Max().Eq(conditions.Brand.Name.Aggregate().Max())

	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		value,
	)
}

func testHavingDynamicNotJoinedInVariable() {
	value := conditions.Brand.Name.Aggregate().Max().Eq(conditions.City.Name.Aggregate().Max()) // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"

	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		value,
	)
}

func testHavingDynamicJoinedWithFunctionOnLeft() {
	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.Brand.Name.Concat(cql.String("asd")).Aggregate().Max().Eq(conditions.Brand.Name.Aggregate().Max()),
	)
}

func testHavingDynamicJoinedWithFunctionOnRight() {
	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.Brand.Name.Aggregate().Max().Eq(conditions.Brand.Name.Concat(cql.String("asd")).Aggregate().Max()),
	)
}

func testHavingDynamicJoinedWithFunctionVariable() {
	value := conditions.Brand.Name.Aggregate().Max().Eq(conditions.Brand.Name.Concat(cql.String("asd")).Aggregate().Max())

	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		value,
	)
}

func testHavingDynamicJoinedWithFunctionOverVariable() {
	value := conditions.Brand.Name

	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.Brand.Name.Aggregate().Max().Eq(value.Concat(cql.String("asd")).Aggregate().Max()),
	)
}

func testHavingDynamicNotJoinedWithFunction() {
	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.Brand.Name.Aggregate().Max().Eq(conditions.City.Name.Concat(cql.String("asd")).Aggregate().Max()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testHavingDynamicNotJoinedWithFunctionVariable() {
	value := conditions.Brand.Name.Aggregate().Max().Eq(conditions.City.Name.Concat(cql.String("asd")).Aggregate().Max()) // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"

	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		value,
	)
}

func testHavingDynamicNotJoinedWithFunctionOverVariable() {
	value := conditions.City.Name

	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		conditions.Brand.Name.Aggregate().Max().Eq(value.Concat(cql.String("asd")).Aggregate().Max()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testHavingDynamicJoinedConditionInList() {
	values := []condition.AggregationCondition{
		conditions.Brand.Name.Aggregate().Max().Eq(conditions.Brand.Name.Concat(cql.String("asd")).Aggregate().Max()),
	}

	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		values...,
	)
}

func testHavingDynamicNotJoinedConditionInList() {
	values := []condition.AggregationCondition{
		conditions.Brand.Name.Aggregate().Max().Eq(conditions.City.Name.Concat(cql.String("asd")).Aggregate().Max()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	}

	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		values...,
	)
}

func testHavingDynamicJoinedConditionInListWithAppend() {
	values := []condition.AggregationCondition{}

	values = append(
		values,
		conditions.Brand.Name.Aggregate().Max().Eq(conditions.Brand.Name.Concat(cql.String("asd")).Aggregate().Max()),
	)

	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		values...,
	)
}

func testHavingDynamicNotJoinedConditionInListWithAppend() {
	values := []condition.AggregationCondition{}

	values = append(
		values,
		conditions.Brand.Name.Aggregate().Max().Eq(conditions.City.Name.Concat(cql.String("asd")).Aggregate().Max()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		values...,
	)
}

func testHavingDynamicNotJoinedConditionInListWithAppendSecond() {
	values := []condition.AggregationCondition{}

	values = append(
		values,
		conditions.Brand.Name.Aggregate().Max().Eq(conditions.Brand.Name.Concat(cql.String("asd")).Aggregate().Max()),
		conditions.Brand.Name.Aggregate().Max().Eq(conditions.City.Name.Concat(cql.String("asd")).Aggregate().Max()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		values...,
	)
}

func testHavingDynamicNotJoinedConditionInListWithAppendMultiple() {
	values := []condition.AggregationCondition{}

	values = append(
		values,
		conditions.Brand.Name.Aggregate().Max().Eq(conditions.Brand.Name.Concat(cql.String("asd")).Aggregate().Max()),
	)

	values = append(
		values,
		conditions.Brand.Name.Aggregate().Max().Eq(conditions.City.Name.Concat(cql.String("asd")).Aggregate().Max()), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Query[models.Phone](
		context.Background(),
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Brand.Name,
	).Having(
		values...,
	)
}
