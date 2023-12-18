package testintegration

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/errors"
	"github.com/ditrit/badaas/orm/unsafe"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type JoinConditionsIntTestSuite struct {
	ORMIntTestSuite
}

func NewJoinConditionsIntTestSuite(
	db *gorm.DB,
) *JoinConditionsIntTestSuite {
	return &JoinConditionsIntTestSuite{
		ORMIntTestSuite: ORMIntTestSuite{
			db: db,
		},
	}
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsUintBelongsTo() {
	brand1 := ts.createBrand("google")
	brand2 := ts.createBrand("apple")

	match := ts.createPhone("pixel", *brand1)
	ts.createPhone("iphone", *brand2)

	entities, err := orm.NewQuery[models.Phone](
		ts.db,
		conditions.Phone.Brand(
			conditions.Brand.NameIs().Eq("google"),
		),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Phone{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsBelongsTo() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	match := ts.createSale(0, product1, nil)
	ts.createSale(0, product2, nil)

	entities, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.Product(
			conditions.Product.IntIs().Eq(1),
		),
	).Find()
	ts.Nil(err)

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

	entities, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.CodeIs().Eq(1),
		conditions.Sale.Product(
			conditions.Product.IntIs().Eq(1),
		),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsHasOneOptional() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.Seller(
			conditions.Seller.NameIs().Eq("franco"),
		),
	).Find()
	ts.Nil(err)

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

	entities, err := orm.NewQuery[models.Employee](
		ts.db,
		conditions.Employee.Boss(
			conditions.Employee.NameIs().Eq("Xavier"),
		),
	).Find()
	ts.Nil(err)

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

	entities, err := orm.NewQuery[models.City](
		ts.db,
		conditions.City.Country(
			conditions.Country.NameIs().Eq("Argentina"),
		),
	).Find()
	ts.Nil(err)

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

	entities, err := orm.NewQuery[models.Country](
		ts.db,
		conditions.Country.Capital(
			conditions.City.NameIs().Eq("Buenos Aires"),
		),
	).Find()
	ts.Nil(err)

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

	entities, err := orm.NewQuery[models.Bicycle](
		ts.db,
		conditions.Bicycle.Owner(
			conditions.Person.NameIs().Eq("franco"),
		),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Bicycle{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsOnHasMany() {
	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	match := ts.createSeller("franco", company1)
	ts.createSeller("agustin", company2)

	entities, err := orm.NewQuery[models.Seller](
		ts.db,
		conditions.Seller.Company(
			conditions.Company.NameIs().Eq("ditrit"),
		),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Seller{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsOnDifferentAttributes() {
	product1 := ts.createProduct("match", 1, 0.0, false, nil)
	product2 := ts.createProduct("match", 2, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.Product(
			conditions.Product.IntIs().Eq(1),
			conditions.Product.StringIs().Eq("match"),
		),
	).Find()
	ts.Nil(err)

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

	entities, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.Product(
			conditions.Product.StringIs().Eq("match"),
		),
	).Find()
	ts.Nil(err)

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

	entities, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.Product(
			conditions.Product.DeletedAtIs().Eq(product1.DeletedAt.Time),
		),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestConditionThatJoinsAndFiltersByNil() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	intProduct2 := 2
	product2 := ts.createProduct("", 2, 0.0, false, &intProduct2)

	match := ts.createSale(0, product1, nil)
	ts.createSale(0, product2, nil)

	entities, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.Product(
			conditions.Product.IntPointerIs().Null(),
		),
	).Find()
	ts.Nil(err)

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

	entities, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.Product(
			conditions.Product.IntIs().Eq(1),
		),
		conditions.Sale.Seller(
			conditions.Seller.NameIs().Eq("franco"),
		),
	).Find()
	ts.Nil(err)

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

	entities, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.Seller(
			conditions.Seller.NameIs().Eq("franco"),
			conditions.Seller.Company(
				conditions.Company.NameIs().Eq("ditrit"),
			),
		),
	).Find()
	ts.Nil(err)

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

	entities, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.Seller(
			conditions.Seller.Company(
				unsafe.NewCondition[models.Company]("%s.name = Seller.name"),
			),
		),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestJoinWithEmptyConnectionConditionMakesNothing() {
	product1 := ts.createProduct("", 1, 0.0, false, nil)
	product2 := ts.createProduct("", 2, 0.0, false, nil)

	match1 := ts.createSale(0, product1, nil)
	match2 := ts.createSale(0, product2, nil)

	entities, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.Product(
			orm.And[models.Product](),
		),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match1, match2}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestJoinWithEmptyContainerConditionReturnsError() {
	_, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.Product(
			orm.Not[models.Product](),
		),
	).Find()
	ts.ErrorIs(err, errors.ErrEmptyConditions)
	ts.ErrorContains(err, "connector: Not; model: models.Product")
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorOver2Tables() {
	company1 := ts.createCompany("ditrit")
	company2 := ts.createCompany("orness")

	seller1 := ts.createSeller("ditrit", company1)
	ts.createSeller("agustin", company2)

	entities, err := orm.NewQuery[models.Seller](
		ts.db,
		conditions.Seller.Company(
			conditions.Company.NameIs().Dynamic().Eq(conditions.Seller.Name),
		),
	).Find()
	ts.Nil(err)

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

	entities, err := orm.NewQuery[models.Sale](
		ts.db,
		conditions.Sale.Seller(
			conditions.Seller.Company(
				conditions.Company.NameIs().Dynamic().Eq(conditions.Seller.Name),
			),
		),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorWithNotJoinedModelReturnsError() {
	_, err := orm.NewQuery[models.Child](
		ts.db,
		conditions.Child.IdIs().Dynamic().Eq(conditions.ParentParent.ID),
	).Find()
	ts.ErrorIs(err, errors.ErrFieldModelNotConcerned)
	ts.ErrorContains(err, "not concerned model: models.ParentParent; operator: Eq; model: models.Child, field: ID")
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorJoinMoreThanOnceWithoutSelectJoinReturnsError() {
	_, err := orm.NewQuery[models.Child](
		ts.db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
		conditions.Child.IdIs().Dynamic().Eq(conditions.ParentParent.ID),
	).Find()
	ts.ErrorIs(err, errors.ErrJoinMustBeSelected)
	ts.ErrorContains(err, "joined multiple times model: models.ParentParent; operator: Eq; model: models.Child, field: ID")
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorJoinMoreThanOnceWithSelectJoin() {
	parentParent := &models.ParentParent{Name: "franco"}
	parent1 := &models.Parent1{ParentParent: *parentParent}
	parent2 := &models.Parent2{ParentParent: *parentParent}
	child := &models.Child{Parent1: *parent1, Parent2: *parent2, Name: "franco"}
	err := ts.db.Create(child).Error
	ts.Nil(err)

	entities, err := orm.NewQuery[models.Child](
		ts.db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
		conditions.Child.NameIs().Dynamic().Eq(conditions.ParentParent.Name).SelectJoin(0, 0),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Child{child}, entities)
}

func (ts *JoinConditionsIntTestSuite) TestDynamicOperatorJoinMoreThanOnceWithoutSelectJoinOnMultivalueOperatorReturnsError() {
	_, err := orm.NewQuery[models.Child](
		ts.db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
		conditions.Child.IdIs().Dynamic().Between(
			conditions.ParentParent.ID,
			conditions.ParentParent.ID,
		),
	).Find()
	ts.ErrorIs(err, errors.ErrJoinMustBeSelected)
	ts.ErrorContains(err, "joined multiple times model: models.ParentParent; operator: Between; model: models.Child, field: ID")
}
