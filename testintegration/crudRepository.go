package testintegration

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"gotest.tools/assert"

	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/ditrit/badaas/testintegration/models"
)

type CRUDRepositoryIntTestSuite struct {
	suite.Suite
	db                    *gorm.DB
	crudProductRepository repository.CRUD[models.Product, model.UUID]
}

func NewCRUDRepositoryIntTestSuite(
	db *gorm.DB,
	crudProductRepository repository.CRUD[models.Product, model.UUID],
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
	_, err := ts.crudProductRepository.GetByID(ts.db, model.NilUUID)
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *CRUDRepositoryIntTestSuite) TestGetByIDReturnsEntityIfIDMatch() {
	product := ts.createProduct(0)
	ts.createProduct(0)
	productReturned, err := ts.crudProductRepository.GetByID(ts.db, product.ID)
	ts.Nil(err)

	assert.DeepEqual(ts.T(), product, productReturned)
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
