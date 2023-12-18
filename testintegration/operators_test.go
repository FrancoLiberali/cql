package testintegration

import (
	"database/sql"
	"strings"

	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/dynamic"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/unsafe"
	"github.com/ditrit/badaas/testintegration/conditions"
	"github.com/ditrit/badaas/testintegration/models"
)

type OperatorsIntTestSuite struct {
	CRUDServiceCommonIntTestSuite
	crudProductService orm.CRUDService[models.Product, model.UUID]
}

func NewOperatorsIntTestSuite(
	db *gorm.DB,
	crudProductService orm.CRUDService[models.Product, model.UUID],
) *OperatorsIntTestSuite {
	return &OperatorsIntTestSuite{
		CRUDServiceCommonIntTestSuite: CRUDServiceCommonIntTestSuite{
			db: db,
		},
		crudProductService: crudProductService,
	}
}

func (ts *OperatorsIntTestSuite) TestEqPointers() {
	intMatch := 1
	match := ts.createProduct("match", 1, 0, false, &intMatch)

	intNotMatch := 2
	ts.createProduct("match", 3, 0, false, &intNotMatch)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductIntPointer(
			orm.Eq(1),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestEqNullableType() {
	match := ts.createProduct("match", 0, 0, false, nil)
	match.NullFloat = sql.NullFloat64{Valid: true, Float64: 1.3}
	err := ts.db.Save(match).Error
	ts.Nil(err)

	notMatch1 := ts.createProduct("not_match", 3, 0, false, nil)
	notMatch1.NullFloat = sql.NullFloat64{Valid: true, Float64: 1.2}
	err = ts.db.Save(notMatch1).Error
	ts.Nil(err)

	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullFloat(
			orm.Eq(1.3),
		),
	)

	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestNotEq() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 3, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.NotEq(2),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestLt() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 2, 0, false, nil)
	ts.createProduct("not_match", 3, 0, false, nil)
	ts.createProduct("not_match", 4, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.Lt(3),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestLtOrEq() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 2, 0, false, nil)
	ts.createProduct("not_match", 3, 0, false, nil)
	ts.createProduct("not_match", 4, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.LtOrEq(2),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestGt() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.Gt(2),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestGtOrEq() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.GtOrEq(3),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestBetween() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 6, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.Between(3, 5),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestNotBetween() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.NotBetween(0, 2),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsNullPointers() {
	match := ts.createProduct("match", 0, 0, false, nil)
	int1 := 1
	int2 := 2

	ts.createProduct("not_match", 0, 0, false, &int1)
	ts.createProduct("not_match", 0, 0, false, &int2)

	entities, err := ts.crudProductService.Query(
		conditions.ProductIntPointer(
			orm.IsNull[int](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsNullNullableTypes() {
	match := ts.createProduct("match", 0, 0, false, nil)

	notMatch := ts.createProduct("not_match", 0, 0, false, nil)
	notMatch.NullFloat = sql.NullFloat64{Valid: true, Float64: 6}
	err := ts.db.Save(notMatch).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullFloat(
			orm.IsNull[float64](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsNotNullPointers() {
	int1 := 1
	match := ts.createProduct("match", 0, 0, false, &int1)
	ts.createProduct("not_match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductIntPointer(
			orm.IsNotNull[int](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsNotNullNullableTypes() {
	match := ts.createProduct("match", 0, 0, false, nil)
	match.NullFloat = sql.NullFloat64{Valid: true, Float64: 6}
	err := ts.db.Save(match).Error
	ts.Nil(err)

	ts.createProduct("not_match", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullFloat(
			orm.IsNotNull[float64](),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsTrue() {
	match := ts.createProduct("match", 0, 0, true, nil)
	ts.createProduct("not_match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductBool(
			orm.IsTrue(),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsFalse() {
	match := ts.createProduct("match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, true, nil)
	ts.createProduct("not_match", 0, 0, true, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductBool(
			orm.IsFalse(),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsNotTrue() {
	match1 := ts.createProduct("match", 0, 0, false, nil)
	match2 := ts.createProduct("match", 0, 0, false, nil)
	match2.NullBool = sql.NullBool{Valid: true, Bool: false}
	err := ts.db.Save(match2).Error
	ts.Nil(err)

	notMatch := ts.createProduct("not_match", 0, 0, false, nil)
	notMatch.NullBool = sql.NullBool{Valid: true, Bool: true}
	err = ts.db.Save(notMatch).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullBool(
			orm.IsNotTrue(),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsNotFalse() {
	match1 := ts.createProduct("match", 0, 0, false, nil)
	match2 := ts.createProduct("match", 0, 0, false, nil)
	match2.NullBool = sql.NullBool{Valid: true, Bool: true}
	err := ts.db.Save(match2).Error
	ts.Nil(err)

	notMatch := ts.createProduct("not_match", 0, 0, false, nil)
	notMatch.NullBool = sql.NullBool{Valid: true, Bool: false}
	err = ts.db.Save(notMatch).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullBool(
			orm.IsNotFalse(),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsUnknown() {
	match := ts.createProduct("match", 0, 0, false, nil)

	notMatch1 := ts.createProduct("match", 0, 0, false, nil)
	notMatch1.NullBool = sql.NullBool{Valid: true, Bool: true}
	err := ts.db.Save(notMatch1).Error
	ts.Nil(err)

	notMatch2 := ts.createProduct("not_match", 0, 0, false, nil)
	notMatch2.NullBool = sql.NullBool{Valid: true, Bool: false}
	err = ts.db.Save(notMatch2).Error
	ts.Nil(err)

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullBool(
			orm.IsUnknown(),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsNotUnknown() {
	match1 := ts.createProduct("", 0, 0, false, nil)
	match1.NullBool = sql.NullBool{Valid: true, Bool: true}
	err := ts.db.Save(match1).Error
	ts.Nil(err)

	match2 := ts.createProduct("", 0, 0, false, nil)
	match2.NullBool = sql.NullBool{Valid: true, Bool: false}
	err = ts.db.Save(match2).Error
	ts.Nil(err)

	ts.createProduct("", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullBool(
			orm.IsNotUnknown(),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsDistinct() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.IsDistinct(2),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsNotDistinct() {
	match := ts.createProduct("match", 3, 0, false, nil)
	ts.createProduct("not_match", 4, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			orm.IsNotDistinct(3),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestArrayIn() {
	match1 := ts.createProduct("s1", 0, 0, false, nil)
	match2 := ts.createProduct("s2", 0, 0, false, nil)

	ts.createProduct("ns1", 0, 0, false, nil)
	ts.createProduct("ns2", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductString(
			orm.ArrayIn("s1", "s2", "s3"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestArrayNotIn() {
	match1 := ts.createProduct("s1", 0, 0, false, nil)
	match2 := ts.createProduct("s2", 0, 0, false, nil)

	ts.createProduct("ns1", 0, 0, false, nil)
	ts.createProduct("ns2", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductString(
			orm.ArrayNotIn("ns1", "ns2"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestLike() {
	match1 := ts.createProduct("basd", 0, 0, false, nil)
	match2 := ts.createProduct("cape", 0, 0, false, nil)

	ts.createProduct("bbsd", 0, 0, false, nil)
	ts.createProduct("bbasd", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductString(
			orm.Like("_a%"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestLikeEscape() {
	match1 := ts.createProduct("ba_sd", 0, 0, false, nil)
	match2 := ts.createProduct("ca_pe", 0, 0, false, nil)

	ts.createProduct("bb_sd", 0, 0, false, nil)
	ts.createProduct("bba_sd", 0, 0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductString(
			orm.Like("_a!_%").Escape('!'),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForBasicType() {
	int1 := 1
	product1 := ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductInt(
			dynamic.Eq(conditions.ProductIntPointerField),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForCustomType() {
	match := ts.createProduct("salut,hola", 1, 0.0, false, nil)
	match.MultiString = models.MultiString{"salut", "hola"}
	err := ts.db.Save(match).Error
	ts.Nil(err)

	ts.createProduct("salut,hola", 1, 0.0, false, nil)
	ts.createProduct("hola", 1, 0.0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductMultiString(
			dynamic.Eq(conditions.ProductMultiStringField),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForBaseModelAttribute() {
	match := ts.createProduct("", 1, 0.0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductDeletedAt(dynamic.IsNotDistinct(conditions.ProductDeletedAtField)),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNotNullTypeCanBeComparedWithNullableType() {
	match := ts.createProduct("", 1, 1.0, false, nil)
	match.NullFloat = sql.NullFloat64{Valid: true, Float64: 1.0}
	err := ts.db.Save(match).Error
	ts.Nil(err)

	ts.createProduct("", 1, 0.0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductFloat(
			dynamic.Eq[float64](conditions.ProductNullFloatField),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestUnsafeOperatorInCaseTypesNotMatchConvertible() {
	// comparisons between types are allowed when they are convertible
	match := ts.createProduct("", 0, 2.1, false, nil)
	ts.createProduct("", 0, 0, false, nil)
	ts.createProduct("", 0, 2, false, nil)
	ts.createProduct("", 0, 2.3, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductFloat(
			unsafe.Eq[float64]("2.1"),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestUnsafeOperatorInCaseTypesNotMatchNotConvertible() {
	_, err := ts.crudProductService.Query(
		conditions.ProductFloat(
			unsafe.Eq[float64]("not_convertible_to_float"),
		),
	)
	ts.ErrorContains(err, "not_convertible_to_float")
}

func (ts *OperatorsIntTestSuite) TestUnsafeOperatorInCaseFieldWithTypesNotMatch() {
	_, err := ts.crudProductService.Query(
		conditions.ProductFloat(
			unsafe.Eq[float64](conditions.ProductStringField),
		),
	)

	ts.True(
		strings.Contains(
			err.Error(),
			"ERROR: operator does not exist: numeric = text (SQLSTATE 42883)", // postgresql
		) || strings.Contains(
			err.Error(),
			"ERROR: unsupported comparison operator: <decimal> = <string> (SQLSTATE 22023)", // cockroachdb
		),
	)
}

func (ts *OperatorsIntTestSuite) TestUnsafeOperatorCanCompareFieldsThatMapToTheSameType() {
	match := ts.createProduct("hola,chau", 1, 1.0, false, nil)
	match.MultiString = models.MultiString{"hola", "chau"}
	err := ts.db.Save(match).Error
	ts.Nil(err)

	notMatch := ts.createProduct("chau", 0, 0.0, false, nil)
	notMatch.MultiString = models.MultiString{"hola", "chau"}
	err = ts.db.Save(notMatch).Error
	ts.Nil(err)

	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := ts.crudProductService.Query(
		conditions.ProductString(
			unsafe.Eq[string](conditions.ProductMultiStringField),
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}
