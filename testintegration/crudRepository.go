package testintegration

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"gotest.tools/assert"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type CRUDRepositoryIntTestSuite struct {
	suite.Suite
	db                    *gorm.DB
	crudProductRepository orm.CRUDRepository[models.Product, orm.UUID]
}

func NewCRUDRepositoryIntTestSuite(
	db *gorm.DB,
	crudProductRepository orm.CRUDRepository[models.Product, orm.UUID],
) *CRUDRepositoryIntTestSuite {
	return &CRUDRepositoryIntTestSuite{
		db:                    db,
		crudProductRepository: crudProductRepository,
	}
}

func (ts *CRUDRepositoryIntTestSuite) SetupTest() {
	CleanDB(ts.db)
}

func (ts *CRUDRepositoryIntTestSuite) TearDownSuite() {
	CleanDB(ts.db)
}

// ------------------------- GetByID --------------------------------

func (ts *CRUDRepositoryIntTestSuite) TestGetByIDReturnsErrorIfIDDontMatch() {
	ts.createProduct(0)
	_, err := ts.crudProductRepository.GetByID(ts.db, orm.NilUUID)
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *CRUDRepositoryIntTestSuite) TestGetByIDReturnsEntityIfIDMatch() {
	product := ts.createProduct(0)
	ts.createProduct(0)
	productReturned, err := ts.crudProductRepository.GetByID(ts.db, product.ID)
	ts.Nil(err)

	assert.DeepEqual(ts.T(), product, productReturned)
}

// ------------------------- QueryOne --------------------------------

func (ts *CRUDRepositoryIntTestSuite) TestQueryOneReturnsErrorIfConditionsDontMatch() {
	ts.createProduct(0)
	_, err := ts.crudProductRepository.QueryOne(
		ts.db,
		conditions.ProductInt(orm.Eq(1)),
	)
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *CRUDRepositoryIntTestSuite) TestQueryOneReturnsEntityIfConditionsMatch() {
	product := ts.createProduct(1)
	productReturned, err := ts.crudProductRepository.QueryOne(
		ts.db,
		conditions.ProductInt(orm.Eq(1)),
	)
	ts.Nil(err)

	assert.DeepEqual(ts.T(), product, productReturned)
}

func (ts *CRUDRepositoryIntTestSuite) TestQueryOneReturnsErrorIfMoreThanOneMatchConditions() {
	ts.createProduct(0)
	ts.createProduct(0)
	_, err := ts.crudProductRepository.QueryOne(
		ts.db,
		conditions.ProductInt(orm.Eq(0)),
	)
	ts.Error(err, orm.ErrMoreThanOneObjectFound)
}

// ------------------------- utils -------------------------

func (ts *CRUDRepositoryIntTestSuite) createProduct(intV int) *models.Product {
	entity := &models.Product{
		Int: intV,
	}
	err := ts.db.Create(entity).Error
	ts.Nil(err)

	return entity
}
