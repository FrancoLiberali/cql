package testintegration

import (
	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
	"gotest.tools/assert"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/errors"
	"github.com/ditrit/badaas/orm/query"
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

// ------------------------- Order --------------------------------

func (ts *QueryIntTestSuite) TestAscendingReturnsResultsInAscendingOrder() {
	product1 := ts.createProduct("", 1, 1.0, false, nil)
	product2 := ts.createProduct("", 1, 2.0, false, nil)
	products, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(1),
	).Ascending(conditions.Product.Float).Find()
	ts.Nil(err)

	ts.Len(products, 2)
	assert.DeepEqual(ts.T(), product1, products[0])
	assert.DeepEqual(ts.T(), product2, products[1])
}

func (ts *QueryIntTestSuite) TestDescendingReturnsResultsInDescendingOrder() {
	product1 := ts.createProduct("", 1, 1.0, false, nil)
	product2 := ts.createProduct("", 1, 2.0, false, nil)
	products, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(1),
	).Descending(conditions.Product.Float).Find()
	ts.Nil(err)

	ts.Len(products, 2)
	assert.DeepEqual(ts.T(), product2, products[0])
	assert.DeepEqual(ts.T(), product1, products[1])
}

func (ts *QueryIntTestSuite) TestOrderByFieldThatIsJoined() {
	product1 := ts.createProduct("", 0, 1.0, false, nil)
	product2 := ts.createProduct("", 0, 2.0, false, nil)

	sale1 := ts.createSale(0, product1, nil)
	sale2 := ts.createSale(0, product2, nil)

	sales, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.Product(),
	).Descending(conditions.Product.Float).Find()
	ts.Nil(err)

	ts.Len(sales, 2)
	assert.DeepEqual(ts.T(), sale2, sales[0])
	assert.DeepEqual(ts.T(), sale1, sales[1])
}

func (ts *QueryIntTestSuite) TestOrderReturnsErrorIfFieldIsNotConcerned() {
	_, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(1),
	).Descending(conditions.Seller.ID).Find()
	ts.ErrorIs(err, errors.ErrFieldModelNotConcerned)
	ts.ErrorContains(err, "not concerned model: models.Seller; method: Descending")
}

func (ts *QueryIntTestSuite) TestOrderReturnsErrorIfFieldIsJoinedMoreThanOnceAndJoinIsNotSelected() {
	_, err := orm.NewQuery[models.Child](
		ts.db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
	).Descending(conditions.ParentParent.ID).Find()
	ts.ErrorIs(err, errors.ErrJoinMustBeSelected)
	ts.ErrorContains(err, "joined multiple times model: models.ParentParent; method: Descending")
}

func (ts *QueryIntTestSuite) TestOrderWorksIfFieldIsJoinedMoreThanOnceAndJoinIsSelected() {
	parentParent1 := &models.ParentParent{Name: "a"}
	parent11 := &models.Parent1{ParentParent: *parentParent1}
	parent12 := &models.Parent2{ParentParent: *parentParent1}
	child1 := &models.Child{Parent1: *parent11, Parent2: *parent12, Name: "franco"}
	err := ts.db.Create(child1).Error
	ts.Nil(err)

	parentParent2 := &models.ParentParent{Name: "b"}
	parent21 := &models.Parent1{ParentParent: *parentParent2}
	parent22 := &models.Parent2{ParentParent: *parentParent2}
	child2 := &models.Child{Parent1: *parent21, Parent2: *parent22, Name: "franco"}
	err = ts.db.Create(child2).Error
	ts.Nil(err)

	children, err := orm.NewQuery[models.Child](
		ts.db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
	).Ascending(conditions.ParentParent.Name, 0).Find()
	ts.Nil(err)

	ts.Len(children, 2)
	assert.DeepEqual(ts.T(), child1, children[0])
	assert.DeepEqual(ts.T(), child2, children[1])
}

// ------------------------- Limit --------------------------------

func (ts *QueryIntTestSuite) TestLimitLimitsTheAmountOfModelsReturned() {
	product1 := ts.createProduct("", 1, 0, false, nil)
	product2 := ts.createProduct("", 1, 0, false, nil)
	products, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(1),
	).Limit(1).Find()
	ts.Nil(err)

	ts.Len(products, 1)
	ts.True(cmp.Equal(products[0], product1) || cmp.Equal(products[0], product2))
}

func (ts *QueryIntTestSuite) TestLimitCanBeCanceled() {
	product1 := ts.createProduct("", 1, 0, false, nil)
	product2 := ts.createProduct("", 1, 0, false, nil)
	products, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(1),
	).Limit(1).Limit(-1).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{product1, product2}, products)
}

// ------------------------- Offset --------------------------------

func (ts *QueryIntTestSuite) TestOffsetSkipsTheModelsReturned() {
	ts.createProduct("", 1, 1, false, nil)
	product2 := ts.createProduct("", 1, 2, false, nil)

	switch getDBDialector() {
	case query.Postgres, query.SQLServer, query.SQLite:
		products, err := orm.NewQuery[models.Product](
			ts.db,
			conditions.Product.IntIs().Eq(1),
		).Ascending(conditions.Product.Float).Offset(1).Find()
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{product2}, products)
	case query.MySQL:
		products, err := orm.NewQuery[models.Product](
			ts.db,
			conditions.Product.IntIs().Eq(1),
		).Ascending(conditions.Product.Float).Offset(1).Limit(10).Find()
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{product2}, products)
	}
}

func (ts *QueryIntTestSuite) TestOffsetReturnsEmptyIfMoreOffsetThanResults() {
	ts.createProduct("", 1, 0, false, nil)
	ts.createProduct("", 1, 0, false, nil)

	switch getDBDialector() {
	case query.Postgres, query.SQLServer, query.SQLite:
		products, err := orm.NewQuery[models.Product](
			ts.db,
			conditions.Product.IntIs().Eq(1),
		).Offset(2).Find()
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{}, products)
	case query.MySQL:
		products, err := orm.NewQuery[models.Product](
			ts.db,
			conditions.Product.IntIs().Eq(1),
		).Offset(2).Limit(10).Find()
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{}, products)
	}
}

func (ts *QueryIntTestSuite) TestOffsetAndLimitWorkTogether() {
	ts.createProduct("", 1, 1, false, nil)
	product2 := ts.createProduct("", 1, 2, false, nil)
	ts.createProduct("", 1, 3, false, nil)
	products, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(1),
	).Ascending(conditions.Product.Float).Offset(1).Limit(1).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{product2}, products)
}
