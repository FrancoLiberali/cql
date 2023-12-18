package test

import (
	"gorm.io/gorm"
	"gotest.tools/assert"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

type DeleteIntTestSuite struct {
	TestSuite
}

func NewDeleteIntTestSuite(
	db *gorm.DB,
) *DeleteIntTestSuite {
	return &DeleteIntTestSuite{
		TestSuite: TestSuite{
			db: db,
		},
	}
}

func (ts *DeleteIntTestSuite) TestDeleteWhenNothingMatchConditions() {
	ts.createProduct("", 0, 0, false, nil)

	deleted, err := cql.Delete[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(1),
	).Exec()
	ts.Nil(err)
	ts.Equal(int64(0), deleted)
}

func (ts *DeleteIntTestSuite) TestDeleteWhenAModelMatchConditions() {
	ts.createProduct("", 0, 0, false, nil)

	deleted, err := cql.Delete[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(0),
	).Exec()
	ts.Nil(err)
	ts.Equal(int64(1), deleted)

	productReturned, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(1),
	).Find()
	ts.Nil(err)
	ts.Len(productReturned, 0)
}

func (ts *DeleteIntTestSuite) TestDeleteWhenMultipleModelsMatchConditions() {
	ts.createProduct("1", 0, 0, false, nil)
	ts.createProduct("2", 0, 0, false, nil)

	deleted, err := cql.Delete[models.Product](
		ts.db,
		conditions.Product.Bool.Is().False(),
	).Exec()
	ts.Nil(err)
	ts.Equal(int64(2), deleted)

	productReturned, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Bool.Is().False(),
	).Find()
	ts.Nil(err)
	ts.Len(productReturned, 0)
}

func (ts *DeleteIntTestSuite) TestDeleteWithJoinInConditions() {
	brand1 := ts.createBrand("google")
	brand2 := ts.createBrand("apple")

	ts.createPhone("pixel", *brand1)
	ts.createPhone("iphone", *brand2)

	deleted, err := cql.Delete[models.Phone](
		ts.db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq("google"),
		),
	).Exec()
	ts.Nil(err)
	ts.Equal(int64(1), deleted)

	phones, err := cql.Query[models.Phone](
		ts.db,
		conditions.Phone.Name.Is().Eq("pixel"),
	).Find()
	ts.Nil(err)
	ts.Len(phones, 0)
}

func (ts *DeleteIntTestSuite) TestDeleteWithJoinDifferentEntitiesInConditions() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	ts.createSale(0, product1, seller1)
	ts.createSale(1, product2, seller2)
	ts.createSale(2, product1, seller2)
	ts.createSale(3, product2, seller1)

	deleted, err := cql.Delete[models.Sale](
		ts.db,
		conditions.Sale.Product(
			conditions.Product.Int.Is().Eq(1),
		),
		conditions.Sale.Seller(
			conditions.Seller.Name.Is().Eq("franco"),
		),
	).Exec()
	ts.Nil(err)
	ts.Equal(int64(1), deleted)

	sales, err := cql.Query[models.Sale](
		ts.db,
		conditions.Sale.Code.Is().Eq(0),
	).Find()
	ts.Nil(err)
	ts.Len(sales, 0)
}

func (ts *DeleteIntTestSuite) TestDeleteWithMultilevelJoinInConditions() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("franco", company1)
	seller2 := ts.createSeller("agustin", company2)

	ts.createSale(0, product1, seller1)
	ts.createSale(1, product2, seller2)

	deleted, err := cql.Delete[models.Sale](
		ts.db,
		conditions.Sale.Seller(
			conditions.Seller.Name.Is().Eq("franco"),
			conditions.Seller.Company(
				conditions.Company.Name.Is().Eq("ditrit"),
			),
		),
	).Exec()
	ts.Nil(err)
	ts.Equal(int64(1), deleted)

	sales, err := cql.Query[models.Sale](
		ts.db,
		conditions.Sale.Code.Is().Eq(0),
	).Find()
	ts.Nil(err)
	ts.Len(sales, 0)
}

func (ts *DeleteIntTestSuite) TestDeleteReturning() {
	switch getDBDialector() {
	// update returning only supported for postgres, sqlite, sqlserver
	case condition.MySQL:
		_, err := cql.Delete[models.Phone](
			ts.db,
		).Returning(nil).Exec()
		ts.ErrorIs(err, condition.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: Returning")
	case condition.Postgres, condition.SQLite, condition.SQLServer:
		product := ts.createProduct("", 0, 0, false, nil)

		productsReturned := []models.Product{}
		deleted, err := cql.Delete[models.Product](
			ts.db,
			conditions.Product.Int.Is().Eq(0),
		).Returning(&productsReturned).Exec()
		ts.Nil(err)
		ts.Equal(int64(1), deleted)

		ts.Len(productsReturned, 1)
		productReturned := productsReturned[0]
		ts.Equal(product.ID, productReturned.ID)

		products, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Int.Is().Eq(0),
		).Find()
		ts.Nil(err)
		ts.Len(products, 0)
	}
}

func (ts *DeleteIntTestSuite) TestDeleteReturningWithPreload() {
	switch getDBDialector() {
	// update returning with preload only supported for postgres
	case condition.SQLite, condition.SQLServer:
		salesReturned := []models.Sale{}
		_, err := cql.Delete[models.Sale](
			ts.db,
			conditions.Sale.Code.Is().Eq(0),
			conditions.Sale.PreloadProduct(),
		).Returning(&salesReturned).Exec()
		ts.ErrorIs(err, condition.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "preloads in returning are not allowed for database")
		ts.ErrorContains(err, "method: Returning")
	case condition.Postgres:
		product1 := ts.createProduct("a_string", 1, 0.0, false, nil)
		product2 := ts.createProduct("", 2, 0.0, false, nil)

		sale1 := ts.createSale(0, product1, nil)
		ts.createSale(1, product2, nil)

		salesReturned := []models.Sale{}
		deleted, err := cql.Delete[models.Sale](
			ts.db,
			conditions.Sale.Code.Is().Eq(0),
			conditions.Sale.PreloadProduct(),
		).Returning(&salesReturned).Exec()
		ts.Nil(err)
		ts.Equal(int64(1), deleted)

		ts.Len(salesReturned, 1)
		saleReturned := salesReturned[0]
		ts.Equal(sale1.ID, saleReturned.ID)
		productPreloaded, err := saleReturned.GetProduct()
		ts.Nil(err)
		assert.DeepEqual(ts.T(), product1, productPreloaded)
	}
}

func (ts *DeleteIntTestSuite) TestDeleteReturningWithPreloadAtSecondLevel() {
	// update returning with preloads only supported for postgres
	if getDBDialector() != condition.Postgres {
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
	deleted, err := cql.Delete[models.Sale](
		ts.db,
		conditions.Sale.Code.Is().Eq(0),
		conditions.Sale.Seller(
			conditions.Seller.PreloadCompany(),
		),
	).Returning(&salesReturned).Exec()
	ts.Nil(err)
	ts.Equal(int64(1), deleted)

	ts.Len(salesReturned, 1)
	saleReturned := salesReturned[0]
	ts.Equal(sale1.ID, saleReturned.ID)
	sellerPreloaded, err := saleReturned.GetSeller()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), withCompany, sellerPreloaded)
	companyPreloaded, err := sellerPreloaded.GetCompany()
	ts.Nil(err)
	assert.DeepEqual(ts.T(), company, companyPreloaded)
}

func (ts *DeleteIntTestSuite) TestDeleteReturningWithPreloadCollection() {
	switch getDBDialector() {
	// update returning only supported for postgres, sqlite, sqlserver
	case condition.Postgres, condition.SQLite, condition.SQLServer:
		company := ts.createCompany("ditrit")
		seller1 := ts.createSeller("1", company)
		seller2 := ts.createSeller("2", company)

		companiesReturned := []models.Company{}
		deleted, err := cql.Delete[models.Company](
			ts.db,
			conditions.Company.Name.Is().Eq("ditrit"),
			conditions.Company.PreloadSellers(),
		).Returning(&companiesReturned).Exec()
		ts.Nil(err)
		ts.Equal(int64(1), deleted)

		ts.Len(companiesReturned, 1)
		companyReturned := companiesReturned[0]
		ts.Equal(company.ID, companyReturned.ID)
		sellersPreloaded, err := companyReturned.GetSellers()
		ts.Nil(err)
		EqualList(&ts.Suite, []models.Seller{*seller1, *seller2}, sellersPreloaded)
	}
}

func (ts *DeleteIntTestSuite) TestDeleteOrderByLimit() {
	// update order by limit only supported for mysql
	if getDBDialector() != condition.MySQL {
		_, err := cql.Delete[models.Product](
			ts.db,
			conditions.Product.Bool.Is().False(),
		).Ascending(
			conditions.Product.String,
		).Limit(1).Exec()
		ts.ErrorIs(err, condition.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: Ascending")
	} else {
		product1 := ts.createProduct("1", 0, 0, false, nil)
		ts.createProduct("2", 0, 0, false, nil)

		deleted, err := cql.Delete[models.Product](
			ts.db,
			conditions.Product.Bool.Is().False(),
		).Descending(
			conditions.Product.String,
		).Limit(1).Exec()
		ts.Nil(err)
		ts.Equal(int64(1), deleted)

		productReturned, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Int.Is().Eq(0),
		).FindOne()
		ts.Nil(err)

		ts.Equal(product1.ID, productReturned.ID)
	}
}

func (ts *DeleteIntTestSuite) TestDeleteLimitWithoutOrderByReturnsError() {
	// update order by limit only supported for mysql
	if getDBDialector() != condition.MySQL {
		_, err := cql.Delete[models.Product](
			ts.db,
			conditions.Product.Bool.Is().False(),
		).Limit(1).Exec()
		ts.ErrorIs(err, condition.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: Limit")
	} else {
		_, err := cql.Delete[models.Product](
			ts.db,
			conditions.Product.Bool.Is().False(),
		).Limit(1).Exec()
		ts.ErrorIs(err, condition.ErrOrderByMustBeCalled)
		ts.ErrorContains(err, "method: Limit")
	}
}
