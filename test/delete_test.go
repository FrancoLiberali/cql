package test

import (
	"context"

	"gotest.tools/assert"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/sql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

type DeleteIntTestSuite struct {
	testSuite
}

func NewDeleteIntTestSuite(
	db *cql.DB,
) *DeleteIntTestSuite {
	return &DeleteIntTestSuite{
		testSuite: testSuite{
			db: db,
		},
	}
}

func (ts *DeleteIntTestSuite) TestDeleteWithTrue() {
	ts.createProductNoTimestamps("", 0, 0, false, nil)
	ts.createProductNoTimestamps("", 1, 0, false, nil)

	deleted, err := cql.Delete[models.ProductNoTimestamps](
		context.Background(),
		ts.db,
		cql.True[models.ProductNoTimestamps](),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(2), deleted)

	productsReturned, err := cql.Query[models.ProductNoTimestamps](
		context.Background(),
		ts.db,
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 0)
}

func (ts *DeleteIntTestSuite) TestDeleteWhenNothingMatchConditions() {
	ts.createProductNoTimestamps("", 0, 0, false, nil)

	deleted, err := cql.Delete[models.ProductNoTimestamps](
		context.Background(),
		ts.db,
		conditions.ProductNoTimestamps.Int.Is().Eq(cql.Int(1)),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(0), deleted)
}

func (ts *DeleteIntTestSuite) TestDeleteWhenAModelMatchConditions() {
	ts.createProductNoTimestamps("", 0, 0, false, nil)

	deleted, err := cql.Delete[models.ProductNoTimestamps](
		context.Background(),
		ts.db,
		conditions.ProductNoTimestamps.Int.Is().Eq(cql.Int(0)),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(1), deleted)

	productReturned, err := cql.Query[models.ProductNoTimestamps](
		context.Background(),
		ts.db,
		conditions.ProductNoTimestamps.Int.Is().Eq(cql.Int(1)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productReturned, 0)
}

func (ts *DeleteIntTestSuite) TestDeleteWhenMultipleModelsMatchConditions() {
	ts.createProductNoTimestamps("1", 0, 0, false, nil)
	ts.createProductNoTimestamps("2", 0, 0, false, nil)

	deleted, err := cql.Delete[models.ProductNoTimestamps](
		context.Background(),
		ts.db,
		conditions.ProductNoTimestamps.Bool.Is().False(),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(2), deleted)

	productReturned, err := cql.Query[models.ProductNoTimestamps](
		context.Background(),
		ts.db,
		conditions.ProductNoTimestamps.Bool.Is().False(),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productReturned, 0)
}

func (ts *DeleteIntTestSuite) TestDeleteWithJoinInConditions() {
	brand1 := ts.createBrand("google")
	brand2 := ts.createBrand("apple")

	ts.createPhoneNoTimestamps("pixel", *brand1)
	ts.createPhoneNoTimestamps("iphone", *brand2)

	deleted, err := cql.Delete[models.PhoneNoTimestamps](
		context.Background(),
		ts.db,
		conditions.PhoneNoTimestamps.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("google")),
		),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(1), deleted)

	// joins delete no funciona, no veo nada en gorm de tests de que esto funcione
	// ts.db.GormDB.
	// 	Unscoped().
	// 	Joins("JOIN brands ON brands.id = phone_no_timestamps.brand_id"). //AND (Brand.name = ?) AND Brand.deleted_at IS NULL").
	// 	// Where("1=1").
	// 	Delete(&models.PhoneNoTimestamps{})

	// ts.db.GormDB.
	// 	Where(
	// 		"id IN (?)",
	// 		ts.db.GormDB.Model(&models.PhoneNoTimestamps{}).
	// 			Joins("JOIN brands ON brands.id = phone_no_timestamps.brand_id").
	// 			// Where("roles.name = ?", "admin").
	// 			Select("phone_no_timestamps.id"),
	// 	).
	// 	Delete(&models.PhoneNoTimestamps{})

	phones, err := cql.Query[models.PhoneNoTimestamps](
		context.Background(),
		ts.db,
		conditions.PhoneNoTimestamps.Name.Is().Eq(cql.String("pixel")),
	).Find()
	ts.Require().NoError(err)
	ts.Len(phones, 0)
}

func (ts *DeleteIntTestSuite) TestDeleteWithJoinDifferentEntitiesInConditions() {
	product1 := ts.createProductNoTimestamps("", 1, 0.0, false, nil)
	product2 := ts.createProductNoTimestamps("", 2, 0.0, false, nil)

	seller1 := ts.createSellerNoTimestamps("franco", nil)
	seller2 := ts.createSellerNoTimestamps("agustin", nil)

	ts.createSaleNoTimestamps(0, product1, seller1)
	ts.createSaleNoTimestamps(1, product2, seller2)
	ts.createSaleNoTimestamps(2, product1, seller2)
	ts.createSaleNoTimestamps(3, product2, seller1)

	deleted, err := cql.Delete[models.SaleNoTimestamps](
		context.Background(),
		ts.db,
		conditions.SaleNoTimestamps.Product(
			conditions.ProductNoTimestamps.Int.Is().Eq(cql.Int(1)),
		),
		conditions.SaleNoTimestamps.Seller(
			conditions.SellerNoTimestamps.Name.Is().Eq(cql.String("franco")),
		),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(1), deleted)

	sales, err := cql.Query[models.SaleNoTimestamps](
		context.Background(),
		ts.db,
		conditions.SaleNoTimestamps.Code.Is().Eq(cql.Int(0)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(sales, 0)
}

func (ts *DeleteIntTestSuite) TestDeleteWithMultilevelJoinInConditions() {
	product1 := ts.createProductNoTimestamps("", 0, 0.0, false, nil)
	product2 := ts.createProductNoTimestamps("", 0, 0.0, false, nil)

	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSellerNoTimestamps("franco", company1)
	seller2 := ts.createSellerNoTimestamps("agustin", company2)

	ts.createSaleNoTimestamps(0, product1, seller1)
	ts.createSaleNoTimestamps(1, product2, seller2)

	deleted, err := cql.Delete[models.SaleNoTimestamps](
		context.Background(),
		ts.db,
		conditions.SaleNoTimestamps.Seller(
			conditions.SellerNoTimestamps.Name.Is().Eq(cql.String("franco")),
			conditions.SellerNoTimestamps.CompanyNoTimestamps(
				conditions.CompanyNoTimestamps.Name.Is().Eq(cql.String("ditrit")),
			),
		),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(1), deleted)

	sales, err := cql.Query[models.SaleNoTimestamps](
		context.Background(),
		ts.db,
		conditions.SaleNoTimestamps.Code.Is().Eq(cql.Int(0)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(sales, 0)
}

func (ts *DeleteIntTestSuite) TestDeleteReturning() {
	switch getDBDialector() {
	// delete returning only supported for postgres, sqlite, sqlserver
	case sql.MySQL:
		_, err := cql.Delete[models.PhoneNoTimestamps](
			context.Background(),
			ts.db,
			conditions.PhoneNoTimestamps.Name.Is().Eq(cql.String("asd")),
		).Returning(nil).Exec()
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: Returning")
	case sql.Postgres, sql.SQLite, sql.SQLServer:
		product := ts.createProductNoTimestamps("", 0, 0, false, nil)

		productsReturned := []models.ProductNoTimestamps{}
		deleted, err := cql.Delete[models.ProductNoTimestamps](
			context.Background(),
			ts.db,
			conditions.ProductNoTimestamps.Int.Is().Eq(cql.Int(0)),
		).Returning(&productsReturned).Exec()
		ts.Require().NoError(err)
		ts.Equal(int64(1), deleted)

		ts.Len(productsReturned, 1)
		productReturned := productsReturned[0]
		ts.Equal(product.ID, productReturned.ID)

		products, err := cql.Query[models.ProductNoTimestamps](
			context.Background(),
			ts.db,
			conditions.ProductNoTimestamps.Int.Is().Eq(cql.Int(0)),
		).Find()
		ts.Require().NoError(err)
		ts.Len(products, 0)
	}
}

func (ts *DeleteIntTestSuite) TestDeleteReturningWithPreload() {
	switch getDBDialector() {
	// delete returning with preload only supported for postgres
	case sql.SQLite, sql.SQLServer:
		salesReturned := []models.SaleNoTimestamps{}
		_, err := cql.Delete[models.SaleNoTimestamps](
			context.Background(),
			ts.db,
			conditions.SaleNoTimestamps.Code.Is().Eq(cql.Int(0)),
			conditions.SaleNoTimestamps.Product().Preload(),
		).Returning(&salesReturned).Exec()
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "preloads in returning are not allowed for database")
		ts.ErrorContains(err, "method: Returning")
	case sql.Postgres:
		product1 := ts.createProductNoTimestamps("a_string", 1, 0.0, false, nil)
		product2 := ts.createProductNoTimestamps("", 2, 0.0, false, nil)

		sale1 := ts.createSaleNoTimestamps(0, product1, nil)
		ts.createSaleNoTimestamps(1, product2, nil)

		salesReturned := []models.SaleNoTimestamps{}
		deleted, err := cql.Delete[models.SaleNoTimestamps](
			context.Background(),
			ts.db,
			conditions.SaleNoTimestamps.Code.Is().Eq(cql.Int(0)),
			conditions.SaleNoTimestamps.Product().Preload(),
		).Returning(&salesReturned).Exec()
		ts.Require().NoError(err)
		ts.Equal(int64(1), deleted)

		ts.Len(salesReturned, 1)
		saleReturned := salesReturned[0]
		ts.Equal(sale1.ID, saleReturned.ID)
		productPreloaded, err := saleReturned.GetProduct()
		ts.Require().NoError(err)
		assert.DeepEqual(ts.T(), product1, productPreloaded)
	}
}

func (ts *DeleteIntTestSuite) TestDeleteReturningWithPreloadAtSecondLevel() {
	// delete returning with preloads only supported for postgres
	if getDBDialector() != sql.Postgres {
		return
	}

	product1 := ts.createProductNoTimestamps("a_string", 1, 0.0, false, nil)
	product2 := ts.createProductNoTimestamps("", 2, 0.0, false, nil)

	company := ts.createCompany("ditrit")

	withCompany := ts.createSellerNoTimestamps("with", company)
	withoutCompany := ts.createSellerNoTimestamps("without", nil)

	sale1 := ts.createSaleNoTimestamps(0, product1, withCompany)
	ts.createSaleNoTimestamps(1, product2, withoutCompany)

	salesReturned := []models.SaleNoTimestamps{}
	deleted, err := cql.Delete[models.SaleNoTimestamps](
		context.Background(),
		ts.db,
		conditions.SaleNoTimestamps.Code.Is().Eq(cql.Int(0)),
		conditions.SaleNoTimestamps.Seller(
			conditions.SellerNoTimestamps.CompanyNoTimestamps().Preload(),
		),
	).Returning(&salesReturned).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(1), deleted)

	ts.Len(salesReturned, 1)
	saleReturned := salesReturned[0]
	ts.Equal(sale1.ID, saleReturned.ID)
	sellerPreloaded, err := saleReturned.GetSeller()
	ts.Require().NoError(err)
	assert.DeepEqual(ts.T(), withCompany, sellerPreloaded)
	companyPreloaded, err := sellerPreloaded.GetCompanyNoTimestamps()
	ts.Require().NoError(err)
	assert.DeepEqual(ts.T(), company, companyPreloaded)
}

func (ts *DeleteIntTestSuite) TestDeleteReturningWithPreloadCollection() {
	switch getDBDialector() {
	// delete returning only supported for postgres, sqlite, sqlserver
	case sql.Postgres, sql.SQLite, sql.SQLServer:
		company := ts.createCompany("ditrit")
		seller1 := ts.createSellerNoTimestamps("1", company)
		seller2 := ts.createSellerNoTimestamps("2", company)

		companiesReturned := []models.CompanyNoTimestamps{}
		deleted, err := cql.Delete[models.CompanyNoTimestamps](
			context.Background(),
			ts.db,
			conditions.CompanyNoTimestamps.Name.Is().Eq(cql.String("ditrit")),
			conditions.CompanyNoTimestamps.Sellers.Preload(),
		).Returning(&companiesReturned).Exec()
		ts.Require().NoError(err)
		ts.Equal(int64(1), deleted)

		ts.Len(companiesReturned, 1)
		companyReturned := companiesReturned[0]
		ts.Equal(company.ID, companyReturned.ID)
		sellersPreloaded, err := companyReturned.GetSellers()
		ts.Require().NoError(err)
		EqualList(&ts.Suite, []models.SellerNoTimestamps{*seller1, *seller2}, sellersPreloaded)
	}
}

func (ts *DeleteIntTestSuite) TestDeleteOrderByLimit() {
	// delete order by limit only supported for mysql
	if getDBDialector() != sql.MySQL {
		_, err := cql.Delete[models.ProductNoTimestamps](
			context.Background(),
			ts.db,
			conditions.ProductNoTimestamps.Bool.Is().False(),
		).Ascending(
			conditions.ProductNoTimestamps.String,
		).Limit(1).Exec()
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: Ascending")
	} else {
		product1 := ts.createProductNoTimestamps("1", 0, 0, false, nil)
		ts.createProductNoTimestamps("2", 0, 0, false, nil)

		deleted, err := cql.Delete[models.ProductNoTimestamps](
			context.Background(),
			ts.db,
			conditions.ProductNoTimestamps.Bool.Is().False(),
		).Descending(
			conditions.ProductNoTimestamps.String,
		).Limit(1).Exec()
		ts.Require().NoError(err)
		ts.Equal(int64(1), deleted)

		productReturned, err := cql.Query[models.ProductNoTimestamps](
			context.Background(),
			ts.db,
			conditions.ProductNoTimestamps.Int.Is().Eq(cql.Int(0)),
		).FindOne()
		ts.Require().NoError(err)

		ts.Equal(product1.ID, productReturned.ID)
	}
}

func (ts *DeleteIntTestSuite) TestDeleteLimitWithoutOrderByReturnsError() {
	// delete order by limit only supported for mysql
	if getDBDialector() != sql.MySQL {
		_, err := cql.Delete[models.ProductNoTimestamps](
			context.Background(),
			ts.db,
			conditions.ProductNoTimestamps.Bool.Is().False(),
		).Limit(1).Exec()
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: Limit")
	} else {
		_, err := cql.Delete[models.ProductNoTimestamps](
			context.Background(),
			ts.db,
			conditions.ProductNoTimestamps.Bool.Is().False(),
		).Limit(1).Exec()
		ts.ErrorIs(err, cql.ErrOrderByMustBeCalled)
		ts.ErrorContains(err, "method: Limit")
	}
}
