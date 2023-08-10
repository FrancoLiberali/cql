package testintegration

import (
	"database/sql"
	"log"
	"strings"

	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/dynamic"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/orm/mysql"
	"github.com/ditrit/badaas/orm/operator"
	"github.com/ditrit/badaas/orm/psql"
	"github.com/ditrit/badaas/orm/sqlite"
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

	var isTrueOperator operator.Operator[bool]

	switch getDBDialector() {
	case postgreSQL, mySQL, sqLite:
		isTrueOperator = orm.IsTrue()
	case sqlServer:
		isTrueOperator = orm.Eq(true)
	}

	entities, err := ts.crudProductService.Query(
		conditions.ProductBool(
			isTrueOperator,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsFalse() {
	match := ts.createProduct("match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, true, nil)
	ts.createProduct("not_match", 0, 0, true, nil)

	var isFalseOperator operator.Operator[bool]

	switch getDBDialector() {
	case postgreSQL, mySQL, sqLite:
		isFalseOperator = orm.IsFalse()
	case sqlServer:
		isFalseOperator = orm.Eq(false)
	}

	entities, err := ts.crudProductService.Query(
		conditions.ProductBool(
			isFalseOperator,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

//nolint:dupl // not really duplicated
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

	var isNotTrueOperator operator.Operator[bool]

	switch getDBDialector() {
	case postgreSQL, mySQL, sqLite:
		isNotTrueOperator = orm.IsNotTrue()
	case sqlServer:
		isNotTrueOperator = orm.IsDistinct(true)
	}

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullBool(
			isNotTrueOperator,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

//nolint:dupl // not duplicated
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

	var isNotFalseOperator operator.Operator[bool]

	switch getDBDialector() {
	case postgreSQL, mySQL, sqLite:
		isNotFalseOperator = orm.IsNotFalse()
	case sqlServer:
		isNotFalseOperator = orm.IsDistinct(false)
	}

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullBool(
			isNotFalseOperator,
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

	var isUnknownOperator operator.Operator[bool]

	switch getDBDialector() {
	case postgreSQL, mySQL:
		isUnknownOperator = orm.IsUnknown()
	case sqlServer, sqLite:
		isUnknownOperator = orm.IsNull[bool]()
	}

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullBool(
			isUnknownOperator,
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

	var isNotUnknownOperator operator.Operator[bool]

	switch getDBDialector() {
	case postgreSQL, mySQL:
		isNotUnknownOperator = orm.IsNotUnknown()
	case sqlServer, sqLite:
		isNotUnknownOperator = orm.IsNotNull[bool]()
	}

	entities, err := ts.crudProductService.Query(
		conditions.ProductNullBool(
			isNotUnknownOperator,
		),
	)
	ts.Nil(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsDistinct() {
	switch getDBDialector() {
	case postgreSQL, sqlServer, sqLite:
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
	case mySQL:
		log.Println("IsDistinct not compatible")
	}
}

func (ts *OperatorsIntTestSuite) TestIsNotDistinct() {
	switch getDBDialector() {
	case postgreSQL, sqlServer, sqLite:
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
	case mySQL:
		log.Println("IsNotDistinct not compatible")
	}
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

func (ts *OperatorsIntTestSuite) TestLikeOnNumeric() {
	switch getDBDialector() {
	case postgreSQL, sqlServer, sqLite:
		log.Println("Like with numeric not compatible")
	case mySQL:
		match1 := ts.createProduct("", 10, 0, false, nil)
		match2 := ts.createProduct("", 100, 0, false, nil)

		ts.createProduct("", 20, 0, false, nil)
		ts.createProduct("", 3, 0, false, nil)

		entities, err := ts.crudProductService.Query(
			conditions.ProductInt(
				mysql.Like[int]("1%"),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	}
}

func (ts *OperatorsIntTestSuite) TestILike() {
	switch getDBDialector() {
	case mySQL, sqlServer, sqLite:
		log.Println("ILike not compatible")
	case postgreSQL:
		match1 := ts.createProduct("basd", 0, 0, false, nil)
		match2 := ts.createProduct("cape", 0, 0, false, nil)
		match3 := ts.createProduct("bAsd", 0, 0, false, nil)

		ts.createProduct("bbsd", 0, 0, false, nil)
		ts.createProduct("bbasd", 0, 0, false, nil)

		entities, err := ts.crudProductService.Query(
			conditions.ProductString(
				psql.ILike("_a%"),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
	}
}

func (ts *OperatorsIntTestSuite) TestSimilarTo() {
	switch getDBDialector() {
	case mySQL, sqlServer, sqLite:
		log.Println("SimilarTo not compatible")
	case postgreSQL:
		match1 := ts.createProduct("abc", 0, 0, false, nil)
		match2 := ts.createProduct("aabcc", 0, 0, false, nil)

		ts.createProduct("aec", 0, 0, false, nil)
		ts.createProduct("aaaaa", 0, 0, false, nil)

		entities, err := ts.crudProductService.Query(
			conditions.ProductString(
				psql.SimilarTo("%(b|d)%"),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	}
}

func (ts *OperatorsIntTestSuite) TestPosixRegexCaseSensitive() {
	match1 := ts.createProduct("ab", 0, 0, false, nil)
	match2 := ts.createProduct("ax", 0, 0, false, nil)

	ts.createProduct("bb", 0, 0, false, nil)
	ts.createProduct("cx", 0, 0, false, nil)
	ts.createProduct("AB", 0, 0, false, nil)

	var posixRegexOperator operator.Operator[string]

	switch getDBDialector() {
	case sqlServer, mySQL:
		log.Println("PosixRegex not compatible")
	case postgreSQL:
		posixRegexOperator = psql.POSIXMatch("^a(b|x)")
	case sqLite:
		posixRegexOperator = sqlite.Glob("a[bx]")
	}

	if posixRegexOperator != nil {
		entities, err := ts.crudProductService.Query(
			conditions.ProductString(
				posixRegexOperator,
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	}
}

func (ts *OperatorsIntTestSuite) TestPosixRegexCaseInsensitive() {
	match1 := ts.createProduct("ab", 0, 0, false, nil)
	match2 := ts.createProduct("ax", 0, 0, false, nil)
	match3 := ts.createProduct("AB", 0, 0, false, nil)

	ts.createProduct("bb", 0, 0, false, nil)
	ts.createProduct("cx", 0, 0, false, nil)

	var posixRegexOperator operator.Operator[string]

	switch getDBDialector() {
	case sqlServer, sqLite:
		log.Println("PosixRegex Case Insensitive not compatible")
	case mySQL:
		posixRegexOperator = mysql.RegexP("^a(b|x)")
	case postgreSQL:
		posixRegexOperator = psql.POSIXIMatch("^a(b|x)")
	}

	if posixRegexOperator != nil {
		entities, err := ts.crudProductService.Query(
			conditions.ProductString(
				posixRegexOperator,
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
	}
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
		conditions.ProductCreatedAt(dynamic.Eq(conditions.ProductCreatedAtField)),
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
	switch getDBDialector() {
	case sqLite:
		// comparisons between types are allowed and matches nothing if not convertible
		ts.createProduct("", 0, 0, false, nil)
		ts.createProduct("", 0, 2, false, nil)
		ts.createProduct("", 0, 2.3, false, nil)

		entities, err := ts.crudProductService.Query(
			conditions.ProductFloat(
				unsafe.Eq[float64]("not_convertible_to_float"),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{}, entities)
	case mySQL:
		// comparisons between types are allowed but matches 0s if not convertible
		match := ts.createProduct("", 0, 0, false, nil)
		ts.createProduct("", 0, 2, false, nil)
		ts.createProduct("", 0, 2.3, false, nil)

		entities, err := ts.crudProductService.Query(
			conditions.ProductFloat(
				unsafe.Eq[float64]("not_convertible_to_float"),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match}, entities)
	case sqlServer:
		// returns an error
		_, err := ts.crudProductService.Query(
			conditions.ProductFloat(
				unsafe.Eq[float64]("not_convertible_to_float"),
			),
		)
		ts.ErrorContains(err, "mssql: Error converting data type nvarchar to float.")
	case postgreSQL:
		// returns an error
		_, err := ts.crudProductService.Query(
			conditions.ProductFloat(
				unsafe.Eq[float64]("not_convertible_to_float"),
			),
		)
		ts.ErrorContains(err, "not_convertible_to_float")
	}
}

func (ts *OperatorsIntTestSuite) TestUnsafeOperatorInCaseFieldWithTypesNotMatch() {
	switch getDBDialector() {
	case sqLite:
		// comparisons between fields with different types are allowed
		match1 := ts.createProduct("0", 0, 0, false, nil)
		match2 := ts.createProduct("1", 0, 1, false, nil)
		ts.createProduct("0", 0, 1, false, nil)
		ts.createProduct("not_convertible", 0, 0, false, nil)

		entities, err := ts.crudProductService.Query(
			conditions.ProductFloat(
				unsafe.Eq[float64](conditions.ProductStringField),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	case mySQL:
		// comparisons between fields with different types are allowed but matches 0s on not convertible
		match1 := ts.createProduct("0", 1, 0, false, nil)
		match2 := ts.createProduct("1", 2, 1, false, nil)
		match3 := ts.createProduct("not_convertible", 2, 0, false, nil)
		ts.createProduct("0.0", 2, 1.0, false, nil)

		entities, err := ts.crudProductService.Query(
			conditions.ProductFloat(
				unsafe.Eq[float64](conditions.ProductStringField),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
	case sqlServer:
		// comparisons between fields with different types are allowed and returns error only if at least one is not convertible
		match1 := ts.createProduct("0", 1, 0, false, nil)
		match2 := ts.createProduct("1", 2, 1, false, nil)

		entities, err := ts.crudProductService.Query(
			conditions.ProductFloat(
				unsafe.Eq[float64](conditions.ProductStringField),
			),
		)
		ts.Nil(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)

		ts.createProduct("not_convertible", 3, 0, false, nil)
		ts.createProduct("0.0", 4, 1.0, false, nil)

		_, err = ts.crudProductService.Query(
			conditions.ProductFloat(
				unsafe.Eq[float64](conditions.ProductStringField),
			),
		)
		ts.ErrorContains(err, "mssql: Error converting data type nvarchar to float.")
	case postgreSQL:
		// returns an error
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
