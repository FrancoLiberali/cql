package test

import (
	"context"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/sql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

type SoftDeleteIntTestSuite struct {
	testSuite
}

func NewSoftDeleteIntTestSuite(
	db *cql.DB,
) *SoftDeleteIntTestSuite {
	return &SoftDeleteIntTestSuite{
		testSuite: testSuite{
			db: db,
		},
	}
}

func (ts *SoftDeleteIntTestSuite) TestSoftDeleteWithTrue() {
	ts.createProduct("", 0, 0, false, nil)
	ts.createProduct("", 1, 0, false, nil)

	deleted, err := cql.Delete[models.Product](
		context.Background(),
		ts.db,
		cql.True[models.Product](),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(2), deleted)

	productsReturned, err := cql.Query[models.Product](
		context.Background(),
		ts.db,
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 0)
}

func (ts *SoftDeleteIntTestSuite) TestSoftDeleteWhenNothingMatchConditions() {
	ts.createProduct("", 0, 0, false, nil)

	deleted, err := cql.Delete[models.Product](
		context.Background(),
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(1)),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(0), deleted)
}

func (ts *SoftDeleteIntTestSuite) TestSoftDeleteWhenAModelMatchConditions() {
	ts.createProduct("", 0, 0, false, nil)

	deleted, err := cql.Delete[models.Product](
		context.Background(),
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(0)),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(1), deleted)

	productReturned, err := cql.Query[models.Product](
		context.Background(),
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(1)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productReturned, 0)
}

func (ts *SoftDeleteIntTestSuite) TestSoftDeleteWhenMultipleModelsMatchConditions() {
	ts.createProduct("1", 0, 0, false, nil)
	ts.createProduct("2", 0, 0, false, nil)

	deleted, err := cql.Delete[models.Product](
		context.Background(),
		ts.db,
		conditions.Product.Bool.Is().False(),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(2), deleted)

	productReturned, err := cql.Query[models.Product](
		context.Background(),
		ts.db,
		conditions.Product.Bool.Is().False(),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productReturned, 0)
}

func (ts *SoftDeleteIntTestSuite) TestSoftDeleteWithJoinInConditions() {
	brand1 := ts.createBrand("google")
	brand2 := ts.createBrand("apple")

	ts.createPhone("pixel", *brand1)
	ts.createPhone("iphone", *brand2)

	deleted, err := cql.Delete[models.Phone](
		context.Background(),
		ts.db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("google")),
		),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(1), deleted)

	phones, err := cql.Query[models.Phone](
		context.Background(),
		ts.db,
		conditions.Phone.Name.Is().Eq(cql.String("pixel")),
	).Find()
	ts.Require().NoError(err)
	ts.Len(phones, 0)
}

func (ts *SoftDeleteIntTestSuite) TestSoftDeleteWithJoinDifferentEntitiesInConditions() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	ts.createSale(0, product1, seller1)
	ts.createSale(1, product2, seller2)
	ts.createSale(2, product1, seller2)
	ts.createSale(3, product2, seller1)

	deleted, err := cql.Delete[models.Sale](
		context.Background(),
		ts.db,
		conditions.Sale.Product(
			conditions.Product.Int.Is().Eq(cql.Int(1)),
		),
		conditions.Sale.Seller(
			conditions.Seller.Name.Is().Eq(cql.String("franco")),
		),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(1), deleted)

	sales, err := cql.Query[models.Sale](
		context.Background(),
		ts.db,
		conditions.Sale.Code.Is().Eq(cql.Int(0)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(sales, 0)
}

func (ts *SoftDeleteIntTestSuite) TestSoftDeleteWithMultilevelJoinInConditions() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("franco", company1)
	seller2 := ts.createSeller("agustin", company2)

	ts.createSale(0, product1, seller1)
	ts.createSale(1, product2, seller2)

	deleted, err := cql.Delete[models.Sale](
		context.Background(),
		ts.db,
		conditions.Sale.Seller(
			conditions.Seller.Name.Is().Eq(cql.String("franco")),
			conditions.Seller.Company(
				conditions.Company.Name.Is().Eq(cql.String("ditrit")),
			),
		),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(1), deleted)

	sales, err := cql.Query[models.Sale](
		context.Background(),
		ts.db,
		conditions.Sale.Code.Is().Eq(cql.Int(0)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(sales, 0)
}

func (ts *SoftDeleteIntTestSuite) TestSoftDeleteReturning() {
	switch getDBDialector() {
	// delete returning only supported for postgres, sqlite, sqlserver
	case sql.MySQL:
		_, err := cql.Delete[models.Phone](
			context.Background(),
			ts.db,
			conditions.Phone.Name.Is().Eq(cql.String("asd")),
		).Returning(nil).Exec()
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: Returning")
	case sql.Postgres, sql.SQLite, sql.SQLServer:
		product := ts.createProduct("", 0, 0, false, nil)

		productsReturned := []models.Product{}
		deleted, err := cql.Delete[models.Product](
			context.Background(),
			ts.db,
			conditions.Product.Int.Is().Eq(cql.Int(0)),
		).Returning(&productsReturned).Exec()
		ts.Require().NoError(err)
		ts.Equal(int64(1), deleted)

		ts.Len(productsReturned, 1)
		productReturned := productsReturned[0]
		ts.Equal(product.ID, productReturned.ID)

		products, err := cql.Query[models.Product](
			context.Background(),
			ts.db,
			conditions.Product.Int.Is().Eq(cql.Int(0)),
		).Find()
		ts.Require().NoError(err)
		ts.Len(products, 0)
	}
}

func (ts *SoftDeleteIntTestSuite) TestSoftDeleteReturningWithPreload() {
	salesReturned := []models.Sale{}
	_, err := cql.Delete[models.Sale](
		context.Background(),
		ts.db,
		conditions.Sale.Code.Is().Eq(cql.Int(0)),
		conditions.Sale.Product().Preload(),
	).Returning(&salesReturned).Exec()
	ts.ErrorIs(err, cql.ErrPreloadsInDeleteReturningNotAllowed)
	ts.ErrorContains(err, "preloads in delete returning are not allowed")
	ts.ErrorContains(err, "method: Returning")
}

func (ts *SoftDeleteIntTestSuite) TestSoftDeleteReturningWithPreloadCollection() {
	companiesReturned := []models.Company{}
	_, err := cql.Delete[models.Company](
		context.Background(),
		ts.db,
		conditions.Company.Name.Is().Eq(cql.String("ditrit")),
		conditions.Company.Sellers.Preload(),
	).Returning(&companiesReturned).Exec()
	ts.ErrorIs(err, cql.ErrPreloadsInDeleteReturningNotAllowed)
	ts.ErrorContains(err, "method: Returning")
}

func (ts *SoftDeleteIntTestSuite) TestSoftDeleteOrderByLimit() {
	// delete order by limit only supported for mysql
	if getDBDialector() != sql.MySQL {
		_, err := cql.Delete[models.Product](
			context.Background(),
			ts.db,
			conditions.Product.Bool.Is().False(),
		).Ascending(
			conditions.Product.String,
		).Limit(1).Exec()
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: Ascending")
	} else {
		product1 := ts.createProduct("1", 0, 0, false, nil)
		ts.createProduct("2", 0, 0, false, nil)

		deleted, err := cql.Delete[models.Product](
			context.Background(),
			ts.db,
			conditions.Product.Bool.Is().False(),
		).Descending(
			conditions.Product.String,
		).Limit(1).Exec()
		ts.Require().NoError(err)
		ts.Equal(int64(1), deleted)

		productReturned, err := cql.Query[models.Product](
			context.Background(),
			ts.db,
			conditions.Product.Int.Is().Eq(cql.Int(0)),
		).FindOne()
		ts.Require().NoError(err)

		ts.Equal(product1.ID, productReturned.ID)
	}
}

func (ts *SoftDeleteIntTestSuite) TestSoftDeleteLimitWithoutOrderByReturnsError() {
	// delete order by limit only supported for mysql
	if getDBDialector() != sql.MySQL {
		_, err := cql.Delete[models.Product](
			context.Background(),
			ts.db,
			conditions.Product.Bool.Is().False(),
		).Limit(1).Exec()
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: Limit")
	} else {
		_, err := cql.Delete[models.Product](
			context.Background(),
			ts.db,
			conditions.Product.Bool.Is().False(),
		).Limit(1).Exec()
		ts.ErrorIs(err, cql.ErrOrderByMustBeCalled)
		ts.ErrorContains(err, "method: Limit")
	}
}
