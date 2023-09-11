package testintegration

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type DeleteIntTestSuite struct {
	ORMIntTestSuite
}

func NewDeleteIntTestSuite(
	db *gorm.DB,
) *DeleteIntTestSuite {
	return &DeleteIntTestSuite{
		ORMIntTestSuite: ORMIntTestSuite{
			db: db,
		},
	}
}

func (ts *DeleteIntTestSuite) SetupTest() {
	CleanDB(ts.db)
}

func (ts *DeleteIntTestSuite) TearDownSuite() {
	CleanDB(ts.db)
}

func (ts *DeleteIntTestSuite) TestDeleteWhenNothingMatchConditions() {
	ts.createProduct("", 0, 0, false, nil)

	deleted, err := orm.Delete[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(1),
	).Exec()
	ts.Nil(err)
	ts.Equal(int64(0), deleted)
}

func (ts *DeleteIntTestSuite) TestDeleteWhenAModelMatchConditions() {
	ts.createProduct("", 0, 0, false, nil)

	deleted, err := orm.Delete[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(0),
	).Exec()
	ts.Nil(err)
	ts.Equal(int64(1), deleted)

	productReturned, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(1),
	).Find()
	ts.Nil(err)
	ts.Len(productReturned, 0)
}

func (ts *DeleteIntTestSuite) TestDeleteWhenMultipleModelsMatchConditions() {
	ts.createProduct("1", 0, 0, false, nil)
	ts.createProduct("2", 0, 0, false, nil)

	deleted, err := orm.Delete[models.Product](
		ts.db,
		conditions.Product.Bool.Is().False(),
	).Exec()
	ts.Nil(err)
	ts.Equal(int64(2), deleted)

	productReturned, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.Bool.Is().False(),
	).Find()
	ts.Nil(err)
	ts.Len(productReturned, 0)
}
