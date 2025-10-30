package test

import (
	"github.com/stretchr/testify/suite"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/test/models"
)

type testSuite struct {
	suite.Suite
	db *cql.DB
}

func (ts *testSuite) SetupTest() {
	CleanDB(ts.db)
}

func create[T any](ts *testSuite, entity T) T {
	err := ts.db.GormDB.Create(entity).Error
	ts.Require().NoError(err)

	return entity
}

func (ts *testSuite) createProduct(stringV string, intV int, floatV float64, boolV bool, intP *int) *models.Product {
	return create(ts, &models.Product{
		String:     stringV,
		Int:        intV,
		Float:      floatV,
		Bool:       boolV,
		IntPointer: intP,
	})
}

func (ts *testSuite) createProductNoTimestamps(stringV string, intV int, floatV float64, boolV bool, intP *int) *models.ProductNoTimestamps {
	return create(ts, &models.ProductNoTimestamps{
		String:     stringV,
		Int:        intV,
		Float:      floatV,
		Bool:       boolV,
		IntPointer: intP,
	})
}

func (ts *testSuite) createSale(code int, product *models.Product, seller *models.Seller) *models.Sale {
	return create(ts, &models.Sale{
		Code:    code,
		Product: *product,
		Seller:  seller,
	})
}

func (ts *testSuite) createSaleNoTimestamps(code int, product *models.ProductNoTimestamps, seller *models.SellerNoTimestamps) *models.SaleNoTimestamps {
	return create(ts, &models.SaleNoTimestamps{
		Code:    code,
		Product: *product,
		Seller:  seller,
	})
}

func (ts *testSuite) createSeller(name string, company *models.Company) *models.Seller {
	var companyID *model.UUID
	if company != nil {
		companyID = &company.ID
	}

	return create(ts, &models.Seller{
		Name:      name,
		CompanyID: companyID,
	})
}

func (ts *testSuite) createSellerNoTimestamps(name string, company *models.CompanyNoTimestamps) *models.SellerNoTimestamps {
	var companyID *model.UUID
	if company != nil {
		companyID = &company.ID
	}

	return create(ts, &models.SellerNoTimestamps{
		Name:                  name,
		CompanyNoTimestampsID: companyID,
	})
}

func (ts *testSuite) createCompany(name string) *models.Company {
	return create(ts, &models.Company{
		Name: name,
	})
}

func (ts *testSuite) createCompanyNoTimestamps(name string) *models.CompanyNoTimestamps {
	return create(ts, &models.CompanyNoTimestamps{
		Name: name,
	})
}

func (ts *testSuite) createCountry(name string, capital models.City) *models.Country {
	return create(ts, &models.Country{
		Name:    name,
		Capital: capital,
	})
}

func (ts *testSuite) createEmployee(name string, boss *models.Employee) *models.Employee {
	return create(ts, &models.Employee{
		Name: name,
		Boss: boss,
	})
}

func (ts *testSuite) createBicycle(name string, owner models.Person) *models.Bicycle {
	return create(ts, &models.Bicycle{
		Name:  name,
		Owner: owner,
	})
}

func (ts *testSuite) createBrand(name string) *models.Brand {
	return create(ts, &models.Brand{
		Name: name,
	})
}

func (ts *testSuite) createPhone(name string, brand models.Brand) *models.Phone {
	return create(ts, &models.Phone{
		Name:  name,
		Brand: brand,
	})
}

func (ts *testSuite) createPhoneNoTimestamps(name string, brand models.Brand) *models.PhoneNoTimestamps {
	return create(ts, &models.PhoneNoTimestamps{
		Name:  name,
		Brand: brand,
	})
}

func (ts *testSuite) createUniversity(name string) *models.University {
	return create(ts, &models.University{
		Name: name,
	})
}
