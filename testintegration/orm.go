package testintegration

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/testintegration/models"
)

type ORMIntTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (ts *ORMIntTestSuite) SetupTest() {
	CleanDB(ts.db)
}

func (ts *ORMIntTestSuite) TearDownSuite() {
	CleanDB(ts.db)
}

func (ts *ORMIntTestSuite) createProduct(stringV string, intV int, floatV float64, boolV bool, intP *int) *models.Product {
	entity := &models.Product{
		String:     stringV,
		Int:        intV,
		Float:      floatV,
		Bool:       boolV,
		IntPointer: intP,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *ORMIntTestSuite) createSale(code int, product *models.Product, seller *models.Seller) *models.Sale {
	entity := &models.Sale{
		Code:    code,
		Product: *product,
		Seller:  seller,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *ORMIntTestSuite) createSeller(name string, company *models.Company) *models.Seller {
	var companyID *model.UUID
	if company != nil {
		companyID = &company.ID
	}

	entity := &models.Seller{
		Name:      name,
		CompanyID: companyID,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *ORMIntTestSuite) createCompany(name string) *models.Company {
	entity := &models.Company{
		Name: name,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *ORMIntTestSuite) createCountry(name string, capital models.City) *models.Country {
	entity := &models.Country{
		Name:    name,
		Capital: capital,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *ORMIntTestSuite) createEmployee(name string, boss *models.Employee) *models.Employee {
	entity := &models.Employee{
		Name: name,
		Boss: boss,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *ORMIntTestSuite) createBicycle(name string, owner models.Person) *models.Bicycle {
	entity := &models.Bicycle{
		Name:  name,
		Owner: owner,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *ORMIntTestSuite) createBrand(name string) *models.Brand {
	entity := &models.Brand{
		Name: name,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *ORMIntTestSuite) createPhone(name string, brand models.Brand) *models.Phone {
	entity := &models.Phone{
		Name:  name,
		Brand: brand,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}

func (ts *ORMIntTestSuite) createUniversity(name string) *models.University {
	entity := &models.University{
		Name: name,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}
