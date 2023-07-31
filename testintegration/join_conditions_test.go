package testintegration

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type JoinConditionsIntTestSuite struct {
	CRUDServiceCommonIntTestSuite
	crudSaleService     orm.CRUDService[models.Sale, orm.UUID]
	crudSellerService   orm.CRUDService[models.Seller, orm.UUID]
	crudCountryService  orm.CRUDService[models.Country, orm.UUID]
	crudCityService     orm.CRUDService[models.City, orm.UUID]
	crudEmployeeService orm.CRUDService[models.Employee, orm.UUID]
	crudBicycleService  orm.CRUDService[models.Bicycle, orm.UUID]
	crudPhoneService    orm.CRUDService[models.Phone, uint]
}

func NewJoinConditionsIntTestSuite(
	db *gorm.DB,
	crudSaleService orm.CRUDService[models.Sale, orm.UUID],
	crudSellerService orm.CRUDService[models.Seller, orm.UUID],
	crudCountryService orm.CRUDService[models.Country, orm.UUID],
	crudCityService orm.CRUDService[models.City, orm.UUID],
	crudEmployeeService orm.CRUDService[models.Employee, orm.UUID],
	crudBicycleService orm.CRUDService[models.Bicycle, orm.UUID],
	crudPhoneService orm.CRUDService[models.Phone, uint],
) *JoinConditionsIntTestSuite {
	return &JoinConditionsIntTestSuite{
		CRUDServiceCommonIntTestSuite: CRUDServiceCommonIntTestSuite{
			db: db,
		},
		crudSaleService:     crudSaleService,
		crudSellerService:   crudSellerService,
		crudCountryService:  crudCountryService,
		crudCityService:     crudCityService,
		crudEmployeeService: crudEmployeeService,
		crudBicycleService:  crudBicycleService,
		crudPhoneService:    crudPhoneService,
	}
}

func (ts *JoinConditionsIntTestSuite) TestQueryWithConditionThatJoinsUintBelongsTo() {
	brand1 := ts.createBrand("google")
	brand2 := ts.createBrand("apple")

	match := ts.createPhone("pixel", *brand1)
	ts.createPhone("iphone", *brand2)

	entities, err := ts.crudPhoneService.Query(
		conditions.PhoneBrand(
			conditions.BrandName("google"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Phone{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestQueryWithConditionThatJoinsBelongsTo() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	match := ts.createSale(0, product1, nil)
	ts.createSale(0, product2, nil)

	entities, err := ts.crudSaleService.Query(
		conditions.SaleProduct(
			conditions.ProductInt(1),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestQueryWithConditionThatJoinsAndFiltersTheMainEntity() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(1, product1, seller1)
	ts.createSale(2, product2, seller2)
	ts.createSale(2, product1, seller2)

	entities, err := ts.crudSaleService.Query(
		conditions.SaleCode(1),
		conditions.SaleProduct(
			conditions.ProductInt(1),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestQueryWithConditionThatJoinsHasOneOptional() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.Query(
		conditions.SaleSeller(
			conditions.SellerName("franco"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestQueryWithConditionThatJoinsHasOneSelfReferential() {
	boss1 := &models.Employee{
		Name: "Xavier",
	}
	boss2 := &models.Employee{
		Name: "Vincent",
	}

	match := ts.createEmployee("franco", boss1)
	ts.createEmployee("pierre", boss2)

	entities, err := ts.crudEmployeeService.Query(
		conditions.EmployeeBoss(
			conditions.EmployeeName("Xavier"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Employee{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestQueryWithConditionThatJoinsOneToOne() {
	capital1 := models.City{
		Name: "Buenos Aires",
	}
	capital2 := models.City{
		Name: "Paris",
	}

	ts.createCountry("Argentina", capital1)
	ts.createCountry("France", capital2)

	entities, err := ts.crudCityService.Query(
		conditions.CityCountry(
			conditions.CountryName("Argentina"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.City{&capital1}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestQueryWithConditionThatJoinsOneToOneReversed() {
	capital1 := models.City{
		Name: "Buenos Aires",
	}
	capital2 := models.City{
		Name: "Paris",
	}

	country1 := ts.createCountry("Argentina", capital1)
	ts.createCountry("France", capital2)

	entities, err := ts.crudCountryService.Query(
		conditions.CountryCapital(
			conditions.CityName("Buenos Aires"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Country{country1}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestQueryWithConditionThatJoinsWithEntityThatDefinesTableName() {
	person1 := models.Person{
		Name: "franco",
	}
	person2 := models.Person{
		Name: "xavier",
	}

	match := ts.createBicycle("BMX", person1)
	ts.createBicycle("Shimano", person2)

	entities, err := ts.crudBicycleService.Query(
		conditions.BicycleOwner(
			conditions.PersonName("franco"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Bicycle{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestQueryWithConditionThatJoinsOnHasMany() {
	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	match := ts.createSeller("franco", company1)
	ts.createSeller("agustin", company2)

	entities, err := ts.crudSellerService.Query(
		conditions.SellerCompany(
			conditions.CompanyName("ditrit"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Seller{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestQueryWithConditionThatJoinsOnDifferentAttributes() {
	product1 := ts.createProduct("match", 1, 0.0, false, nil)
	product2 := ts.createProduct("match", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.Query(
		conditions.SaleProduct(
			conditions.ProductInt(1),
			conditions.ProductString("match"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestQueryWithConditionThatJoinsAddsDeletedAtAutomatically() {
	product1 := ts.createProduct("match", 1, 0.0, false, nil)
	product2 := ts.createProduct("match", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	ts.Nil(ts.db.Delete(product2).Error)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.Query(
		conditions.SaleProduct(
			conditions.ProductString("match"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestQueryWithConditionThatJoinsOnDeletedAt() {
	product1 := ts.createProduct("match", 1, 0.0, false, nil)
	product2 := ts.createProduct("match", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	ts.Nil(ts.db.Delete(product1).Error)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.Query(
		conditions.SaleProduct(
			conditions.ProductDeletedAt(product1.DeletedAt),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestQueryWithConditionThatJoinsDifferentEntities() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)
	ts.createSale(0, product1, seller2)
	ts.createSale(0, product2, seller1)

	entities, err := ts.crudSaleService.Query(
		conditions.SaleProduct(
			conditions.ProductInt(1),
		),
		conditions.SaleSeller(
			conditions.SellerName("franco"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestQueryWithConditionThatJoinsMultipleTimes() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("franco", company1)
	seller2 := ts.createSeller("agustin", company2)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.Query(
		conditions.SaleSeller(
			conditions.SellerName("franco"),
			conditions.SellerCompany(
				conditions.CompanyName("ditrit"),
			),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}
