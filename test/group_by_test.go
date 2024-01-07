package test

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/sql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

type GroupByIntTestSuite struct {
	testSuite
}

func NewGroupByIntTestSuite(
	db *gorm.DB,
) *GroupByIntTestSuite {
	return &GroupByIntTestSuite{
		testSuite: testSuite{
			db: db,
		},
	}
}

type ResultAlone struct {
	Int int
}

type ResultInt struct {
	Int         int
	Aggregation int
}

type ResultIntPointer struct {
	Int         int
	Aggregation *int
}

type ResultFloat struct {
	Int         int
	Aggregation float64
}

type ResultBool struct {
	Int         int
	Aggregation bool
}

func (ts *GroupByIntTestSuite) TestGroupByNoSelect() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results := []ResultAlone{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(conditions.Product.Int).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultAlone{{Int: 1}, {Int: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByWithConditionsNoSelect() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results := []ResultAlone{}

	err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(1),
	).GroupBy(conditions.Product.Int).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultAlone{{Int: 1}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectSum() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation: 2}, {Int: 0, Aggregation: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectCount() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.Int.Aggregate().Count(), "aggregation",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation: 2}, {Int: 0, Aggregation: 1}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectCountWithNulls() {
	int1 := 1

	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, &int1)
	ts.createProduct("3", 0, 0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.IntPointer.Aggregate().Count(), "aggregation",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation: 1}, {Int: 0, Aggregation: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectCountAll() {
	int1 := 1

	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, &int1)
	ts.createProduct("3", 0, 0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		cql.CountAll(), "aggregation",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation: 2}, {Int: 0, Aggregation: 1}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectAverage() {
	ts.createProduct("1", 1, 0.25, false, nil)
	ts.createProduct("2", 1, 0.75, false, nil)
	ts.createProduct("3", 0, 1, false, nil)

	results := []ResultFloat{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.Float.Aggregate().Average(), "aggregation",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultFloat{{Int: 1, Aggregation: 0.5}, {Int: 0, Aggregation: 1}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectMin() {
	ts.createProduct("1", 1, 0.25, false, nil)
	ts.createProduct("2", 1, 0.75, false, nil)
	ts.createProduct("3", 0, 1, false, nil)

	results := []ResultFloat{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.Float.Aggregate().Min(), "aggregation",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultFloat{{Int: 1, Aggregation: 0.25}, {Int: 0, Aggregation: 1}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectMax() {
	ts.createProduct("1", 1, 0.25, false, nil)
	ts.createProduct("2", 1, 0.75, false, nil)
	ts.createProduct("3", 0, 1, false, nil)

	results := []ResultFloat{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.Float.Aggregate().Max(), "aggregation",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultFloat{{Int: 1, Aggregation: 0.75}, {Int: 0, Aggregation: 1}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectAll() {
	ts.createProduct("1", 1, 0.25, true, nil)
	ts.createProduct("2", 1, 0.75, false, nil)
	ts.createProduct("1", 2, 0.25, true, nil)
	ts.createProduct("2", 2, 0.75, true, nil)
	ts.createProduct("3", 0, 1, true, nil)

	results := []ResultBool{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.Bool.Aggregate().All(), "aggregation",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultBool{{Int: 1, Aggregation: false}, {Int: 2, Aggregation: true}, {Int: 0, Aggregation: true}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectAny() {
	ts.createProduct("1", 1, 0.25, true, nil)
	ts.createProduct("2", 1, 0.75, false, nil)
	ts.createProduct("1", 2, 0.25, false, nil)
	ts.createProduct("2", 2, 0.75, false, nil)
	ts.createProduct("3", 0, 1, true, nil)

	results := []ResultBool{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.Bool.Aggregate().Any(), "aggregation",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultBool{{Int: 1, Aggregation: true}, {Int: 2, Aggregation: false}, {Int: 0, Aggregation: true}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectNone() {
	ts.createProduct("1", 1, 0.25, true, nil)
	ts.createProduct("2", 1, 0.75, false, nil)
	ts.createProduct("1", 2, 0.25, false, nil)
	ts.createProduct("2", 2, 0.75, false, nil)
	ts.createProduct("3", 0, 1, false, nil)

	results := []ResultBool{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.Bool.Aggregate().None(), "aggregation",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultBool{{Int: 1, Aggregation: false}, {Int: 2, Aggregation: true}, {Int: 0, Aggregation: true}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectAnd() {
	results := []ResultInt{}

	switch getDBDialector() {
	case sql.Postgres:
		int1 := 1
		int2 := 3
		int3 := 3

		ts.createProduct("1", 1, 0.25, false, &int1)
		ts.createProduct("2", 1, 0.75, false, &int2)
		ts.createProduct("1", 2, 0.25, false, &int3)
		ts.createProduct("3", 0, 1, false, nil)
		ts.createProduct("3", 3, 1, false, &int1)
		ts.createProduct("3", 3, 1, false, nil)

		err := cql.Query[models.Product](
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Select(
			conditions.Product.IntPointer.Aggregate().And(), "aggregation",
		).Into(&results)

		ts.Require().NoError(err)
		EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation: 1}, {Int: 2, Aggregation: 3}, {Int: 0, Aggregation: 0}, {Int: 3, Aggregation: 1}}, results)
	case sql.MySQL:
		int1 := 1
		int2 := 3
		int3 := 3

		ts.createProduct("1", 1, 0.25, false, &int1)
		ts.createProduct("2", 1, 0.75, false, &int2)
		ts.createProduct("1", 2, 0.25, false, &int3)
		ts.createProduct("3", 3, 1, false, &int1)

		err := cql.Query[models.Product](
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Select(
			conditions.Product.IntPointer.Aggregate().And(), "aggregation",
		).Into(&results)

		ts.Require().NoError(err)
		EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation: 1}, {Int: 2, Aggregation: 3}, {Int: 3, Aggregation: 1}}, results)
	case sql.SQLite, sql.SQLServer:
		err := cql.Query[models.Product](
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Select(
			conditions.Product.IntPointer.Aggregate().And(), "aggregation",
		).Into(&results)

		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "function: And")
		ts.ErrorContains(err, "method: Select")
	}
}

func (ts *GroupByIntTestSuite) TestGroupBySelectOr() {
	results := []ResultInt{}

	switch getDBDialector() {
	case sql.Postgres, sql.MySQL:
		int1 := 1
		int2 := 2
		int3 := 3

		ts.createProduct("1", 1, 0.25, false, &int1)
		ts.createProduct("2", 1, 0.75, false, &int2)
		ts.createProduct("1", 2, 0.25, false, &int3)
		ts.createProduct("3", 0, 1, false, nil)
		ts.createProduct("3", 3, 1, false, &int1)
		ts.createProduct("3", 3, 1, false, nil)

		err := cql.Query[models.Product](
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Select(
			conditions.Product.IntPointer.Aggregate().Or(), "aggregation",
		).Into(&results)

		ts.Require().NoError(err)
		EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation: 3}, {Int: 2, Aggregation: 3}, {Int: 0, Aggregation: 0}, {Int: 3, Aggregation: 1}}, results)
	case sql.SQLite, sql.SQLServer:
		err := cql.Query[models.Product](
			ts.db,
		).GroupBy(
			conditions.Product.Int,
		).Select(
			conditions.Product.IntPointer.Aggregate().Or(), "aggregation",
		).Into(&results)

		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "function: Or")
		ts.ErrorContains(err, "method: Select")
	}
}

// TODO multiple selects, joins (ver cuando tienen atributo que se llama igual deberia romperse lo que esta ahora), errores
