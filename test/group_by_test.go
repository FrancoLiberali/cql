package test

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
	"gorm.io/gorm"
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

type Result1 struct {
	Int int
}

func (ts *GroupByIntTestSuite) TestGroupByNoSelect() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results := []Result1{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(conditions.Product.Int).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result1{{Int: 1}, {Int: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByWithConditionsNoSelect() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results := []Result1{}

	err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(1),
	).GroupBy(conditions.Product.Int).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result1{{Int: 1}}, results)
}

type Result2 struct {
	Int int
	Sum int
}

func (ts *GroupByIntTestSuite) TestGroupBySelectAggregation() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results := []Result2{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "sum",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []Result2{{Int: 1, Sum: 2}, {Int: 0, Sum: 0}}, results)
}

// TODO multiple selects, joins (ver cuando tienen atributo que se llama igual deberia romperse lo que esta ahora), errores
