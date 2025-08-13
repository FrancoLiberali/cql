package test

import (
	"gorm.io/gorm"

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
	Int          int
	Aggregation1 int
	Aggregation2 int
}

type ResultName struct {
	Name        string
	Aggregation int
}

type ResultCode struct {
	Code         int
	Aggregation1 int
	Aggregation2 int
}

type ResultIntPointer struct {
	Int         int
	Aggregation *int
}

type ResultFloat struct {
	Int         int
	Aggregation float64
}

type ResultString struct {
	Int         int
	Aggregation string
}

type ResultBool struct {
	Int         int
	Aggregation bool
}

type ResultIntAndFloat struct {
	Int          int
	Float        float64
	Aggregation1 int
	Aggregation2 float64
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

func (ts *GroupByIntTestSuite) TestGroupByFieldNotPresentReturnsError() {
	results := []ResultAlone{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(conditions.Sale.SellerID).Into(&results)

	ts.ErrorIs(err, cql.ErrFieldModelNotConcerned)
	ts.ErrorContains(err, "field's model is not concerned by the query (not joined); not concerned model: models.Sale")
}

func (ts *GroupByIntTestSuite) TestGroupBySelectIntoStructWithMoreFieldsWork() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(conditions.Product.Int).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1}, {Int: 0}}, results)
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

func (ts *GroupByIntTestSuite) TestGroupBySelectFieldNotPresentReturnsError() {
	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Sale.ID.Aggregate().Count(), "aggregation1",
	).Into(&results)

	ts.ErrorIs(err, cql.ErrFieldModelNotConcerned)
	ts.ErrorContains(err, "field's model is not concerned by the query (not joined); not concerned model: models.Sale")
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
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 2}, {Int: 0, Aggregation1: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectIntoCastIntToFloatWorks() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results := []ResultFloat{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.Int.Aggregate().Count(), "aggregation",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultFloat{{Int: 1, Aggregation: 2}, {Int: 0, Aggregation: 1}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectIntoCastIntToStringWorks() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results := []ResultString{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.Int.Aggregate().Count(), "aggregation",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultString{{Int: 1, Aggregation: "2"}, {Int: 0, Aggregation: "1"}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupBySelectIntoStructWithLessFieldsWork() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results := []ResultAlone{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultAlone{{Int: 1}, {Int: 0}}, results)
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
		conditions.Product.Int.Aggregate().Count(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 2}, {Int: 0, Aggregation1: 1}}, results)
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
		conditions.Product.IntPointer.Aggregate().Count(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 1}, {Int: 0, Aggregation1: 0}}, results)
}

// func (ts *GroupByIntTestSuite) TestGroupBySelectCountAll() {
// 	int1 := 1

// 	ts.createProduct("1", 1, 0, false, nil)
// 	ts.createProduct("2", 1, 0, false, &int1)
// 	ts.createProduct("3", 0, 0, false, nil)

// 	results := []ResultInt{}

// 	err := cql.Query[models.Product](
// 		ts.db,
// 	).GroupBy(
// 		conditions.Product.Int,
// 	).Select(
// 		cql.CountAll(), "aggregation1",
// 	).Into(&results)

// 	ts.Require().NoError(err)
// 	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 2}, {Int: 0, Aggregation1: 1}}, results)
// }

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

func (ts *GroupByIntTestSuite) TestGroupBySelectMinForNotNumericField() {
	ts.createProduct("1", 1, 0.25, false, nil)
	ts.createProduct("2", 1, 0.75, false, nil)
	ts.createProduct("3", 0, 1, false, nil)

	results := []ResultString{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.String.Aggregate().Min(), "aggregation",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultString{{Int: 1, Aggregation: "1"}, {Int: 0, Aggregation: "3"}}, results)
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

func (ts *GroupByIntTestSuite) TestGroupBySelectMaxForNotNumericField() {
	ts.createProduct("1", 1, 0.25, false, nil)
	ts.createProduct("2", 1, 0.75, false, nil)
	ts.createProduct("3", 0, 1, false, nil)

	results := []ResultString{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.String.Aggregate().Max(), "aggregation",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultString{{Int: 1, Aggregation: "2"}, {Int: 0, Aggregation: "3"}}, results)
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
			conditions.Product.IntPointer.Aggregate().And(), "aggregation1",
		).Into(&results)

		ts.Require().NoError(err)
		EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 1}, {Int: 2, Aggregation1: 3}, {Int: 0, Aggregation1: 0}, {Int: 3, Aggregation1: 1}}, results)
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
			conditions.Product.IntPointer.Aggregate().And(), "aggregation1",
		).Into(&results)

		ts.Require().NoError(err)
		EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 1}, {Int: 2, Aggregation1: 3}, {Int: 3, Aggregation1: 1}}, results)
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
			conditions.Product.IntPointer.Aggregate().Or(), "aggregation1",
		).Into(&results)

		ts.Require().NoError(err)
		EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 3}, {Int: 2, Aggregation1: 3}, {Int: 0, Aggregation1: 0}, {Int: 3, Aggregation1: 1}}, results)
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

func (ts *GroupByIntTestSuite) TestGroupBySelectMoreThanOne() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 0, false, nil)
	ts.createProduct("3", 0, 0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation2",
	).Select(
		conditions.Product.Int.Aggregate().Count(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 2, Aggregation2: 2}, {Int: 0, Aggregation1: 1, Aggregation2: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByMoreThanOne() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("3", 1, 1.1, false, nil)
	ts.createProduct("4", 1, 1.1, false, nil)
	ts.createProduct("5", 0, 0, false, nil)

	results := []ResultIntAndFloat{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
		conditions.Product.Float,
	).Select(
		conditions.Product.Int.Aggregate().Count(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultIntAndFloat{
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

	results := []ResultIntAndFloat{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int, conditions.Product.Float,
	).Select(
		conditions.Product.Int.Aggregate().Count(), "aggregation1",
	).Select(
		conditions.Product.Float.Aggregate().Sum(), "aggregation2",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultIntAndFloat{
		{Int: 1, Float: 0, Aggregation1: 1, Aggregation2: 0},
		{Int: 1, Float: 1, Aggregation1: 1, Aggregation2: 1},
		{Int: 1, Float: 1.1, Aggregation1: 2, Aggregation2: 2.2},
		{Int: 0, Float: 0, Aggregation1: 1, Aggregation2: 0},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByJoinedField() {
	product1 := ts.createProduct("1", 2, 0, false, nil)
	product2 := ts.createProduct("2", 3, 1, false, nil)

	ts.createSale(1, product1, nil)
	ts.createSale(1, product2, nil)
	ts.createSale(2, product1, nil)

	results := []ResultInt{}

	switch getDBDialector() {
	// TODO group by joined field doesn't work for Postgres by bug in gorm
	case sql.MySQL, sql.SQLServer, sql.SQLite:
		err := cql.Query[models.Sale](
			ts.db,
			conditions.Sale.Product(),
		).GroupBy(
			conditions.Product.Int,
		).Select(
			conditions.Sale.Code.Aggregate().Sum(), "aggregation1",
		).Into(&results)

		ts.Require().NoError(err)
		EqualList(&ts.Suite, []ResultInt{
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

	results := []ResultCode{}

	err := cql.Query[models.Sale](
		ts.db,
		conditions.Sale.Product(),
	).GroupBy(
		conditions.Sale.Code,
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultCode{
		{Code: 1, Aggregation1: 5},
		{Code: 2, Aggregation1: 2},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByJoinedFieldAndWithJoinedFieldInSelect() {
	product1 := ts.createProduct("1", 2, 1, false, nil)
	product2 := ts.createProduct("2", 3, 3, false, nil)

	ts.createSale(1, product1, nil)
	ts.createSale(1, product2, nil)
	ts.createSale(2, product1, nil)

	results := []ResultInt{}

	switch getDBDialector() {
	// TODO group by joined field doesn't work for Postgres by bug in gorm
	case sql.MySQL, sql.SQLServer, sql.SQLite:
		err := cql.Query[models.Sale](
			ts.db,
			conditions.Sale.Product(),
		).GroupBy(
			conditions.Product.Int,
		).Select(
			conditions.Product.Float.Aggregate().Sum(), "aggregation1",
		).Into(&results)

		ts.Require().NoError(err)
		EqualList(&ts.Suite, []ResultInt{
			{Int: 2, Aggregation1: 2},
			{Int: 3, Aggregation1: 3},
		}, results)
	}
}

func (ts *GroupByIntTestSuite) TestGroupByFieldPresentInMultipleTables() {
	company := ts.createCompany("name1")
	ts.createSeller("name1", company)
	ts.createSeller("name2", company)

	results := []ResultName{}

	err := cql.Query[models.Seller](
		ts.db,
		conditions.Seller.Company(),
	).GroupBy(
		conditions.Seller.Name,
	).Select(
		conditions.Company.Name.Aggregate().Count(), "aggregation",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultName{
		{Name: "name1", Aggregation: 1},
		{Name: "name2", Aggregation: 1},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByJoinedMultipleTimesFieldReturnsError() {
	results := []ResultInt{}

	err := cql.Query[models.Child](
		ts.db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
	).GroupBy(
		conditions.ParentParent.Name,
	).Into(&results)

	ts.ErrorIs(err, cql.ErrAppearanceMustBeSelected)
	ts.ErrorContains(err, "field's model appears more than once, select which one you want to use with Appearance; model: models.ParentParent")
}

func (ts *GroupByIntTestSuite) TestGroupByWithConditionsBefore() {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)
	ts.createProduct("4", 0, 2.0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Float.Is().Eq(1.0),
	).GroupBy(
		conditions.Product.Int,
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 2}, {Int: 0, Aggregation1: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithSameCondition() {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)
	ts.createProduct("4", 0, 2.0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Having(
		conditions.Product.Int.Aggregate().Sum().Eq(cql.Int(2)),
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 2}}, results)
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

func (ts *GroupByIntTestSuite) TestGroupByHavingWithDifferentConditionUIntPTR() {
	ts.internalTestGroupByHavingWithDifferentCondition(cql.UIntPTR(2))
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithDifferentConditionFloat32() {
	ts.internalTestGroupByHavingWithDifferentCondition(cql.Float32(2))
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithDifferentConditionFloat64() {
	ts.internalTestGroupByHavingWithDifferentCondition(cql.Float64(2))
}

func (ts *GroupByIntTestSuite) internalTestGroupByHavingWithDifferentCondition(value condition.NumericAggregationComparable) {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Having(
		conditions.Product.Int.Aggregate().Count().Eq(value),
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithComparisonWithAggregationNumeric() {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)
	ts.createProduct("4", 0, 2.0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Having(
		conditions.Product.Int.Aggregate().Count().Eq(conditions.Product.Int.Aggregate().Sum()),
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingWithComparisonWithAggregationNumericOtherType() {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)
	ts.createProduct("4", 0, 2.0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Having(
		conditions.Product.Int.Aggregate().Sum().Eq(conditions.Product.Float.Aggregate().Sum()),
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 2}}, results)
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

	results := []ResultCode{}

	err := cql.Query[models.Sale](
		ts.db,
		conditions.Sale.Product(),
	).GroupBy(
		conditions.Sale.Code,
	).Having(
		conditions.Product.Float.Aggregate().Sum().Eq(conditions.Sale.ID.Aggregate().Count()),
	).Select(
		conditions.Product.Float.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultCode{
		{Code: 1, Aggregation1: 3},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByMultipleHaving() {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)
	ts.createProduct("4", 0, 2.0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Having(
		conditions.Product.Int.Aggregate().Count().Eq(cql.Int(2)),
		conditions.Product.Float.Aggregate().Sum().Eq(cql.Int(2)),
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByMultipleHavingWithAndConnection() {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)
	ts.createProduct("4", 0, 2.0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Having(
		cql.AndHaving(
			conditions.Product.Int.Aggregate().Count().Eq(cql.Int(2)),
			conditions.Product.Float.Aggregate().Sum().Eq(cql.Int(2)),
		),
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByMultipleHavingWithOrConnection() {
	ts.createProduct("1", 1, 1.0, false, nil)
	ts.createProduct("2", 1, 1.0, false, nil)
	ts.createProduct("3", 0, 1.0, false, nil)
	ts.createProduct("5", 2, 3.0, false, nil)
	ts.createProduct("6", 3, 2.0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Having(
		cql.OrHaving(
			conditions.Product.Int.Aggregate().Count().Eq(cql.Int(2)),
			conditions.Product.Float.Aggregate().Sum().Eq(cql.Int(1)),
		),
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{
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

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Having(
		cql.NotHaving(
			conditions.Product.Int.Aggregate().Count().Eq(cql.Int(2)),
		),
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{
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

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Having(
		cql.NotHaving(
			conditions.Product.Int.Aggregate().Count().Eq(cql.Int(2)),
			conditions.Product.Float.Aggregate().Sum().Eq(cql.Int(2)),
		),
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{
		{Int: 0, Aggregation1: 0},
		{Int: 3, Aggregation1: 3},
	}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingBoolean() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("3", 0, 1.0, true, nil)
	ts.createProduct("4", 0, 2.0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Having(
		conditions.Product.Bool.Aggregate().All().Eq(cql.Bool(true)),
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingBooleanCompareWithAnotherAggregation() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("3", 0, 1.0, true, nil)
	ts.createProduct("4", 0, 2.0, false, nil)
	ts.createProduct("4", 2, 2.0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Having(
		conditions.Product.Bool.Aggregate().All().Eq(conditions.Product.Bool.Aggregate().Any()),
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 1, Aggregation1: 2}, {Int: 2, Aggregation1: 2}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingOtherType() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("3", 0, 1.0, true, nil)
	ts.createProduct("4", 0, 2.0, false, nil)
	ts.createProduct("4", 2, 2.0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Having(
		conditions.Product.String.Aggregate().Max().Eq(cql.String("4")),
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 2, Aggregation1: 2}, {Int: 0, Aggregation1: 0}}, results)
}

func (ts *GroupByIntTestSuite) TestGroupByHavingOtherTypeCompareWithAnotherAggregation() {
	ts.createProduct("1", 1, 1.0, true, nil)
	ts.createProduct("2", 1, 1.0, true, nil)
	ts.createProduct("3", 0, 1.0, true, nil)
	ts.createProduct("4", 0, 2.0, false, nil)
	ts.createProduct("4", 2, 2.0, false, nil)

	results := []ResultInt{}

	err := cql.Query[models.Product](
		ts.db,
	).GroupBy(
		conditions.Product.Int,
	).Having(
		conditions.Product.String.Aggregate().Max().Eq(conditions.Product.String.Aggregate().Min()),
	).Select(
		conditions.Product.Int.Aggregate().Sum(), "aggregation1",
	).Into(&results)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{{Int: 2, Aggregation1: 2}}, results)
}
