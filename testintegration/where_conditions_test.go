package testintegration

import (
	"log"

	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/cql"
	"github.com/ditrit/badaas/orm/mysql"
	"github.com/ditrit/badaas/orm/unsafe"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type WhereConditionsIntTestSuite struct {
	ORMIntTestSuite
}

func NewWhereConditionsIntTestSuite(
	db *gorm.DB,
) *WhereConditionsIntTestSuite {
	return &WhereConditionsIntTestSuite{
		ORMIntTestSuite: ORMIntTestSuite{
			db: db,
		},
	}
}

func (ts *WhereConditionsIntTestSuite) TestQueryReturnsEmptyIfNotEntitiesCreated() {
	entities, err := orm.Query[models.Product](ts.db).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryReturnsTheOnlyOneIfOneEntityCreated() {
	match := ts.createProduct("", 0, 0, false, nil)

	entities, err := orm.Query[models.Product](ts.db).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestQueryReturnsTheListWhenMultipleCreated() {
	match1 := ts.createProduct("", 0, 0, false, nil)
	match2 := ts.createProduct("", 0, 0, false, nil)
	match3 := ts.createProduct("", 0, 0, false, nil)

	entities, err := orm.Query[models.Product](ts.db).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionsReturnsEmptyIfNotEntitiesCreated() {
	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.String.Is().Eq("not_created"),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionsReturnsEmptyIfNothingMatch() {
	ts.createProduct("something_else", 0, 0, false, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.String.Is().Eq("not_match"),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionsReturnsOneIfOnlyOneMatch() {
	match := ts.createProduct("match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.String.Is().Eq("match"),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionsReturnsMultipleIfMultipleMatch() {
	match1 := ts.createProduct("match", 0, 0, false, nil)
	match2 := ts.createProduct("match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.String.Is().Eq("match"),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfIntType() {
	match := ts.createProduct("match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(1),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfFloatType() {
	match := ts.createProduct("match", 0, 1.1, false, nil)
	ts.createProduct("not_match", 0, 2.2, false, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.Float.Is().Eq(1.1),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfBoolType() {
	match := ts.createProduct("match", 0, 0.0, true, nil)
	ts.createProduct("not_match", 0, 0.0, false, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.Bool.Is().True(),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestMultipleConditionsOfDifferentTypesWorks() {
	match1 := ts.createProduct("match", 1, 0.0, true, nil)
	match2 := ts.createProduct("match", 1, 0.0, true, nil)

	ts.createProduct("not_match", 1, 0.0, true, nil)
	ts.createProduct("match", 2, 0.0, true, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.String.Is().Eq("match"),
		conditions.Product.Int.Is().Eq(1),
		conditions.Product.Bool.Is().True(),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfID() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.ID.Is().Eq(match.ID),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfCreatedAt() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.CreatedAt.Is().Eq(match.CreatedAt),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestDeletedAtConditionIsAddedAutomatically() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	deleted := ts.createProduct("", 0, 0.0, false, nil)

	ts.Nil(ts.db.Delete(deleted).Error)

	entities, err := orm.Query[models.Product](ts.db).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfDeletedAt() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	ts.Nil(ts.db.Delete(match).Error)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.DeletedAt.Is().Eq(match.DeletedAt.Time),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfEmbedded() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	match.EmbeddedInt = 1

	err := ts.db.Save(match).Error
	ts.Nil(err)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.ToBeEmbeddedEmbeddedInt.Is().Eq(1),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfGormEmbedded() {
	match := ts.createProduct("", 0, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	match.GormEmbedded.Int = 1

	err := ts.db.Save(match).Error
	ts.Nil(err)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.GormEmbeddedInt.Is().Eq(1),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfPointerTypeWithValue() {
	intMatch := 1
	match := ts.createProduct("match", 1, 0, false, &intMatch)
	intNotMatch := 2
	ts.createProduct("not_match", 2, 0, false, &intNotMatch)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.IntPointer.Is().Eq(1),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfByteArrayWithContent() {
	match := ts.createProduct("match", 1, 0, false, nil)
	notMatch1 := ts.createProduct("not_match", 2, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	match.ByteArray = []byte{1, 2}
	notMatch1.ByteArray = []byte{2, 3}

	err := ts.db.Save(match).Error
	ts.Nil(err)

	err = ts.db.Save(notMatch1).Error
	ts.Nil(err)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.ByteArray.Is().Eq([]byte{1, 2}),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfByteArrayEmpty() {
	match := ts.createProduct("match", 1, 0, false, nil)
	notMatch1 := ts.createProduct("not_match", 2, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	match.ByteArray = []byte{}
	notMatch1.ByteArray = []byte{2, 3}

	err := ts.db.Save(match).Error
	ts.Nil(err)

	err = ts.db.Save(notMatch1).Error
	ts.Nil(err)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.ByteArray.Is().Eq([]byte{}),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfCustomType() {
	match := ts.createProduct("match", 1, 0, false, nil)
	notMatch1 := ts.createProduct("not_match", 2, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	match.MultiString = models.MultiString{"salut", "hola"}
	notMatch1.MultiString = models.MultiString{"salut", "hola", "hello"}

	err := ts.db.Save(match).Error
	ts.Nil(err)

	err = ts.db.Save(notMatch1).Error
	ts.Nil(err)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.MultiString.Is().Eq(models.MultiString{"salut", "hola"}),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfRelationType() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := orm.Query[models.Sale](
		ts.db,
		conditions.Sale.ProductID.Is().Eq(product1.ID),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfRelationTypeOptionalWithValue() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	seller1 := ts.createSeller("franco", nil)
	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, seller1)
	ts.createSale(0, product2, seller2)

	entities, err := orm.Query[models.Sale](
		ts.db,
		conditions.Sale.SellerID.Is().Eq(seller1.ID),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionOfRelationTypeOptionalByNil() {
	product1 := ts.createProduct("", 0, 0.0, false, nil)
	product2 := ts.createProduct("", 0, 0.0, false, nil)

	seller2 := ts.createSeller("agustin", nil)

	match := ts.createSale(0, product1, nil)
	ts.createSale(0, product2, seller2)

	entities, err := orm.Query[models.Sale](
		ts.db,
		conditions.Sale.SellerID.Is().Null(),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Sale{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestConditionsOnUIntModel() {
	match := ts.createBrand("match")
	ts.createBrand("not_match")

	entities, err := orm.Query[models.Brand](
		ts.db,
		conditions.Brand.Name.Is().Eq("match"),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Brand{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestMultipleConditionsAreConnectedByAnd() {
	match := ts.createProduct("match", 3, 0, false, nil)
	ts.createProduct("not_match", 5, 0, false, nil)
	ts.createProduct("not_match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().GtOrEq(3),
		conditions.Product.Int.Is().LtOrEq(4),
		conditions.Product.String.Is().Eq("match"),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestNot() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 3, 0, false, nil)

	ts.createProduct("not_match", 2, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		orm.Not(conditions.Product.Int.Is().Eq(2)),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestNotWithMultipleConditionsAreConnectedByAnd() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 5, 0, false, nil)

	ts.createProduct("not_match", 2, 0, false, nil)
	ts.createProduct("not_match", 3, 0, false, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		orm.Not(
			conditions.Product.Int.Is().Gt(1),
			conditions.Product.Int.Is().Lt(4),
		),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestOr() {
	match1 := ts.createProduct("match", 2, 0, false, nil)
	match2 := ts.createProduct("match", 3, 0, false, nil)
	match3 := ts.createProduct("match_3", 3, 0, false, nil)

	ts.createProduct("not_match", 1, 0, false, nil)
	ts.createProduct("not_match", 4, 0, false, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		orm.Or(
			conditions.Product.Int.Is().Eq(2),
			conditions.Product.Int.Is().Eq(3),
			conditions.Product.String.Is().Eq("match_3"),
		),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestNotOr() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 5, 0, false, nil)
	match3 := ts.createProduct("match", 4, 0, false, nil)

	ts.createProduct("not_match", 2, 0, false, nil)
	ts.createProduct("not_match_string", 3, 0, false, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		orm.Not[models.Product](
			orm.Or(
				conditions.Product.Int.Is().Eq(2),
				conditions.Product.String.Is().Eq("not_match_string"),
			),
		),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestXor() {
	switch getDBDialector() {
	case cql.Postgres, cql.SQLite, cql.SQLServer:
		log.Println("Xor not compatible")
	case cql.MySQL:
		match1 := ts.createProduct("", 1, 0, false, nil)
		match2 := ts.createProduct("", 7, 0, false, nil)

		ts.createProduct("", 5, 0, false, nil)
		ts.createProduct("", 4, 0, false, nil)

		entities, err := orm.Query[models.Product](
			ts.db,
			mysql.Xor(
				conditions.Product.Int.Is().Lt(6),
				conditions.Product.Int.Is().Gt(3),
			),
		).Find()
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	}
}

func (ts *WhereConditionsIntTestSuite) TestMultipleConditionsDifferentOperators() {
	match1 := ts.createProduct("match", 1, 0.0, true, nil)
	match2 := ts.createProduct("match", 1, 0.0, true, nil)

	ts.createProduct("not_match", 1, 0.0, true, nil)
	ts.createProduct("match", 2, 0.0, true, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		conditions.Product.String.Is().Eq("match"),
		conditions.Product.Int.Is().Lt(2),
		conditions.Product.Bool.Is().True(),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestUnsafeCondition() {
	match1 := ts.createProduct("match", 1, 0.0, true, nil)
	match2 := ts.createProduct("match", 1, 0.0, true, nil)

	ts.createProduct("not_match", 2, 0.0, true, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		unsafe.NewCondition[models.Product]("%s.int = ?", 1),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestEmptyConnectionConditionMakesNothing() {
	match1 := ts.createProduct("match", 1, 0.0, true, nil)
	match2 := ts.createProduct("match", 1, 0.0, true, nil)

	entities, err := orm.Query[models.Product](
		ts.db,
		orm.And[models.Product](),
	).Find()
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *WhereConditionsIntTestSuite) TestEmptyContainerConditionReturnsError() {
	_, err := orm.Query[models.Product](
		ts.db,
		orm.Not[models.Product](),
	).Find()
	ts.ErrorIs(err, cql.ErrEmptyConditions)
}
