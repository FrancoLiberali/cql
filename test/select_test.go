package test

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

type SelectIntTestSuite struct {
	testSuite
}

func NewSelectIntTestSuite(
	db *gorm.DB,
) *SelectIntTestSuite {
	return &SelectIntTestSuite{
		testSuite: testSuite{
			db: db,
		},
	}
}

// TODO hacer lo mismo para los selects del groupby

func (ts *SelectIntTestSuite) TestOneSelect() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 0, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			ts.db,
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *ResultInt) {
			result.Int = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{
		{Int: 0},
		{Int: 1},
		{Int: 1},
	}, results)
}

func (ts *SelectIntTestSuite) TestTwoSelectSameValue() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 0, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			ts.db,
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *ResultInt) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *ResultInt) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{
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
			ts.db,
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *ResultInt) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Float, func(value float64, result *ResultInt) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{
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
			ts.db,
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *ResultInt) {
			result.Int = int(value) + 1
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{
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
			ts.db,
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int.Plus(1), func(value float64, result *ResultInt) {
			result.Int = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{
		{Int: 1},
		{Int: 2},
		{Int: 3},
	}, results)
}

func (ts *SelectIntTestSuite) TestSelectMultipleWithFunction() {
	ts.createProduct("1", 0, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 2, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			ts.db,
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *ResultIntAndFloat) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Plus(1), func(value float64, result *ResultIntAndFloat) {
			result.Aggregation1 = int(value)
		}),
		cql.ValueInto(conditions.Product.Float.Minus(1.5), func(value float64, result *ResultIntAndFloat) {
			result.Aggregation2 = value
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultIntAndFloat{
		{Int: 0, Aggregation1: 1, Aggregation2: -1.5},
		{Int: 1, Aggregation1: 2, Aggregation2: -0.5},
		{Int: 2, Aggregation1: 3, Aggregation2: 0.5},
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
			ts.db,
			conditions.Sale.Product(),
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *ResultInt) {
			result.Int = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{
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
			ts.db,
			conditions.Sale.Product(),
		).Descending(conditions.Sale.Code),
		cql.ValueInto(conditions.Sale.Code, func(value float64, result *ResultInt) {
			result.Int = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{
		{Int: 1},
		{Int: 1},
		{Int: 1},
		{Int: 2},
		{Int: 2},
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
			ts.db,
		).Descending(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *ResultInt) {
			result.Int = int(value)
		}),
	)

	ts.ErrorIs(err, cql.ErrFieldModelNotConcerned)
	ts.ErrorContains(err, "field's model is not concerned by the query (not joined); not concerned model: models.Product")
}
