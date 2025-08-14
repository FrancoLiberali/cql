package testintegration

import (
	"gorm.io/gorm"
	"gotest.tools/assert"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/cql"
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

	updated, err := orm.Update[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(1),
	).Set(
		conditions.Product.Int.Set().Eq(0),
	)
	ts.Nil(err)
	ts.Equal(int64(0), updated)
}

func (ts *UpdateIntTestSuite) TestUpdateWhenAModelMatchConditions() {
	product := ts.createProduct("", 0, 0, false, nil)

	updated, err := orm.Update[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(0),
	).Set(
		conditions.Product.Int.Set().Eq(1),
	)
	ts.Nil(err)
	ts.Equal(int64(1), updated)

	productReturned, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(1),
	).FindOne()
	ts.Nil(err)

	ts.Equal(product.ID, productReturned.ID)
	ts.Equal(1, productReturned.Int)
	ts.NotEqual(product.UpdatedAt.UnixMicro(), productReturned.UpdatedAt.UnixMicro())
}

func (ts *UpdateIntTestSuite) TestUpdateWhenMultipleModelsMatchConditions() {
	product1 := ts.createProduct("1", 0, 0, false, nil)
	product2 := ts.createProduct("2", 0, 0, false, nil)

	updated, err := orm.Update[models.Product](
		ts.db,
		conditions.Product.Bool.Is().False(),
	).Set(
		conditions.Product.Int.Set().Eq(1),
	)
	ts.Nil(err)
	ts.Equal(int64(2), updated)

	products, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(1),
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

	updated, err := orm.Update[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(0),
	).Set(
		conditions.Product.Int.Set().Eq(1),
		conditions.Product.Bool.Set().Eq(true),
	)
	ts.Nil(err)
	ts.Equal(int64(1), updated)

	productReturned, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(1),
		conditions.Product.Bool.Is().True(),
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

	updated, err := orm.Update[models.Phone](
		ts.db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq("google"),
		),
	).Set(
		conditions.Phone.Name.Set().Eq("pixel 7"),
	)
	ts.Nil(err)
	ts.Equal(int64(1), updated)

	pixel7, err := orm.Query[models.Phone](
		ts.db,
		conditions.Phone.Name.Is().Eq("pixel 7"),
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

	updated, err := orm.Update[models.Sale](
		ts.db,
		conditions.Sale.Product(
			conditions.Product.Int.Is().Eq(1),
		),
		conditions.Sale.Seller(
			conditions.Seller.Name.Is().Eq("franco"),
		),
	).Set(
		conditions.Sale.Code.Set().Eq(1),
	)
	ts.Nil(err)
	ts.Equal(int64(1), updated)

	sale, err := orm.Query[models.Sale](
		ts.db,
		conditions.Sale.Code.Is().Eq(1),
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

	updated, err := orm.Update[models.Sale](
		ts.db,
		conditions.Sale.Seller(
			conditions.Seller.Name.Is().Eq("franco"),
			conditions.Seller.Company(
				conditions.Company.Name.Is().Eq("ditrit"),
			),
		),
	).Set(
		conditions.Sale.Code.Set().Eq(1),
	)
	ts.Nil(err)
	ts.Equal(int64(1), updated)

	sale, err := orm.Query[models.Sale](
		ts.db,
		conditions.Sale.Code.Is().Eq(1),
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

	updated, err := orm.Update[models.Phone](
		ts.db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq("google"),
		),
	).Set(
		conditions.Phone.Name.Set().Dynamic(conditions.Brand.Name),
	)

	ts.Nil(err)
	ts.Equal(int64(1), updated)

	phoneReturned, err := orm.Query[models.Phone](
		ts.db,
		conditions.Phone.Name.Is().Eq("google"),
	).FindOne()
	ts.Nil(err)

	ts.Equal(pixel.ID, phoneReturned.ID)
	ts.Equal("google", phoneReturned.Name)
	ts.NotEqual(pixel.UpdatedAt.UnixMicro(), phoneReturned.UpdatedAt.UnixMicro())
}

func (ts *UpdateIntTestSuite) TestUpdateDynamicWithoutJoinNumberReturnsErrorIfJoinedMoreThanOnce() {
	_, err := orm.Update[models.Child](
		ts.db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
	).Set(
		conditions.Child.Name.Set().Dynamic(conditions.ParentParent.Name),
	)

	ts.ErrorIs(err, cql.ErrJoinMustBeSelected)
	ts.ErrorContains(err, "joined multiple times model: models.ParentParent; method: Set")
}

func (ts *UpdateIntTestSuite) TestUpdateDynamicWithJoinNumber() {
	parentParent := &models.ParentParent{Name: "franco"}
	parent1 := &models.Parent1{ParentParent: *parentParent}
	parent2 := &models.Parent2{ParentParent: *parentParent}
	child := &models.Child{Parent1: *parent1, Parent2: *parent2, Name: "not_franco"}
	err := ts.db.Create(child).Error
	ts.Nil(err)

	updated, err := orm.Update[models.Child](
		ts.db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
	).Set(
		conditions.Child.Name.Set().Dynamic(conditions.ParentParent.Name, 0),
	)
	ts.Nil(err)
	ts.Equal(int64(1), updated)

	childReturned, err := orm.Query[models.Child](
		ts.db,
		conditions.Child.Name.Is().Eq("franco"),
	).FindOne()
	ts.Nil(err)

	ts.Equal(child.ID, childReturned.ID)
	ts.Equal("franco", childReturned.Name)
	ts.NotEqual(child.UpdatedAt.UnixMicro(), childReturned.UpdatedAt.UnixMicro())
}

func (ts *UpdateIntTestSuite) TestUpdateUnsafe() {
	product := ts.createProduct("", 0, 0, false, nil)

	updated, err := orm.Update[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(0),
	).Set(
		conditions.Product.Int.Set().Unsafe("1"),
	)
	ts.Nil(err)
	ts.Equal(int64(1), updated)

	productReturned, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(1),
	).FindOne()
	ts.Nil(err)

	ts.Equal(product.ID, productReturned.ID)
	ts.Equal(1, productReturned.Int)
	ts.NotEqual(product.UpdatedAt.UnixMicro(), productReturned.UpdatedAt.UnixMicro())
}

func (ts *UpdateIntTestSuite) TestUpdateReturning() {
	switch getDBDialector() {
	// update returning only supported for postgres, sqlite, sqlserver
	case cql.MySQL:
		_, err := orm.Update[models.Phone](
			ts.db,
		).Returning(nil).Set()
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: Returning")
	case cql.Postgres, cql.SQLite, cql.SQLServer:
		product := ts.createProduct("", 0, 0, false, nil)

		productsReturned := []models.Product{}
		updated, err := orm.Update[models.Product](
			ts.db,
			conditions.Product.Int.Is().Eq(0),
		).Returning(&productsReturned).Set(
			conditions.Product.Int.Set().Eq(1),
		)
		ts.Nil(err)
		ts.Equal(int64(1), updated)

		ts.Len(productsReturned, 1)
		productReturned := productsReturned[0]
		ts.Equal(product.ID, productReturned.ID)
		ts.Equal(1, productReturned.Int)
		ts.NotEqual(product.UpdatedAt.UnixMicro(), productReturned.UpdatedAt.UnixMicro())
	}
}

func (ts *UpdateIntTestSuite) TestUpdateReturningWithPreload() {
	switch getDBDialector() {
	// update returning with preload only supported for postgres
	case cql.SQLite, cql.SQLServer:
		salesReturned := []models.Sale{}
		_, err := orm.Update[models.Sale](
			ts.db,
			conditions.Sale.Code.Is().Eq(0),
			conditions.Sale.PreloadProduct(),
		).Returning(&salesReturned).Set(
			conditions.Sale.Code.Set().Eq(2),
		)
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "preloads in returning are not allowed for database")
		ts.ErrorContains(err, "method: Returning")
	case cql.Postgres:
		product1 := ts.createProduct("a_string", 1, 0.0, false, nil)
		product2 := ts.createProduct("", 2, 0.0, false, nil)

		sale1 := ts.createSale(0, product1, nil)
		ts.createSale(1, product2, nil)

		salesReturned := []models.Sale{}
		updated, err := orm.Update[models.Sale](
			ts.db,
			conditions.Sale.Code.Is().Eq(0),
			conditions.Sale.PreloadProduct(),
		).Returning(&salesReturned).Set(
			conditions.Sale.Code.Set().Eq(2),
		)
		ts.Nil(err)
		ts.Equal(int64(1), updated)

		ts.Len(salesReturned, 1)
		saleReturned := salesReturned[0]
		ts.Equal(sale1.ID, saleReturned.ID)
		ts.Equal(2, saleReturned.Code)
		ts.NotEqual(sale1.UpdatedAt.UnixMicro(), saleReturned.UpdatedAt.UnixMicro())
		productPreloaded, err := saleReturned.GetProduct()
		ts.Nil(err)
		assert.DeepEqual(ts.T(), product1, productPreloaded)
	}
}

func (ts *UpdateIntTestSuite) TestUpdateReturningWithPreloadAtSecondLevel() {
	// update returning with preloads only supported for postgres
	if getDBDialector() != cql.Postgres {
		return
	}

	product1 := ts.createProduct("a_string", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	company := ts.createCompany("ditrit")

	withCompany := ts.createSeller("with", company)
	withoutCompany := ts.createSeller("without", nil)

	sale1 := ts.createSale(0, product1, withCompany)
	ts.createSale(1, product2, withoutCompany)

	salesReturned := []models.Sale{}
	updated, err := orm.Update[models.Sale](
		ts.db,
		conditions.Sale.Code.Is().Eq(0),
		conditions.Sale.Seller(
			conditions.Seller.PreloadCompany(),
		),
	).Returning(&salesReturned).Set(
		conditions.Sale.Code.Set().Eq(2),
	)
	ts.Nil(err)
	ts.Equal(int64(1), updated)

	ts.Len(salesReturned, 1)
	saleReturned := salesReturned[0]
	ts.Equal(sale1.ID, saleReturned.ID)
	ts.Equal(2, saleReturned.Code)
	ts.NotEqual(sale1.UpdatedAt.UnixMicro(), saleReturned.UpdatedAt.UnixMicro())
	sellerPreloaded, err := saleReturned.GetSeller()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), withCompany, sellerPreloaded)
	companyPreloaded, err := sellerPreloaded.GetCompany()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), company, companyPreloaded)
}

func (ts *UpdateIntTestSuite) TestUpdateReturningWithPreloadCollection() {
	switch getDBDialector() {
	// update returning only supported for postgres, sqlite, sqlserver
	case cql.Postgres, cql.SQLite, cql.SQLServer:
		company := ts.createCompany("ditrit")
		seller1 := ts.createSeller("1", company)
		seller2 := ts.createSeller("2", company)

		companiesReturned := []models.Company{}
		updated, err := orm.Update[models.Company](
			ts.db,
			conditions.Company.Name.Is().Eq("ditrit"),
			conditions.Company.PreloadSellers(),
		).Returning(&companiesReturned).Set(
			conditions.Company.Name.Set().Eq("orness"),
		)
		ts.Nil(err)
		ts.Equal(int64(1), updated)

		ts.Len(companiesReturned, 1)
		companyReturned := companiesReturned[0]
		ts.Equal(company.ID, companyReturned.ID)
		ts.Equal("orness", companyReturned.Name)
		ts.NotEqual(company.UpdatedAt.UnixMicro(), companyReturned.UpdatedAt.UnixMicro())
		sellersPreloaded, err := companyReturned.GetSellers()
		ts.Nil(err)
		EqualList(&ts.Suite, []models.Seller{*seller1, *seller2}, sellersPreloaded)
	}
}

func (ts *UpdateIntTestSuite) TestUpdateMultipleTables() {
	// update join only supported for mysql
	if getDBDialector() != cql.MySQL {
		_, err := orm.Update[models.Phone](
			ts.db,
		).SetMultiple()
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: SetMultiple")
	} else {
		brand1 := ts.createBrand("google")
		brand2 := ts.createBrand("apple")

		pixel := ts.createPhone("pixel", *brand1)
		ts.createPhone("iphone", *brand2)

		updated, err := orm.Update[models.Phone](
			ts.db,
			conditions.Phone.Brand(
				conditions.Brand.Name.Is().Eq("google"),
			),
		).SetMultiple(
			conditions.Phone.Name.Set().Eq("7"),
			conditions.Brand.Name.Set().Eq("google pixel"),
		)
		ts.Nil(err)
		ts.Equal(int64(2), updated)

		pixel7, err := orm.Query[models.Phone](
			ts.db,
			conditions.Phone.Name.Is().Eq("7"),
		).FindOne()
		ts.Nil(err)

		ts.Equal(pixel.ID, pixel7.ID)
		ts.Equal("7", pixel7.Name)
		ts.NotEqual(pixel.UpdatedAt.UnixMicro(), pixel7.UpdatedAt.UnixMicro())

		googlePixel, err := orm.Query[models.Brand](
			ts.db,
			conditions.Brand.Name.Is().Eq("google pixel"),
		).FindOne()
		ts.Nil(err)

		ts.Equal(brand1.ID, googlePixel.ID)
		ts.Equal("google pixel", googlePixel.Name)
		ts.NotEqual(brand1.UpdatedAt.UnixMicro(), googlePixel.UpdatedAt.UnixMicro())
	}
}

func (ts *UpdateIntTestSuite) TestUpdateMultipleTablesReturnsErrorIfTableNotJoined() {
	// update join only supported for mysql
	if getDBDialector() != cql.MySQL {
		return
	}

	_, err := orm.Update[models.Phone](
		ts.db,
	).SetMultiple(
		conditions.Phone.Name.Set().Eq("7"),
		conditions.Brand.Name.Set().Eq("google pixel"),
	)
	ts.ErrorIs(err, cql.ErrFieldModelNotConcerned)
	ts.ErrorContains(err, "not concerned model: models.Brand; method: Set")
}

func (ts *UpdateIntTestSuite) TestUpdateOrderByLimit() {
	// update order by limit only supported for mysql
	if getDBDialector() != cql.MySQL {
		_, err := orm.Update[models.Product](
			ts.db,
			conditions.Product.Bool.Is().False(),
		).Ascending(
			conditions.Product.String,
		).Limit(1).Set(
			conditions.Product.Int.Set().Eq(1),
		)
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: Ascending")
	} else {
		product1 := ts.createProduct("1", 0, 0, false, nil)
		ts.createProduct("2", 0, 0, false, nil)

		updated, err := orm.Update[models.Product](
			ts.db,
			conditions.Product.Bool.Is().False(),
		).Ascending(
			conditions.Product.String,
		).Limit(1).Set(
			conditions.Product.Int.Set().Eq(1),
		)
		ts.Nil(err)
		ts.Equal(int64(1), updated)

		productReturned, err := orm.Query[models.Product](
			ts.db,
			conditions.Product.Int.Is().Eq(1),
		).FindOne()
		ts.Nil(err)

		ts.Equal(product1.ID, productReturned.ID)
		ts.Equal(1, productReturned.Int)
		ts.NotEqual(product1.UpdatedAt.UnixMicro(), productReturned.UpdatedAt.UnixMicro())
	}
}

func (ts *UpdateIntTestSuite) TestUpdateLimitWithoutOrderByReturnsError() {
	// update order by limit only supported for mysql
	if getDBDialector() != cql.MySQL {
		_, err := orm.Update[models.Product](
			ts.db,
			conditions.Product.Bool.Is().False(),
		).Limit(1).Set(
			conditions.Product.Int.Set().Eq(1),
		)
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: Limit")
	} else {
		_, err := orm.Update[models.Product](
			ts.db,
			conditions.Product.Bool.Is().False(),
		).Limit(1).Set(
			conditions.Product.Int.Set().Eq(1),
		)
		ts.ErrorIs(err, cql.ErrOrderByMustBeCalled)
		ts.ErrorContains(err, "method: Limit")
	}
}
