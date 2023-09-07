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
	ts.NotEqual(product.UpdatedAt.UnixMicro(), productReturned.UpdatedAt.UnixMicro())
}

func (ts *UpdateIntTestSuite) TestUpdateWhenMultipleModelsMatchConditions() {
	product1 := ts.createProduct("1", 0, 0, false, nil)
	product2 := ts.createProduct("2", 0, 0, false, nil)

	updated, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.BoolIs().Eq(false),
	).Update(
		conditions.Product.IntSet().Eq(1),
	)
	ts.Nil(err)
	ts.Equal(int64(2), updated)

	products, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(1),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{product1, product2}, products)
	ts.Equal(1, products[0].Int)
	ts.Equal(1, products[1].Int)
	ts.NotEqual(product1.UpdatedAt.UnixMicro(), products[0].UpdatedAt.UnixMicro())
	ts.NotEqual(product2.UpdatedAt.UnixMicro(), products[1].UpdatedAt.UnixMicro())
}

func (ts *UpdateIntTestSuite) TestUpdateMultipleFieldsAtTheSameTime() {
	product := ts.createProduct("", 0, 0, false, nil)

	updated, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(0),
	).Update(
		conditions.Product.IntSet().Eq(1),
		conditions.Product.BoolSet().Eq(true),
	)
	ts.Nil(err)
	ts.Equal(int64(1), updated)

	productReturned, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(1),
		conditions.Product.BoolIs().Eq(true),
	).FindOne()
	ts.Nil(err)

	ts.Equal(product.ID, productReturned.ID)
	ts.Equal(1, productReturned.Int)
	ts.True(productReturned.Bool)
	ts.NotEqual(product.UpdatedAt.UnixMicro(), productReturned.UpdatedAt.UnixMicro())
}

func (ts *UpdateIntTestSuite) TestUpdateWithJoinInConditions() {
	brand1 := ts.createBrand("google")
	brand2 := ts.createBrand("apple")

	pixel := ts.createPhone("pixel", *brand1)
	ts.createPhone("iphone", *brand2)

	updated, err := orm.NewQuery[models.Phone](
		ts.db,
		conditions.Phone.Brand(
			conditions.Brand.NameIs().Eq("google"),
		),
	).Update(
		conditions.Phone.NameSet().Eq("pixel 7"),
	)
	ts.Nil(err)
	ts.Equal(int64(1), updated)

	pixel7, err := orm.NewQuery[models.Phone](
		ts.db,
		conditions.Phone.NameIs().Eq("pixel 7"),
	).FindOne()
	ts.Nil(err)

	ts.Equal(pixel.ID, pixel7.ID)
	ts.Equal("pixel 7", pixel7.Name)
	ts.NotEqual(pixel.UpdatedAt.UnixMicro(), pixel7.UpdatedAt.UnixMicro())
}

func (ts *UpdateIntTestSuite) TestUpdateWithJoinDifferentEntitiesInConditions() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)
	ts.createSale(0, product1, seller2)
	ts.createSale(0, product2, seller1)

	updated, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.Product(
			conditions.Product.IntIs().Eq(1),
		),
		conditions.Sale.Seller(
			conditions.Seller.NameIs().Eq("franco"),
		),
	).Update(
		conditions.Sale.CodeSet().Eq(1),
	)
	ts.Nil(err)
	ts.Equal(int64(1), updated)

	sale, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.CodeIs().Eq(1),
	).FindOne()
	ts.Nil(err)

	ts.Equal(match.ID, sale.ID)
	ts.Equal(1, sale.Code)
	ts.NotEqual(match.UpdatedAt.UnixMicro(), sale.UpdatedAt.UnixMicro())
}

func (ts *UpdateIntTestSuite) TestUpdateWithMultilevelJoinInConditions() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("franco", company1)
	seller2 := ts.createSeller("agustin", company2)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	updated, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.Seller(
			conditions.Seller.NameIs().Eq("franco"),
			conditions.Seller.Company(
				conditions.Company.NameIs().Eq("ditrit"),
			),
		),
	).Update(
		conditions.Sale.CodeSet().Eq(1),
	)
	ts.Nil(err)
	ts.Equal(int64(1), updated)

	sale, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.CodeIs().Eq(1),
	).FindOne()
	ts.Nil(err)

	ts.Equal(match.ID, sale.ID)
	ts.Equal(1, sale.Code)
	ts.NotEqual(match.UpdatedAt.UnixMicro(), sale.UpdatedAt.UnixMicro())
}

func (ts *UpdateIntTestSuite) TestUpdateDynamic() {
	google := ts.createBrand("google")
	apple := ts.createBrand("apple")

	pixel := ts.createPhone("pixel", *google)
	ts.createPhone("iphone", *apple)

	updated, err := orm.NewQuery[models.Phone](
		ts.db,
		conditions.Phone.Brand(
			conditions.Brand.NameIs().Eq("google"),
		),
	).Update(
		conditions.Phone.NameSet().Dynamic(conditions.Brand.Name),
	)

	ts.Nil(err)
	ts.Equal(int64(1), updated)

	phoneReturned, err := orm.NewQuery[models.Phone](
		ts.db,
		conditions.Phone.NameIs().Eq("google"),
	).FindOne()
	ts.Nil(err)

	ts.Equal(pixel.ID, phoneReturned.ID)
	ts.Equal("google", phoneReturned.Name)
	ts.NotEqual(pixel.UpdatedAt.UnixMicro(), phoneReturned.UpdatedAt.UnixMicro())
}

func (ts *UpdateIntTestSuite) TestUpdateUnsafe() {
	product := ts.createProduct("", 0, 0, false, nil)

	updated, err := orm.NewQuery[models.Product](
		ts.db,
		conditions.Product.IntIs().Eq(0),
	).Update(
		conditions.Product.IntSet().Unsafe("1"),
		// TODO se pueden repetir? mirar si da error en la base o que hace
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
	ts.NotEqual(product.UpdatedAt.UnixMicro(), productReturned.UpdatedAt.UnixMicro())
}
