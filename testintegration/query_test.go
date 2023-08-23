package testintegration

import (
	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
	"gotest.tools/assert"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/errors"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type QueryIntTestSuite struct {
	ORMIntTestSuite
}

func NewQueryIntTestSuite(
	db *gorm.DB,
) *QueryIntTestSuite {
	return &QueryIntTestSuite{
		ORMIntTestSuite: ORMIntTestSuite{
			db: db,
		},
	}
}

func (ts *QueryIntTestSuite) SetupTest() {
	CleanDB(ts.db)
}

func (ts *QueryIntTestSuite) TearDownSuite() {
	CleanDB(ts.db)
}

// ------------------------- FindOne --------------------------------

func (ts *QueryIntTestSuite) TestFindOneReturnsErrorIfConditionsDontMatch() {
	ts.createProduct("", 0, 0, false, nil)
	_, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(1),
	).FindOne()
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *QueryIntTestSuite) TestFindOneReturnsEntityIfConditionsMatch() {
	product := ts.createProduct("", 1, 0, false, nil)
	productReturned, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(1),
	).FindOne()
	ts.Nil(err)

	assert.DeepEqual(ts.T(), product, productReturned)
}

func (ts *QueryIntTestSuite) TestFindOneReturnsErrorIfMoreThanOneMatchConditions() {
	ts.createProduct("", 0, 0, false, nil)
	ts.createProduct("", 0, 0, false, nil)
	_, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(0),
	).FindOne()
	ts.Error(err, errors.ErrMoreThanOneObjectFound)
}

// ------------------------- First --------------------------------

func (ts *QueryIntTestSuite) TestFirstReturnsErrorIfConditionsDontMatch() {
	ts.createProduct("", 0, 0, false, nil)
	_, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(1),
	).First()
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *QueryIntTestSuite) TestFirstReturnsFirstEntityIfConditionsMatch() {
	brand1 := ts.createBrand("a")
	ts.createBrand("a")

	brandReturned, err := orm.NewQuery[models.Brand](
		ts.db,
		conditions.Brand.NameIs().Eq("a"),
	).First()
	ts.Nil(err)

	assert.DeepEqual(ts.T(), brand1, brandReturned)
}

// ------------------------- Last --------------------------------

func (ts *QueryIntTestSuite) TestLastReturnsErrorIfConditionsDontMatch() {
	ts.createProduct("", 0, 0, false, nil)
	_, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(1),
	).Last()
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *QueryIntTestSuite) TestLastReturnsLastEntityIfConditionsMatch() {
	ts.createBrand("a")
	brand2 := ts.createBrand("a")

	brandReturned, err := orm.NewQuery[models.Brand](
		ts.db,
		conditions.Brand.NameIs().Eq("a"),
	).Last()
	ts.Nil(err)

	assert.DeepEqual(ts.T(), brand2, brandReturned)
}

// ------------------------- Take --------------------------------

func (ts *QueryIntTestSuite) TestTakeReturnsErrorIfConditionsDontMatch() {
	ts.createProduct("", 0, 0, false, nil)
	_, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(1),
	).Take()
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *QueryIntTestSuite) TestTakeReturnsFirstCreatedEntityIfConditionsMatch() {
	product1 := ts.createProduct("", 1, 0, false, nil)
	product2 := ts.createProduct("", 1, 0, false, nil)
	productReturned, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(1),
	).Take()
	ts.Nil(err)

	ts.True(cmp.Equal(productReturned, product1) || cmp.Equal(productReturned, product2))
}
