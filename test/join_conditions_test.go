package test

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
	"github.com/FrancoLiberali/cql/unsafe"
)

type JoinConditionsIntTestSuite struct {
	testSuite
}

func NewJoinConditionsIntTestSuite(
	db *gorm.DB,
) *JoinConditionsIntTestSuite {
	return &JoinConditionsIntTestSuite{
		testSuite: testSuite{
			db: db,
		},
	}
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsUintBelongsTo() {
	brand1 := ts.createBrand("google")
	brand2 := ts.createBrand("apple")

	match := ts.createPhone("pixel", *brand1)
	ts.createPhone("iphone", *brand2)

	entities, err := cql.Query[models.Phone](
		ts.db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("google")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Phone{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsBelongsTo() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	match := ts.createSale(0, product1, nil)
	ts.createSale(0, product2, nil)

	entities, err := cql.Query[models.Sale](
		ts.db,
		conditions.Sale.Product(
			conditions.Product.Int.Is().Eq(cql.Int(1)),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsAndFiltersTheMainEntity() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(1, product1, seller1)
	ts.createSale(2, product2, seller2)
	ts.createSale(2, product1, seller2)

	entities, err := cql.Query[models.Sale](
		ts.db,
		conditions.Sale.Code.Is().Eq(cql.Int(1)),
		conditions.Sale.Product(
			conditions.Product.Int.Is().Eq(cql.Int(1)),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsHasOneOptional() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := cql.Query[models.Sale](
		ts.db,
		conditions.Sale.Seller(
			conditions.Seller.Name.Is().Eq(cql.String("franco")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsHasOneSelfReferential() {
	boss1 := &models.Employee{
		Name: "Xavier",
	}
	boss2 := &models.Employee{
		Name: "Vincent",
	}

	match := ts.createEmployee("franco", boss1)
	ts.createEmployee("pierre", boss2)

	entities, err := cql.Query[models.Employee](
		ts.db,
		conditions.Employee.Boss(
			conditions.Employee.Name.Is().Eq(cql.String("Xavier")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Employee{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsOneToOne() {
	capital1 := models.City{
		Name: "Buenos Aires",
	}
	capital2 := models.City{
		Name: "Paris",
	}

	ts.createCountry("Argentina", capital1)
	ts.createCountry("France", capital2)

	entities, err := cql.Query[models.City](
		ts.db,
		conditions.City.Country(
			conditions.Country.Name.Is().Eq(cql.String("Argentina")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.City{&capital1}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsOneToOneReversed() {
	capital1 := models.City{
		Name: "Buenos Aires",
	}
	capital2 := models.City{
		Name: "Paris",
	}

	country1 := ts.createCountry("Argentina", capital1)
	ts.createCountry("France", capital2)

	entities, err := cql.Query[models.Country](
		ts.db,
		conditions.Country.Capital(
			conditions.City.Name.Is().Eq(cql.String("Buenos Aires")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Country{country1}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsWithEntityThatDefinesTableName() {
	person1 := models.Person{
		Name: "franco",
	}
	person2 := models.Person{
		Name: "xavier",
	}

	match := ts.createBicycle("BMX", person1)
	ts.createBicycle("Shimano", person2)

	entities, err := cql.Query[models.Bicycle](
		ts.db,
		conditions.Bicycle.Owner(
			conditions.Person.Name.Is().Eq(cql.String("franco")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Bicycle{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsOnHasMany() {
	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	match := ts.createSeller("franco", company1)
	ts.createSeller("agustin", company2)

	entities, err := cql.Query[models.Seller](
		ts.db,
		conditions.Seller.Company(
			conditions.Company.Name.Is().Eq(cql.String("ditrit")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Seller{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsOnDifferentAttributes() {
	product1 := ts.createProduct("match", 1, 0.0, false, nil)
	product2 := ts.createProduct("match", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := cql.Query[models.Sale](
		ts.db,
		conditions.Sale.Product(
			conditions.Product.Int.Is().Eq(cql.Int(1)),
			conditions.Product.String.Is().Eq(cql.String("match")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsAddsDeletedAtAutomatically() {
	product1 := ts.createProduct("match", 1, 0.0, false, nil)
	product2 := ts.createProduct("match", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	ts.Nil(ts.db.Delete(product2).Error)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := cql.Query[models.Sale](
		ts.db,
		conditions.Sale.Product(
			conditions.Product.String.Is().Eq(cql.String("match")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsOnDeletedAt() {
	product1 := ts.createProduct("match", 1, 0.0, false, nil)
	product2 := ts.createProduct("match", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	ts.Nil(ts.db.Delete(product1).Error)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := cql.Query[models.Sale](
		ts.db,
		conditions.Sale.Product(
			conditions.Product.DeletedAt.Is().Eq(cql.Time(product1.DeletedAt.Time)),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsAndFiltersByNil() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	intProduct2 := 2
	product2 := ts.createProduct("", 2, 0.0, false, &intProduct2)

	match := ts.createSale(0, product1, nil)
	ts.createSale(0, product2, nil)

	entities, err := cql.Query[models.Sale](
		ts.db,
		conditions.Sale.Product(
			conditions.Product.IntPointer.Is().Null(),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsDifferentEntities() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)
	ts.createSale(0, product1, seller2)
	ts.createSale(0, product2, seller1)

	entities, err := cql.Query[models.Sale](
		ts.db,
		conditions.Sale.Product(
			conditions.Product.Int.Is().Eq(cql.Int(1)),
		),
		conditions.Sale.Seller(
			conditions.Seller.Name.Is().Eq(cql.String("franco")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsMultipleTimes() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("franco", company1)
	seller2 := ts.createSeller("agustin", company2)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := cql.Query[models.Sale](
		ts.db,
		conditions.Sale.Seller(
			conditions.Seller.Name.Is().Eq(cql.String("franco")),
			conditions.Seller.Company(
				conditions.Company.Name.Is().Eq(cql.String("ditrit")),
			),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestJoinWithUnsafeCondition() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("ditrit", company1)
	seller2 := ts.createSeller("agustin", company2)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := cql.Query[models.Sale](
		ts.db,
		conditions.Sale.Seller(
			conditions.Seller.Company(
				unsafe.NewCondition[models.Company]("%s.name = Seller.name"),
			),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorOver2Tables() {
	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("ditrit", company1)
	ts.createSeller("agustin", company2)

	entities, err := cql.Query[models.Seller](
		ts.db,
		conditions.Seller.Company(
			conditions.Company.Name.Is().Eq(conditions.Seller.Name),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Seller{seller1}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorOver2TablesAtMoreLevel() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("ditrit", company1)
	seller2 := ts.createSeller("agustin", company2)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := cql.Query[models.Sale](
		ts.db,
		conditions.Sale.Seller(
			conditions.Seller.Company(
				conditions.Company.Name.Is().Eq(conditions.Seller.Name),
			),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorWithNotJoinedModelReturnsError() {
	_, err := cql.Query[models.Child](
		ts.db,
		conditions.Child.ID.Is().Eq(conditions.ParentParent.ID),
	).Find()
	ts.ErrorIs(err, cql.ErrFieldModelNotConcerned)
	ts.ErrorContains(err, "not concerned model: models.ParentParent; operator: Eq; model: models.Child, field: ID")
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorWithJoinedInTheFutureModelReturnsError() {
	_, err := cql.Query[models.Child](
		ts.db,
		conditions.Child.ID.Is().Eq(conditions.Parent1.ID),
		conditions.Child.Parent1(),
	).Find()
	ts.ErrorIs(err, cql.ErrFieldModelNotConcerned)
	ts.ErrorContains(err, "not concerned model: models.Parent1; operator: Eq; model: models.Child, field: ID")
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorJoinMoreThanOnceWithoutAppearanceReturnsError() {
	_, err := cql.Query[models.Child](
		ts.db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
		conditions.Child.ID.Is().Eq(conditions.ParentParent.ID),
	).Find()
	ts.ErrorIs(err, cql.ErrAppearanceMustBeSelected)
	ts.ErrorContains(err, "model: models.ParentParent; operator: Eq; model: models.Child, field: ID")
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorJoinMoreThanOnceWithAppearance() {
	parentParent := &models.ParentParent{Name: "franco"}
	parent1 := &models.Parent1{ParentParent: *parentParent}
	parent2 := &models.Parent2{ParentParent: *parentParent}
	child := &models.Child{Parent1: *parent1, Parent2: *parent2, Name: "franco"}
	err := ts.db.Create(child).Error
	ts.Require().NoError(err)

	entities, err := cql.Query[models.Child](
		ts.db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
		conditions.Child.Name.Is().Eq(conditions.ParentParent.Name.Appearance(0)),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Child{child}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorJoinMoreThanOnceWithoutAppearanceOnMultivalueOperatorReturnsError() {
	_, err := cql.Query[models.Child](
		ts.db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
		conditions.Child.ID.Is().Between(
			conditions.ParentParent.ID,
			conditions.ParentParent.ID,
		),
	).Find()
	ts.ErrorIs(err, cql.ErrAppearanceMustBeSelected)
	ts.ErrorContains(err, "model: models.ParentParent; operator: Between; model: models.Child, field: ID")
}

func (ts *JoinConditionsIntTestSuite) TestCollectionAnyReturnsEmptyWhenNothingMatches() {
	company1 := ts.createCompany("ditrit")
	ts.createCompany("orness")

	ts.createSeller("franco", company1)
	ts.createSeller("agustin", company1)

	entities, err := cql.Query[models.Company](
		ts.db,
		conditions.Company.Sellers.Any(
			conditions.Seller.Name.Is().Eq(cql.String("not")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Company{}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestCollectionAnyReturnsIfOneMatches() {
	company1 := ts.createCompany("ditrit")
	ts.createCompany("orness")

	ts.createSeller("franco", company1)
	ts.createSeller("agustin", company1)

	entities, err := cql.Query[models.Company](
		ts.db,
		conditions.Company.Sellers.Any(
			conditions.Seller.Name.Is().Eq(cql.String("franco")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Company{company1}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestCollectionAnyReturnsIfMoreThanOneMatches() {
	company1 := ts.createCompany("ditrit")
	ts.createCompany("orness")

	ts.createSeller("franco", company1)
	ts.createSeller("agustin", company1)

	entities, err := cql.Query[models.Company](
		ts.db,
		conditions.Company.Sellers.Any(
			cql.Or(
				conditions.Seller.Name.Is().Eq(cql.String("franco")),
				conditions.Seller.Name.Is().Eq(cql.String("agustin")),
			),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Company{company1}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestCollectionNoneReturnsWhenIsEmpty() {
	company1 := ts.createCompany("ditrit")

	entities, err := cql.Query[models.Company](
		ts.db,
		conditions.Company.Sellers.None(
			conditions.Seller.Name.Is().Eq(cql.String("not")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Company{company1}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestCollectionNoneReturnsWhenNothingMatches() {
	company1 := ts.createCompany("ditrit")

	ts.createSeller("franco", company1)
	ts.createSeller("agustin", company1)

	entities, err := cql.Query[models.Company](
		ts.db,
		conditions.Company.Sellers.None(
			conditions.Seller.Name.Is().Eq(cql.String("not")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Company{company1}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestCollectionNoneReturnsEmptyIfOneMatches() {
	company1 := ts.createCompany("ditrit")

	ts.createSeller("franco", company1)
	ts.createSeller("agustin", company1)

	entities, err := cql.Query[models.Company](
		ts.db,
		conditions.Company.Sellers.None(
			conditions.Seller.Name.Is().Eq(cql.String("franco")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Company{}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestCollectionNoneReturnsEmptyIfMoreThanOneMatches() {
	company1 := ts.createCompany("ditrit")

	ts.createSeller("franco", company1)
	ts.createSeller("agustin", company1)

	entities, err := cql.Query[models.Company](
		ts.db,
		conditions.Company.Sellers.None(
			cql.Or(
				conditions.Seller.Name.Is().Eq(cql.String("franco")),
				conditions.Seller.Name.Is().Eq(cql.String("agustin")),
			),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Company{}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestCollectionAllReturnsWhenIsEmpty() {
	company1 := ts.createCompany("ditrit")

	entities, err := cql.Query[models.Company](
		ts.db,
		conditions.Company.Sellers.All(
			conditions.Seller.Name.Is().Eq(cql.String("not")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Company{company1}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestCollectionAllReturnsEmptyWhenNothingMatches() {
	company1 := ts.createCompany("ditrit")

	ts.createSeller("franco", company1)
	ts.createSeller("agustin", company1)

	entities, err := cql.Query[models.Company](
		ts.db,
		conditions.Company.Sellers.All(
			conditions.Seller.Name.Is().Eq(cql.String("not")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Company{}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestCollectionReturnsEmptyIfOneMatches() {
	company1 := ts.createCompany("ditrit")

	ts.createSeller("franco", company1)
	ts.createSeller("agustin", company1)

	entities, err := cql.Query[models.Company](
		ts.db,
		conditions.Company.Sellers.All(
			conditions.Seller.Name.Is().Eq(cql.String("franco")),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Company{}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestCollectionAllReturnsIfAllMatch() {
	company1 := ts.createCompany("ditrit")

	ts.createSeller("franco", company1)
	ts.createSeller("agustin", company1)

	entities, err := cql.Query[models.Company](
		ts.db,
		conditions.Company.Sellers.All(
			cql.Or(
				conditions.Seller.Name.Is().Eq(cql.String("franco")),
				conditions.Seller.Name.Is().Eq(cql.String("agustin")),
			),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Company{company1}, entities)
}
