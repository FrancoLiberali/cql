package test

import (
	"context"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/sql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

type GroupByIntTestSuite struct {
	testSuite
}

func NewGroupByIntTestSuite(
	db *cql.DB,
) *GroupByIntTestSuite {
	return &GroupByIntTestSuite{
		testSuite: testSuite{
			db: db,
		},
	}
}

type Result struct {
	Int          int
	Float        float64
	String       string
	Aggregation1 int
	Aggregation2 int
	Aggregation3 float64
	Aggregation4 bool
	Aggregation5 string
}

func (ts *GroupByIntTestSuite) TestGroupByNoSelect() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1}, {Int: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByFieldNotPresentReturnsError() {
	_, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(conditions.Sale.SellerID),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
	)

	ts.ErrorIs(err, cql.ErrFieldModelNotConcerned)
	ts.ErrorContains(err, "field's model is not concerned by the query (not joined); not concerned model: models.Sale")
}

func (ts *GroupByIntTestSuite) TestGroupByWithConditionsNoSelect() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
			conditions.Product.Int.Is().Eq(cql.Int(1)),
		).GroupBy(conditions.Product.Int),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectFieldNotPresentReturnsError() {
	_, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Sale.ID.Aggregate().Count(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.ErrorIs(err, cql.ErrFieldModelNotConcerned)
	ts.ErrorContains(err, "field's model is not concerned by the query (not joined); not concerned model: models.Sale")
}

func (ts *GroupByIntTestSuite) TestGroupBySelectSum() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2}, {Int: 0, Aggregation1: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectCount() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Count(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2}, {Int: 0, Aggregation1: 1}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectCountWithNulls() {
	int1 := 1

	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, &int1)
	ts.createProduct("3", 0, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.IntPointer.Aggregate().Count(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 1}, {Int: 0, Aggregation1: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectCountAll() {
	int1 := 1

	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, &int1)
	ts.createProduct("3", 0, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(cql.CountAll(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2}, {Int: 0, Aggregation1: 1}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectAverage() {
	ts.createProduct("1", 1, 0.25, false, nil)
	ts.createProduct("2", 1, 0.75, false, nil)
	ts.createProduct("3", 0, 1, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Float.Aggregate().Average(), func(value float64, result *Result) {
			result.Aggregation3 = value
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation3: 0.5}, {Int: 0, Aggregation3: 1}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectMin() {
	ts.createProduct("1", 1, 0.25, false, nil)
	ts.createProduct("2", 1, 0.75, false, nil)
	ts.createProduct("3", 0, 1, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Float.Aggregate().Min(), func(value float64, result *Result) {
			result.Aggregation3 = value
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation3: 0.25}, {Int: 0, Aggregation3: 1}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectMinForNotNumericField() {
	ts.createProduct("1", 1, 0.25, false, nil)
	ts.createProduct("2", 1, 0.75, false, nil)
	ts.createProduct("3", 0, 1, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.String.Aggregate().Min(), func(value string, result *Result) {
			result.Aggregation5 = value
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation5: "1"}, {Int: 0, Aggregation5: "3"}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectMax() {
	ts.createProduct("1", 1, 0.25, false, nil)
	ts.createProduct("2", 1, 0.75, false, nil)
	ts.createProduct("3", 0, 1, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Float.Aggregate().Max(), func(value float64, result *Result) {
			result.Aggregation3 = value
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation3: 0.75}, {Int: 0, Aggregation3: 1}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectMaxForNotNumericField() {
	ts.createProduct("1", 1, 0.25, false, nil)
	ts.createProduct("2", 1, 0.75, false, nil)
	ts.createProduct("3", 0, 1, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.String.Aggregate().Max(), func(value string, result *Result) {
			result.Aggregation5 = value
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation5: "2"}, {Int: 0, Aggregation5: "3"}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectAll() {
	ts.createProduct("1", 1, 0.25, true, nil)
	ts.createProduct("2", 1, 0.75, false, nil)
	ts.createProduct("1", 2, 0.25, true, nil)
	ts.createProduct("2", 2, 0.75, true, nil)
	ts.createProduct("3", 0, 1, true, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Bool.Aggregate().All(), func(value bool, result *Result) {
			result.Aggregation4 = value
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation4: false}, {Int: 2, Aggregation4: true}, {Int: 0, Aggregation4: true}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectAny() {
	ts.createProduct("1", 1, 0.25, true, nil)
	ts.createProduct("2", 1, 0.75, false, nil)
	ts.createProduct("1", 2, 0.25, false, nil)
	ts.createProduct("2", 2, 0.75, false, nil)
	ts.createProduct("3", 0, 1, true, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Bool.Aggregate().Any(), func(value bool, result *Result) {
			result.Aggregation4 = value
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation4: true}, {Int: 2, Aggregation4: false}, {Int: 0, Aggregation4: true}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectNone() {
	ts.createProduct("1", 1, 0.25, true, nil)
	ts.createProduct("2", 1, 0.75, false, nil)
	ts.createProduct("1", 2, 0.25, false, nil)
	ts.createProduct("2", 2, 0.75, false, nil)
	ts.createProduct("3", 0, 1, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Bool.Aggregate().None(), func(value bool, result *Result) {
			result.Aggregation4 = value
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation4: false}, {Int: 2, Aggregation4: true}, {Int: 0, Aggregation4: true}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectAnd() {
	int1 := 1
	int2 := 3
	int3 := 3

	ts.createProduct("1", 1, 1, false, &int1)
	ts.createProduct("2", 1, 3, false, &int2)
	ts.createProduct("1", 2, 3, false, &int3)
	ts.createProduct("3", 3, 1, false, &int1)
	ts.createProduct("3", 3, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.IntPointer.Aggregate().And(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	switch getDBDialector() {
	case sql.Postgres, sql.MySQL:
		ts.Require().NoError(err)
		EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 1}, {Int: 2, Aggregation1: 3}, {Int: 3, Aggregation1: 1}}, results)
	case sql.SQLite, sql.SQLServer:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "function: And")
	}
}

func (ts *GroupByIntTestSuite) TestGroupBySelectOr() {
	int1 := 1
	int2 := 2
	int3 := 3

	ts.createProduct("1", 1, 0.25, false, &int1)
	ts.createProduct("2", 1, 0.75, false, &int2)
	ts.createProduct("1", 2, 0.25, false, &int3)
	ts.createProduct("3", 3, 1, false, &int1)
	ts.createProduct("3", 3, 1, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.IntPointer.Aggregate().Or(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	switch getDBDialector() {
	case sql.Postgres, sql.MySQL:
		ts.Require().NoError(err)
		EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 3}, {Int: 2, Aggregation1: 3}, {Int: 3, Aggregation1: 1}}, results)
	case sql.SQLite, sql.SQLServer:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "function: Or")
	}
}

func (ts *GroupByIntTestSuite) TestGroupBySelectMoreThanOne() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Count(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation2 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2, Aggregation2: 2}, {Int: 0, Aggregation1: 1, Aggregation2: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByMoreThanOne() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("3", 1, 1.1, false, nil)
	ts.createProduct("4", 1, 1.1, false, nil)
	ts.createProduct("5", 0, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
			conditions.Product.Float,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Float, func(value float64, result *Result) {
			result.Float = value
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Count(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 1, Float: 0, Aggregation1: 1},
		{Int: 1, Float: 1, Aggregation1: 1},
		{Int: 1, Float: 1.1, Aggregation1: 2},
		{Int: 0, Float: 0, Aggregation1: 1},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByMoreThanOneSelectMoreThanOne() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("3", 1, 1.1, false, nil)
	ts.createProduct("4", 1, 1.1, false, nil)
	ts.createProduct("5", 0, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int, conditions.Product.Float,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Float, func(value float64, result *Result) {
			result.Float = value
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Count(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
		cql.ValueInto(conditions.Product.Float.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation3 = value
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 1, Float: 0, Aggregation1: 1, Aggregation3: 0},
		{Int: 1, Float: 1, Aggregation1: 1, Aggregation3: 1},
		{Int: 1, Float: 1.1, Aggregation1: 2, Aggregation3: 2.2},
		{Int: 0, Float: 0, Aggregation1: 1, Aggregation3: 0},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByJoinedField() {
	product1 := ts.createProduct("1", 2, 0, false, nil)
	product2 := ts.createProduct("2", 3, 1, false, nil)

	ts.createSale(1, product1, nil)
	ts.createSale(1, product2, nil)
	ts.createSale(2, product1, nil)

	switch getDBDialector() {
	// TODO group by joined field doesn't work for Postgres by bug in gorm
	case sql.MySQL, sql.SQLServer, sql.SQLite:
		results, err := cql.Select(
			cql.Query[models.Sale](
				context.Background(),
				ts.db,
				conditions.Sale.Product(),
			).GroupBy(
				conditions.Product.Int,
			),
			cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
				result.Int = int(value)
			}),
			cql.ValueInto(conditions.Sale.Code.Aggregate().Sum(), func(value float64, result *Result) {
				result.Aggregation1 = int(value)
			}),
		)

		ts.Require().NoError(err)
		EqualList(&ts.Suite, []Result{
			{Int: 2, Aggregation1: 3},
			{Int: 3, Aggregation1: 1},
		}, results)
	}
}

func (ts *GroupByIntTestSuite) TestGroupByWithJoinedFieldInSelect() {
	product1 := ts.createProduct("1", 2, 0, false, nil)
	product2 := ts.createProduct("2", 3, 1, false, nil)

	ts.createSale(1, product1, nil)
	ts.createSale(1, product2, nil)
	ts.createSale(2, product1, nil)

	results, err := cql.Select(
		cql.Query[models.Sale](
			context.Background(),
			ts.db,
			conditions.Sale.Product(),
		).GroupBy(
			conditions.Sale.Code,
		),
		cql.ValueInto(conditions.Sale.Code, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 1, Aggregation1: 5},
		{Int: 2, Aggregation1: 2},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByJoinedFieldAndWithJoinedFieldInSelect() {
	product1 := ts.createProduct("1", 2, 1, false, nil)
	product2 := ts.createProduct("2", 3, 3, false, nil)

	ts.createSale(1, product1, nil)
	ts.createSale(1, product2, nil)
	ts.createSale(2, product1, nil)

	switch getDBDialector() {
	// TODO group by joined field doesn't work for Postgres by bug in gorm
	case sql.MySQL, sql.SQLServer, sql.SQLite:
		results, err := cql.Select(
			cql.Query[models.Sale](
				context.Background(),
				ts.db,
				conditions.Sale.Product(),
			).GroupBy(
				conditions.Product.Int,
			),
			cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
				result.Int = int(value)
			}),
			cql.ValueInto(conditions.Product.Float.Aggregate().Sum(), func(value float64, result *Result) {
				result.Aggregation3 = value
			}),
		)

		ts.Require().NoError(err)
		EqualList(&ts.Suite, []Result{
			{Int: 2, Aggregation3: 2},
			{Int: 3, Aggregation3: 3},
		}, results)
	}
}

func (ts *GroupByIntTestSuite) TestGroupByFieldPresentInMultipleTables() {
	company := ts.createCompany("name1")
	ts.createSeller("name1", company)
	ts.createSeller("name2", company)

	results, err := cql.Select(
		cql.Query[models.Seller](
			context.Background(),
			ts.db,
			conditions.Seller.Company(),
		).GroupBy(
			conditions.Seller.Name,
		),
		cql.ValueInto(conditions.Seller.Name, func(value string, result *Result) {
			result.String = value
		}),
		cql.ValueInto(conditions.Company.Name.Aggregate().Count(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{String: "name1", Aggregation1: 1},
		{String: "name2", Aggregation1: 1},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByJoinedMultipleTimesFieldReturnsError() {
	_, err := cql.Select(
		cql.Query[models.Child](
			context.Background(),
			ts.db,
			conditions.Child.Parent1(
				conditions.Parent1.ParentParent(),
			),
			conditions.Child.Parent2(
				conditions.Parent2.ParentParent(),
			),
		).GroupBy(
			conditions.ParentParent.Name,
		),
		cql.ValueInto(conditions.ParentParent.Name, func(value string, result *Result) {
			result.String = value
		}),
	)

	ts.ErrorIs(err, cql.ErrAppearanceMustBeSelected)
	ts.ErrorContains(err, "field's model appears more than once, select which one you want to use with Appearance; model: models.ParentParent")
}

func (ts *GroupByIntTestSuite) TestGroupByWithConditionsBefore() {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)
	ts.createProduct("4", 0, 2.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
			conditions.Product.Float.Is().Eq(cql.Float64(1.0)),
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2}, {Int: 0, Aggregation1: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithSameCondition() {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 2.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)
	ts.createProduct("4", 0, 2.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Int.Aggregate().Sum().Eq(cql.Int(2)),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithDifferentCondition() {
	ts.internalTestGroupByHavingWithDifferentCondition(cql.Int(2))
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithDifferentConditionInt8() {
	ts.internalTestGroupByHavingWithDifferentCondition(cql.Int8(2))
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithDifferentConditionInt16() {
	ts.internalTestGroupByHavingWithDifferentCondition(cql.Int16(2))
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithDifferentConditionInt32() {
	ts.internalTestGroupByHavingWithDifferentCondition(cql.Int32(2))
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithDifferentConditionInt64() {
	ts.internalTestGroupByHavingWithDifferentCondition(cql.Int64(2))
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithDifferentConditionUInt() {
	ts.internalTestGroupByHavingWithDifferentCondition(cql.UInt(2))
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithDifferentConditionUInt8() {
	ts.internalTestGroupByHavingWithDifferentCondition(cql.UInt8(2))
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithDifferentConditionUInt16() {
	ts.internalTestGroupByHavingWithDifferentCondition(cql.UInt16(2))
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithDifferentConditionUInt32() {
	ts.internalTestGroupByHavingWithDifferentCondition(cql.UInt32(2))
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithDifferentConditionUInt64() {
	ts.internalTestGroupByHavingWithDifferentCondition(cql.UInt64(2))
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithDifferentConditionFloat32() {
	ts.internalTestGroupByHavingWithDifferentCondition(cql.Float32(2))
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithDifferentConditionFloat64() {
	ts.internalTestGroupByHavingWithDifferentCondition(cql.Float64(2))
}

func (ts *GroupByIntTestSuite) internalTestGroupByHavingWithDifferentCondition(value condition.ValueOfType[float64]) {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Int.Aggregate().Count().Eq(value),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithComparisonWithAggregationNumeric() {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)
	ts.createProduct("4", 0, 2.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Int.Aggregate().Count().Eq(conditions.Product.Int.Aggregate().Sum()),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithComparisonWithAggregationNumericOtherType() {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)
	ts.createProduct("4", 0, 2.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Int.Aggregate().Sum().Eq(conditions.Product.Float.Aggregate().Sum()),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithComparisonWithAggregationOfAnotherTable() {
	product0 := ts.createProduct("1", 2, 0, false, nil)
	product1 := ts.createProduct("1", 2, 1, false, nil)
	product2 := ts.createProduct("2", 3, 2, false, nil)

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
		).GroupBy(
			conditions.Sale.Code,
		).Having(
			conditions.Product.Float.Aggregate().Sum().Eq(conditions.Sale.ID.Aggregate().Count()),
		),
		cql.ValueInto(conditions.Sale.Code, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Float.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 1, Aggregation1: 3},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByMultipleHaving() {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)
	ts.createProduct("4", 0, 2.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Int.Aggregate().Count().Eq(cql.Int(2)),
			conditions.Product.Float.Aggregate().Sum().Eq(cql.Int(2)),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByMultipleHavingWithAndConnection() {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)
	ts.createProduct("4", 0, 2.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			cql.AndHaving(
				conditions.Product.Int.Aggregate().Count().Eq(cql.Int(2)),
				conditions.Product.Float.Aggregate().Sum().Eq(cql.Int(2)),
			),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByMultipleHavingWithOrConnection() {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)
	ts.createProduct("5", 2, 3.0, false, nil)
	ts.createProduct("6", 3, 2.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			cql.OrHaving(
				conditions.Product.Int.Aggregate().Count().Eq(cql.Int(2)),
				conditions.Product.Float.Aggregate().Sum().Eq(cql.Int(1)),
			),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 1, Aggregation1: 2},
		{Int: 0, Aggregation1: 0},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByWithNotConnection() {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)
	ts.createProduct("5", 2, 3.0, false, nil)
	ts.createProduct("6", 3, 2.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			cql.NotHaving(
				conditions.Product.Int.Aggregate().Count().Eq(cql.Int(2)),
			),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 0, Aggregation1: 0},
		{Int: 2, Aggregation1: 2},
		{Int: 3, Aggregation1: 3},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByWithNotConnectionMultiple() {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)
	ts.createProduct("5", 0, 3.0, false, nil)
	ts.createProduct("6", 3, 2.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			cql.NotHaving(
				conditions.Product.Int.Aggregate().Count().Eq(cql.Int(2)),
				conditions.Product.Float.Aggregate().Sum().Eq(cql.Int(2)),
			),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 0, Aggregation1: 0},
		{Int: 3, Aggregation1: 3},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingBoolean() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("3", 0, 1.0, true, nil)
	ts.createProduct("4", 0, 2.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Bool.Aggregate().All().Eq(cql.Bool(true)),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingBooleanCompareWithAnotherAggregation() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("3", 0, 1.0, true, nil)
	ts.createProduct("4", 0, 2.0, false, nil)
	ts.createProduct("4", 2, 2.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Bool.Aggregate().All().Eq(conditions.Product.Bool.Aggregate().Any()),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 1, Aggregation1: 2},
		{Int: 2, Aggregation1: 2},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingOtherType() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("3", 0, 1.0, true, nil)
	ts.createProduct("4", 0, 2.0, false, nil)
	ts.createProduct("4", 2, 2.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.String.Aggregate().Max().Eq(cql.String("4")),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 2, Aggregation1: 2}, {Int: 0, Aggregation1: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingOtherTypeCompareWithAnotherAggregation() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("3", 0, 1.0, true, nil)
	ts.createProduct("4", 0, 2.0, false, nil)
	ts.createProduct("4", 2, 2.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.String.Aggregate().Max().Eq(conditions.Product.String.Aggregate().Min()),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 2, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingNotEq() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("3", 0, 1.0, true, nil)
	ts.createProduct("4", 0, 2.0, false, nil)
	ts.createProduct("4", 2, 3.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Float.Aggregate().Max().NotEq(cql.Float64(2.0)),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2}, {Int: 2, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingLt() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("3", 0, 1.0, true, nil)
	ts.createProduct("4", 0, 2.0, false, nil)
	ts.createProduct("4", 2, 3.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Float.Aggregate().Max().Lt(cql.Float64(2.0)),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingLtOrEq() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("3", 0, 1.0, true, nil)
	ts.createProduct("4", 0, 2.0, false, nil)
	ts.createProduct("4", 2, 3.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Float.Aggregate().Max().LtOrEq(cql.Float64(2.0)),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2}, {Int: 0, Aggregation1: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingGt() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("3", 0, 1.0, true, nil)
	ts.createProduct("4", 0, 2.0, false, nil)
	ts.createProduct("4", 2, 3.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Float.Aggregate().Max().Gt(cql.Float64(2.0)),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 2, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingGtOrEq() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("3", 0, 1.0, true, nil)
	ts.createProduct("4", 0, 2.0, false, nil)
	ts.createProduct("4", 2, 3.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Float.Aggregate().Max().GtOrEq(cql.Float64(2.0)),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 2, Aggregation1: 2}, {Int: 0, Aggregation1: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingIn() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("3", 0, 1.0, true, nil)
	ts.createProduct("4", 0, 2.0, false, nil)
	ts.createProduct("4", 2, 3.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Float.Aggregate().Max().In([]float64{2.0, 3.0}),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 2, Aggregation1: 2}, {Int: 0, Aggregation1: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingNotIn() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("3", 0, 1.0, true, nil)
	ts.createProduct("4", 0, 2.0, false, nil)
	ts.createProduct("4", 2, 3.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Float.Aggregate().Max().NotIn([]float64{2.0, 3.0}),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingLike() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("03", 0, 1.0, true, nil)
	ts.createProduct("24", 0, 2.0, false, nil)
	ts.createProduct("14", 2, 3.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.String.Aggregate().Max().Like("_4"),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 2, Aggregation1: 2}, {Int: 0, Aggregation1: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithField() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("03", 0, 1.0, true, nil)
	ts.createProduct("24", 0, 2.0, false, nil)
	ts.createProduct("14", 2, 3.0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Float.Aggregate().Max().Eq(conditions.Product.Int),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{{Int: 1, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByMultipleHavingWithField() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("3", 1, 1.1, false, nil)
	ts.createProduct("4", 1, 1.1, false, nil)
	ts.createProduct("5", 0, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
			conditions.Product.Float,
		).Having(
			conditions.Product.Float.Aggregate().Max().Eq(conditions.Product.Int),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Float, func(value float64, result *Result) {
			result.Float = value
		}),
		cql.ValueInto(conditions.Product.Int.Aggregate().Sum(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 1, Float: 1, Aggregation1: 1},
		{Int: 0, Float: 0, Aggregation1: 0},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByAggregateAfterFunction() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 0, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Plus(cql.Int(123)).Aggregate().Max(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 1, Aggregation1: 124},
		{Int: 0, Aggregation1: 123},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByAggregateAfterFunctionDynamic() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 0, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Plus(conditions.Product.Int).Aggregate().Max(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 1, Aggregation1: 2},
		{Int: 0, Aggregation1: 0},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByAggregateHavingAfterFunction() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 0, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Int.Plus(cql.Int(12)).Aggregate().Max().Eq(cql.Int(13)),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Plus(cql.Int(123)).Aggregate().Max(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 1, Aggregation1: 124},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByAggregateHavingAfterFunctionDynamic() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 0, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			context.Background(),
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Having(
			conditions.Product.Int.Plus(conditions.Product.Float).Aggregate().Max().Eq(cql.Int(2)),
		),
		cql.ValueInto(conditions.Product.Int, func(value float64, result *Result) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Plus(conditions.Product.Int).Aggregate().Max(), func(value float64, result *Result) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result{
		{Int: 1, Aggregation1: 2},
		{Int: 0, Aggregation1: 0},
	}, results)
}
