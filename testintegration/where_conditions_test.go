package testintegration

import (
	"gorm.io/gorm"
	"gotest.tools/assert"

	"github.com/google/uuid"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type WhereConditionsIntTestSuite struct {
	CRUDServiceCommonIntTestSuite
	crudProductService orm.CRUDService[models.Product, orm.UUID]
	crudSaleService    orm.CRUDService[models.Sale, orm.UUID]
	crudBrandService   orm.CRUDService[models.Brand, uint]
}

func NewWhereConditionsIntTestSuite(
	db *gorm.DB,
	crudProductService orm.CRUDService[models.Product, orm.UUID],
	crudSaleService orm.CRUDService[models.Sale, orm.UUID],
	crudBrandService orm.CRUDService[models.Brand, uint],
) *WhereConditionsIntTestSuite {
	return &WhereConditionsIntTestSuite{
		CRUDServiceCommonIntTestSuite: CRUDServiceCommonIntTestSuite{
			db: db,
		},
		crudProductService: crudProductService,
		crudSaleService:    crudSaleService,
		crudBrandService:   crudBrandService,
	}
}

// ------------------------- GetByID --------------------------------

func (ts *WhereConditionsIntTestSuite) TestGetByIDReturnsErrorIfNotEntityCreated() {
	_, err := ts.crudProductService.GetByID(orm.NilUUID)
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *WhereConditionsIntTestSuite) TestGetByIDReturnsErrorIfNotEntityMatch() {
	ts.createProduct("", 0, 0, false, nil)

	_, err := ts.crudProductService.GetByID(orm.UUID(uuid.New()))
	ts.Error(err, gorm.ErrRecordNotFound)
}

func (ts *WhereConditionsIntTestSuite) TestGetByIDReturnsTheEntityIfItIsCreate() {
	match := ts.createProduct("", 0, 0, false, nil)

	entity, err := ts.crudProductService.GetByID(match.ID)
	ts.Nil(err)

	assert.DeepEqual(ts.T(), match, entity)
}

// ------------------------- Query --------------------------------

func (ts *WhereConditionsIntTestSuite) TestQueryWithoutConditionsReturnsEmptyIfNotEntitiesCreated() {
	entities, err := ts.crudProductService.Query()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithoutConditionsReturnsTheOnlyOneIfOneEntityCreated() {
	match := ts.createProduct("", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithoutConditionsReturnsTheListWhenMultipleCreated() {
	match1 := ts.createProduct("", 0, 0, false, nil)
	match2 := ts.createProduct("", 0, 0, false, nil)
	match3 := ts.createProduct("", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionsReturnsEmptyIfNotEntitiesCreated() {
	entities, err := ts.crudProductService.Query(
		conditions.ProductString("not_created"),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionsReturnsEmptyIfNothingMatch() {
	ts.createProduct("something_else", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductString("not_match"),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionsReturnsOneIfOnlyOneMatch() {
	match := ts.createProduct("match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductString("match"),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionsReturnsMultipleIfMultipleMatch() {
	match1 := ts.createProduct("match", 0, 0, false, nil)
	match2 := ts.createProduct("match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductString("match"),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionOfIntType() {
	match := ts.createProduct("match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(1),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionOfFloatType() {
	match := ts.createProduct("match", 0, 1.1, false, nil)
	ts.createProduct("not_match", 0, 2.2, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductFloat(1.1),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionOfBoolType() {
	match := ts.createProduct("match", 0, 0.0, true, nil)
	ts.createProduct("not_match", 0, 0.0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductBool(true),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithMultipleConditionsOfDifferentTypesWorks() {
	match1 := ts.createProduct("match", 1, 0.0, true, nil)
	match2 := ts.createProduct("match", 1, 0.0, true, nil)

	ts.createProduct("not_match", 1, 0.0, true, nil)
	ts.createProduct("match", 2, 0.0, true, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductString("match"),
		conditions.ProductInt(1),
		conditions.ProductBool(true),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionOfID() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductId(match.ID),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionOfCreatedAt() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductCreatedAt(match.CreatedAt),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryDeletedAtConditionIsAddedAutomatically() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	deleted := ts.createProduct("", 0, 0.0, false, nil)

	ts.Nil(ts.db.Delete(deleted).Error)

	entities, err := ts.crudProductService.Query()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

// TODO DeletedAt with nil value but not automatic

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionOfDeletedAtNotNil() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	ts.Nil(ts.db.Delete(match).Error)

	entities, err := ts.crudProductService.Query(
		conditions.ProductDeletedAt(match.DeletedAt),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionOfEmbedded() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)
	match.EmbeddedInt = 1

	err := ts.db.Save(match).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.Query(
		conditions.ProductEmbeddedInt(1),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionOfGormEmbedded() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)
	match.GormEmbedded.Int = 1

	err := ts.db.Save(match).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.Query(
		conditions.ProductGormEmbeddedInt(1),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionOfPointerTypeWithValue() {
	intMatch := 1
	match := ts.createProduct("match", 1, 0, false, &intMatch)
	intNotMatch := 2
	ts.createProduct("not_match", 2, 0, false, &intNotMatch)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductIntPointer(intMatch),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionOfByteArrayWithContent() {
	match := ts.createProduct("match", 1, 0, false, nil)
	notMatch1 := ts.createProduct("not_match", 2, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)
	match.ByteArray = []byte{1, 2}
	notMatch1.ByteArray = []byte{2, 3}

	err := ts.db.Save(match).Error
	ts.Nil(err)

	err = ts.db.Save(notMatch1).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.Query(
		conditions.ProductByteArray([]byte{1, 2}),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionOfByteArrayEmpty() {
	match := ts.createProduct("match", 1, 0, false, nil)
	notMatch1 := ts.createProduct("not_match", 2, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)
	match.ByteArray = []byte{}
	notMatch1.ByteArray = []byte{2, 3}

	err := ts.db.Save(match).Error
	ts.Nil(err)

	err = ts.db.Save(notMatch1).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.Query(
		conditions.ProductByteArray([]byte{}),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionOfCustomType() {
	match := ts.createProduct("match", 1, 0, false, nil)
	notMatch1 := ts.createProduct("not_match", 2, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)
	match.MultiString = models.MultiString{"salut", "hola"}
	notMatch1.MultiString = models.MultiString{"salut", "hola", "hello"}

	err := ts.db.Save(match).Error
	ts.Nil(err)

	err = ts.db.Save(notMatch1).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.Query(
		conditions.ProductMultiString(models.MultiString{"salut", "hola"}),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionOfRelationType() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.Query(
		conditions.SaleProductId(product1.ID),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionOfRelationTypeOptionalWithValue() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := ts.crudSaleService.Query(
		conditions.SaleSellerId(seller1.ID),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryWithConditionsOnUIntModel() {
	match := ts.createBrand("match")
	ts.createBrand("not_match")

	entities, err := ts.crudBrandService.Query(
		conditions.BrandName("match"),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Brand{match}, entities)
}
