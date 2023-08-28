package testintegration

import (
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
