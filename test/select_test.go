package test

import (
	"context"
	"reflect"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/sql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

type SelectIntTestSuite struct {
	testSuite
}

func NewSelectIntTestSuite(
	db *cql.DB,
) *SelectIntTestSuite {
	return &SelectIntTestSuite{
		testSuite: testSuite{
			db: db,
		},
	}
}

func (ts *SelectIntTestSuite) TestSelectOneSelect() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 0, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 0},
		{Int: 1},
		{Int: 1},
	}, results)
}

func (ts *SelectIntTestSuite) TestSelectWithOrder() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 0, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 0},
		{Int: 1},
		{Int: 1},
	}, results)
}

func (ts *SelectIntTestSuite) TestSelectWithMultipleOrder() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 1, true, nil)
	ts.createProduct("5", 0, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).Descending(conditions.Product.Int).Descending(conditions.Product.Bool),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
	)

	ts.Require().NoError(err)
	ts.Assert().True(reflect.DeepEqual(results, []Result{
		{Int: 1},
		{Int: 1},
		{Int: 0},
	}))
}

func (ts *SelectIntTestSuite) TestSelectWithOrderNotSelected() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("5", 0, 2, true, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).Descending(conditions.Product.Bool),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
	)

	ts.Require().NoError(err)
	ts.Assert().True(reflect.DeepEqual(results, []Result{
		{Int: 0},
		{Int: 1},
	}))
}

func (ts *SelectIntTestSuite) TestTwoSelectSameValue() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 0, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 0, Aggregation1: 0},
		{Int: 1, Aggregation1: 1},
		{Int: 1, Aggregation1: 1},
	}, results)
}

func (ts *SelectIntTestSuite) TestTwoSelectDifferentValue() {
	ts.createProduct("1", 2, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 0, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Float, func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 0, Aggregation1: 2},
		{Int: 1, Aggregation1: 1},
		{Int: 2, Aggregation1: 0},
	}, results)
}

func (ts *SelectIntTestSuite) TestOneSelectWithFunctionInGo() {
	ts.createProduct("1", 0, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 2, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value) + 1
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 1},
		{Int: 2},
		{Int: 3},
	}, results)
}

func (ts *SelectIntTestSuite) TestOneSelectWithFunctionInCQL() {
	ts.createProduct("1", 0, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 2, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int.Plus(cql.Int(1)), func(value float64, result *Result) {
			result.Int = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 1},
		{Int: 2},
		{Int: 3},
	}, results)
}

func (ts *SelectIntTestSuite) TestOneSelectWithFunctionDynamic() {
	ts.createProduct("1", 0, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 2, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int.Plus(conditions.Product.Float), func(value float64, result *Result) {
			result.Int = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 0},
		{Int: 2},
		{Int: 4},
	}, results)
}

func (ts *SelectIntTestSuite) TestSelectMultipleWithFunction() {
	ts.createProduct("1", 0, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 2, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Plus(cql.Int(1)), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
		cql.ValueInto(conditions.Product.Float.Minus(cql.Float64(1.5)), func(value float64, result *Result) {
			result.Aggregation3 = value
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 0, Aggregation1: 1, Aggregation3: -1.5},
		{Int: 1, Aggregation1: 2, Aggregation3: -0.5},
		{Int: 2, Aggregation1: 3, Aggregation3: 0.5},
	}, results)
}

func (ts *SelectIntTestSuite) TestSelectFromJoinedModel() {
	product0 := ts.createProduct("1", 0, 0, false, nil)
	product1 := ts.createProduct("2", 1, 1, false, nil)
	product2 := ts.createProduct("5", 2, 2, false, nil)

	ts.createSale(1, product0, nil)
	ts.createSale(1, product1, nil)
	ts.createSale(1, product2, nil)
	ts.createSale(2, product1, nil)
	ts.createSale(2, product2, nil)

	results, err := cql.Select(
		cql.Query[models.Sale](
			context.Background(),
			ts.db,
			conditions.Sale.Product(),
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 0},
		{Int: 1},
		{Int: 1},
		{Int: 2},
		{Int: 2},
	}, results)
}

func (ts *SelectIntTestSuite) TestSelectFromMainModelAfterJoin() {
	product0 := ts.createProduct("1", 0, 0, false, nil)
	product1 := ts.createProduct("2", 1, 1, false, nil)
	product2 := ts.createProduct("5", 2, 2, false, nil)

	ts.createSale(1, product0, nil)
	ts.createSale(1, product1, nil)
	ts.createSale(1, product2, nil)
	ts.createSale(2, product1, nil)
	ts.createSale(2, product2, nil)

	results, err := cql.Select(
		cql.Query[models.Sale](
			context.Background(),
			ts.db,
			conditions.Sale.Product(),
		).Descending(conditions.Sale.Code),
		cql.ValueInto(conditions.Sale.Code, func(value float64, result *Result) {
			result.Int = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 1},
		{Int: 1},
		{Int: 1},
		{Int: 2},
		{Int: 2},
	}, results)
}

func (ts *SelectIntTestSuite) TestSelectFromMainAndJoinedModel() {
	product0 := ts.createProduct("1", 0, 0, false, nil)
	product1 := ts.createProduct("2", 1, 1, false, nil)
	product2 := ts.createProduct("5", 2, 2, false, nil)

	ts.createSale(1, product0, nil)
	ts.createSale(1, product1, nil)
	ts.createSale(1, product2, nil)
	ts.createSale(2, product1, nil)
	ts.createSale(2, product2, nil)

	results, err := cql.Select(
		cql.Query[models.Sale](
			context.Background(),
			ts.db,
			conditions.Sale.Product(),
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Sale.Code, func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 0, Aggregation1: 1},
		{Int: 1, Aggregation1: 1},
		{Int: 1, Aggregation1: 1},
		{Int: 2, Aggregation1: 2},
		{Int: 2, Aggregation1: 2},
	}, results)
}

func (ts *SelectIntTestSuite) TestSelectFromNotJoinedModelReturnsError() {
	product0 := ts.createProduct("1", 0, 0, false, nil)
	product1 := ts.createProduct("2", 1, 1, false, nil)
	product2 := ts.createProduct("5", 2, 2, false, nil)

	ts.createSale(1, product0, nil)
	ts.createSale(1, product1, nil)
	ts.createSale(1, product2, nil)
	ts.createSale(2, product1, nil)
	ts.createSale(2, product2, nil)

	_, err := cql.Select(
		cql.Query[models.Sale](
			context.Background(),
			ts.db,
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
	)

	ts.ErrorIs(err, cql.ErrFieldModelNotConcerned)
	ts.ErrorContains(err, "field's model is not concerned by the query (not joined); not concerned model: models.Product")
}

// TestSelectIntoPointerFormStrict exercises the Go 1.27 pointer-form .Into:
//   - a plain int column binds straight into an int field (no float64 detour,
//     no int(value) conversion at the call site);
//   - a string column binds into a string field;
//   - an aggregation binds into its honest float64 field.
func (ts *SelectIntTestSuite) TestSelectIntoPointerFormStrict() {
	ts.createProduct("a", 2, 0, false, nil)
	ts.createProduct("b", 1, 0, false, nil)
	ts.createProduct("c", 3, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).Ascending(conditions.Product.Int),
		conditions.Product.Int.Into(func(r *Result) *int { return &r.Int }),
		conditions.Product.String.Into(func(r *Result) *string { return &r.String }),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 1, String: "b"},
		{Int: 2, String: "a"},
		{Int: 3, String: "c"},
	}, results)
}

func (ts *SelectIntTestSuite) TestSelectIntoPointerFormAggregation() {
	ts.createProduct("1", 4, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("5", 1, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		),
		conditions.Product.Int.Aggregate().Sum().Into(func(r *Result) *float64 { return &r.Aggregation3 }),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Aggregation3: 6},
	}, results)
}

// A numeric arithmetic expression binds *float64 (SQL widens); covers the
// NotUpdatableNumericField.Into path.
func (ts *SelectIntTestSuite) TestSelectIntoPointerFormNumericExpression() {
	ts.createProduct("1", 10, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		),
		conditions.Product.Int.Plus(cql.Int(1)).Into(func(r *Result) *float64 { return &r.Aggregation3 }),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Aggregation3: 11},
	}, results)
}

func (ts *SelectIntTestSuite) TestSelectAggregations() {
	ts.createProduct("1", 4, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 1, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		),
		cql.ValueInto(conditions.Product.Int.Aggregate().Max(), func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Min(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Count(), func(value float64, result *Result) {
			result.Aggregation2 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 4, Aggregation1: 1, Aggregation2: 3},
	}, results)
}

func (ts *SelectIntTestSuite) TestSelectAggregationsAndNotAggregations() {
	// mix of aggregations and not aggregations only supported by sqlite
	if getDBDialector() == sql.SQLite {
		ts.createProduct("1", 4, 0, false, nil)
		ts.createProduct("2", 1, 1, false, nil)
		ts.createProduct("5", 1, 2, false, nil)

		results, err := cql.Select(
			cql.Query[models.Product](
				context.Background(),
				ts.db,
			).Descending(conditions.Product.Int),
			cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
				result.Int = int(value)
			}),
			cql.ValueInto(conditions.Product.Int.Aggregate().Min(), func(value float64, result *Result) {
				result.Aggregation1 = int(value)
			}),
			cql.ValueInto(conditions.Product.Int.Aggregate().Count(), func(value float64, result *Result) {
				result.Aggregation2 = int(value)
			}),
		)

		ts.Require().NoError(err)
		EqualList(&ts.Suite, []Result{
			{Int: 1, Aggregation1: 1, Aggregation2: 3},
		}, results)
	}
}
