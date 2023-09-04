package testintegration

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type UpdateIntTestSuite struct {
	ORMIntTestSuite
}

func NewUpdateIntTestSuite(
	db *gorm.DB,
) *UpdateIntTestSuite {
	return &UpdateIntTestSuite{
		ORMIntTestSuite: ORMIntTestSuite{
			db: db,
		},
	}
}

func (ts *UpdateIntTestSuite) SetupTest() {
	CleanDB(ts.db)
}

func (ts *UpdateIntTestSuite) TearDownSuite() {
	CleanDB(ts.db)
}

func (ts *UpdateIntTestSuite) TestUpdateWhenNothingMatchConditions() {
	ts.createProduct("", 0, 0, false, nil)

	updated, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(1),
	).Update(
		conditions.Product.IntSet().Eq(0),
	)
	ts.Nil(err)
	ts.Equal(int64(0), updated)
}

func (ts *UpdateIntTestSuite) TestUpdateWhenAModelMatchConditions() {
	product := ts.createProduct("", 0, 0, false, nil)

	updated, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(0),
	).Update(
		conditions.Product.IntSet().Eq(1),
		// conditions.Product.IntSet().Dynamic(conditions.Sale.Code),
		// conditions.Product.IntSet().Unsafe("1"),
		// se pueden repetir? mirar si da error en la base o que hace
	)
	ts.Nil(err)
	ts.Equal(int64(1), updated)

	productReturned, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(1),
	).FindOne()
	ts.Nil(err)

	ts.Equal(product.ID, productReturned.ID)
	ts.Equal(1, productReturned.Int)
}
