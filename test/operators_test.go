package test

import (
	"database/sql"
	"log"
	"strings"

	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/mysql"
	"github.com/FrancoLiberali/cql/psql"
	cqlSQL "github.com/FrancoLiberali/cql/sql"
	"github.com/FrancoLiberali/cql/sqlite"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

type OperatorsIntTestSuite struct {
	testSuite
}

func NewOperatorsIntTestSuite(
	db *gorm.DB,
) *OperatorsIntTestSuite {
	return &OperatorsIntTestSuite{
		testSuite: testSuite{
			db: db,
		},
	}
}

func (ts *OperatorsIntTestSuite) TestEqPointers() {
	intMatch := 1
	match := ts.createProduct("match", 1, 0, false, &intMatch)

	intNotMatch := 2
	ts.createProduct("match", 3, 0, false, &intNotMatch)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.IntPointer.Is().Eq(1),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestEqNullableType() {
	match := ts.createProduct("match", 0, 0, false, nil)
	match.NullFloat = sql.NullFloat64{Valid: true, Float64: 1.3}
	err := ts.db.Save(match).Error
	ts.Require().NoError(err)

	notMatch1 := ts.createProduct("not_match", 3, 0, false, nil)
	notMatch1.NullFloat = sql.NullFloat64{Valid: true, Float64: 1.2}
	err = ts.db.Save(notMatch1).Error
	ts.Require().NoError(err)

	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.NullFloat.Is().Eq(1.3),
	).Find()

	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestNotEq() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 3, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().NotEq(2),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestLt() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 2, 0, false, nil)
	ts.createProduct("not_match", 3, 0, false, nil)
	ts.createProduct("not_match", 4, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Lt(3),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestLtOrEq() {
	match1 := ts.createProduct("match", 1, 0, false, nil)
	match2 := ts.createProduct("match", 2, 0, false, nil)
	ts.createProduct("not_match", 3, 0, false, nil)
	ts.createProduct("not_match", 4, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().LtOrEq(2),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestGt() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Gt(2),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestGtOrEq() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().GtOrEq(3),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestBetween() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 6, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Between(3, 5),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestNotBetween() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 1, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().NotBetween(0, 2),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsNullPointers() {
	match := ts.createProduct("match", 0, 0, false, nil)
	int1 := 1
	int2 := 2

	ts.createProduct("not_match", 0, 0, false, &int1)
	ts.createProduct("not_match", 0, 0, false, &int2)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.IntPointer.Is().Null(),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsNullNullableTypes() {
	match := ts.createProduct("match", 0, 0, false, nil)

	notMatch := ts.createProduct("not_match", 0, 0, false, nil)
	notMatch.NullFloat = sql.NullFloat64{Valid: true, Float64: 6}
	err := ts.db.Save(notMatch).Error
	ts.Require().NoError(err)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.NullFloat.Is().Null(),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsNotNullPointers() {
	int1 := 1
	match := ts.createProduct("match", 0, 0, false, &int1)
	ts.createProduct("not_match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.IntPointer.Is().NotNull(),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsNotNullNullableTypes() {
	match := ts.createProduct("match", 0, 0, false, nil)
	match.NullFloat = sql.NullFloat64{Valid: true, Float64: 6}
	err := ts.db.Save(match).Error
	ts.Require().NoError(err)

	ts.createProduct("not_match", 0, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.NullFloat.Is().NotNull(),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsTrue() {
	match := ts.createProduct("match", 0, 0, true, nil)
	ts.createProduct("not_match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, false, nil)

	var err error

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Bool.Is().True(),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsFalse() {
	match := ts.createProduct("match", 0, 0, false, nil)
	ts.createProduct("not_match", 0, 0, true, nil)
	ts.createProduct("not_match", 0, 0, true, nil)

	var err error

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Bool.Is().False(),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsNotTrue() {
	match1 := ts.createProduct("match", 0, 0, false, nil)
	match2 := ts.createProduct("match", 0, 0, false, nil)
	match2.NullBool = sql.NullBool{Valid: true, Bool: false}
	err := ts.db.Save(match2).Error
	ts.Require().NoError(err)

	notMatch := ts.createProduct("not_match", 0, 0, false, nil)
	notMatch.NullBool = sql.NullBool{Valid: true, Bool: true}
	err = ts.db.Save(notMatch).Error
	ts.Require().NoError(err)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.NullBool.Is().NotTrue(),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsNotFalse() {
	match1 := ts.createProduct("match", 0, 0, false, nil)
	match2 := ts.createProduct("match", 0, 0, false, nil)
	match2.NullBool = sql.NullBool{Valid: true, Bool: true}
	err := ts.db.Save(match2).Error
	ts.Require().NoError(err)

	notMatch := ts.createProduct("not_match", 0, 0, false, nil)
	notMatch.NullBool = sql.NullBool{Valid: true, Bool: false}
	err = ts.db.Save(notMatch).Error
	ts.Require().NoError(err)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.NullBool.Is().NotFalse(),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsUnknown() {
	match := ts.createProduct("match", 0, 0, false, nil)

	notMatch1 := ts.createProduct("match", 0, 0, false, nil)
	notMatch1.NullBool = sql.NullBool{Valid: true, Bool: true}
	err := ts.db.Save(notMatch1).Error
	ts.Require().NoError(err)

	notMatch2 := ts.createProduct("not_match", 0, 0, false, nil)
	notMatch2.NullBool = sql.NullBool{Valid: true, Bool: false}
	err = ts.db.Save(notMatch2).Error
	ts.Require().NoError(err)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.NullBool.Is().Unknown(),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsNotUnknown() {
	match1 := ts.createProduct("", 0, 0, false, nil)
	match1.NullBool = sql.NullBool{Valid: true, Bool: true}
	err := ts.db.Save(match1).Error
	ts.Require().NoError(err)

	match2 := ts.createProduct("", 0, 0, false, nil)
	match2.NullBool = sql.NullBool{Valid: true, Bool: false}
	err = ts.db.Save(match2).Error
	ts.Require().NoError(err)

	ts.createProduct("", 0, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.NullBool.Is().NotUnknown(),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsDistinct() {
	match1 := ts.createProduct("match", 3, 0, false, nil)
	match2 := ts.createProduct("match", 4, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Distinct(2),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestIsNotDistinct() {
	match := ts.createProduct("match", 3, 0, false, nil)
	ts.createProduct("not_match", 4, 0, false, nil)
	ts.createProduct("not_match", 2, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().NotDistinct(3),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestArrayIn() {
	match1 := ts.createProduct("s1", 0, 0, false, nil)
	match2 := ts.createProduct("s2", 0, 0, false, nil)

	ts.createProduct("ns1", 0, 0, false, nil)
	ts.createProduct("ns2", 0, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.String.Is().In("s1", "s2", "s3"),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestArrayNotIn() {
	match1 := ts.createProduct("s1", 0, 0, false, nil)
	match2 := ts.createProduct("s2", 0, 0, false, nil)

	ts.createProduct("ns1", 0, 0, false, nil)
	ts.createProduct("ns2", 0, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.String.Is().NotIn("ns1", "ns2"),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestLike() {
	match1 := ts.createProduct("basd", 0, 0, false, nil)
	match2 := ts.createProduct("cape", 0, 0, false, nil)

	ts.createProduct("bbsd", 0, 0, false, nil)
	ts.createProduct("bbasd", 0, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.String.Is().Like("_a%"),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestLikeEscape() {
	match1 := ts.createProduct("ba_sd", 0, 0, false, nil)
	match2 := ts.createProduct("ca_pe", 0, 0, false, nil)

	ts.createProduct("bb_sd", 0, 0, false, nil)
	ts.createProduct("bba_sd", 0, 0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.String.Is().Custom(
			condition.Like("_a!_%").Escape('!'),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
}

func (ts *OperatorsIntTestSuite) TestLikeOnNumeric() {
	switch getDBDialector() {
	case cqlSQL.Postgres, cqlSQL.SQLServer, cqlSQL.SQLite:
		log.Println("Like with numeric not compatible")
	case cqlSQL.MySQL:
		match1 := ts.createProduct("", 10, 0, false, nil)
		match2 := ts.createProduct("", 100, 0, false, nil)

		ts.createProduct("", 20, 0, false, nil)
		ts.createProduct("", 3, 0, false, nil)

		entities, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Int.Is().Custom(
				mysql.Like[int]("1%"),
			),
		).Find()
		ts.Require().NoError(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	}
}

func (ts *OperatorsIntTestSuite) TestILike() {
	switch getDBDialector() {
	case cqlSQL.MySQL, cqlSQL.SQLServer, cqlSQL.SQLite:
		_, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.String.Is().Custom(
				psql.ILike("_a%"),
			),
		).Find()
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "operator: psql.ILike")
	case cqlSQL.Postgres:
		match1 := ts.createProduct("basd", 0, 0, false, nil)
		match2 := ts.createProduct("cape", 0, 0, false, nil)
		match3 := ts.createProduct("bAsd", 0, 0, false, nil)

		ts.createProduct("bbsd", 0, 0, false, nil)
		ts.createProduct("bbasd", 0, 0, false, nil)

		entities, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.String.Is().Custom(
				psql.ILike("_a%"),
			),
		).Find()
		ts.Require().NoError(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
	}
}

func (ts *OperatorsIntTestSuite) TestSimilarTo() {
	switch getDBDialector() {
	case cqlSQL.MySQL, cqlSQL.SQLServer, cqlSQL.SQLite:
		_, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.String.Is().Custom(
				psql.SimilarTo("%(b|d)%"),
			),
		).Find()
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "operator: psql.SimilarTo")
	case cqlSQL.Postgres:
		match1 := ts.createProduct("abc", 0, 0, false, nil)
		match2 := ts.createProduct("aabcc", 0, 0, false, nil)

		ts.createProduct("aec", 0, 0, false, nil)
		ts.createProduct("aaaaa", 0, 0, false, nil)

		entities, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.String.Is().Custom(
				psql.SimilarTo("%(b|d)%"),
			),
		).Find()
		ts.Require().NoError(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	}
}

func (ts *OperatorsIntTestSuite) TestPosixRegexCaseSensitive() {
	match1 := ts.createProduct("ab", 0, 0, false, nil)
	match2 := ts.createProduct("ax", 0, 0, false, nil)

	ts.createProduct("bb", 0, 0, false, nil)
	ts.createProduct("cx", 0, 0, false, nil)
	ts.createProduct("AB", 0, 0, false, nil)

	var posixRegexOperator condition.Operator[string]

	switch getDBDialector() {
	case cqlSQL.SQLServer, cqlSQL.MySQL:
		_, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.String.Is().Custom(
				psql.POSIXMatch("^a(b|x)"),
			),
		).Find()

		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "operator: psql.POSIXMatch")

		_, err = cql.Query[models.Product](
			ts.db,
			conditions.Product.String.Is().Custom(
				sqlite.Glob("a[bx]"),
			),
		).Find()

		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "operator: sqlite.Glob")
	case cqlSQL.Postgres:
		posixRegexOperator = psql.POSIXMatch("^a(b|x)")
	case cqlSQL.SQLite:
		posixRegexOperator = sqlite.Glob("a[bx]")
	}

	if posixRegexOperator != nil {
		entities, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.String.Is().Custom(
				posixRegexOperator,
			),
		).Find()
		ts.Require().NoError(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	}
}

func (ts *OperatorsIntTestSuite) TestPosixRegexCaseInsensitive() {
	match1 := ts.createProduct("ab", 0, 0, false, nil)
	match2 := ts.createProduct("ax", 0, 0, false, nil)
	match3 := ts.createProduct("AB", 0, 0, false, nil)

	ts.createProduct("bb", 0, 0, false, nil)
	ts.createProduct("cx", 0, 0, false, nil)

	var posixRegexOperator condition.Operator[string]

	switch getDBDialector() {
	case cqlSQL.SQLServer, cqlSQL.SQLite:
		_, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.String.Is().Custom(
				mysql.Regexp("^a(b|x)"),
			),
		).Find()

		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "operator: mysql.Regexp")

		_, err = cql.Query[models.Product](
			ts.db,
			conditions.Product.String.Is().Custom(
				psql.POSIXIMatch("^a(b|x)"),
			),
		).Find()

		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "operator: psql.POSIXIMatch")
	case cqlSQL.MySQL:
		posixRegexOperator = mysql.Regexp("^a(b|x)")
	case cqlSQL.Postgres:
		posixRegexOperator = psql.POSIXIMatch("^a(b|x)")
	}

	if posixRegexOperator != nil {
		entities, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.String.Is().Custom(
				posixRegexOperator,
			),
		).Find()
		ts.Require().NoError(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
	}
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericTypeWithSameType() {
	int1 := 1
	product1 := ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value()),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericTypeWithDifferentType() {
	int1 := 1
	product1 := ts.createProduct("", 1, 1.0, false, &int1)
	ts.createProduct("", 2, 2.1, false, &int1)
	ts.createProduct("", 0, 2.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.IsDynamic().Eq(conditions.Product.Float.Value()),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForCustomType() {
	match := ts.createProduct("salut,hola", 1, 0.0, false, nil)
	match.MultiString = models.MultiString{"salut", "hola"}
	err := ts.db.Save(match).Error
	ts.Require().NoError(err)

	ts.createProduct("salut,hola", 1, 0.0, false, nil)
	ts.createProduct("hola", 1, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.MultiString.IsDynamic().Eq(conditions.Product.MultiString.Value()),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForBaseModelAttribute() {
	match := ts.createProduct("", 1, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.CreatedAt.IsDynamic().Eq(conditions.Product.CreatedAt.Value()),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNotNullTypeCanBeComparedWithNullableType() {
	match := ts.createProduct("", 1, 1.0, false, nil)
	match.NullFloat = sql.NullFloat64{Valid: true, Float64: 1.0}
	err := ts.db.Save(match).Error
	ts.Require().NoError(err)

	ts.createProduct("", 1, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Float.IsDynamic().Eq(conditions.Product.NullFloat.Value()),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestUnsafeOperatorInCaseTypesNotMatchConvertible() {
	// comparisons between types are allowed when they are convertible
	match := ts.createProduct("", 0, 2.1, false, nil)
	ts.createProduct("", 0, 0, false, nil)
	ts.createProduct("", 0, 2, false, nil)
	ts.createProduct("", 0, 2.3, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Float.IsUnsafe().Eq("2.1"),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestUnsafeOperatorInCaseTypesNotMatchNotConvertible() {
	switch getDBDialector() {
	case cqlSQL.SQLite:
		// comparisons between types are allowed and matches nothing if not convertible
		ts.createProduct("", 0, 0, false, nil)
		ts.createProduct("", 0, 2, false, nil)
		ts.createProduct("", 0, 2.3, false, nil)

		entities, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Float.IsUnsafe().Eq("not_convertible_to_float"),
		).Find()
		ts.Require().NoError(err)

		EqualList(&ts.Suite, []*models.Product{}, entities)
	case cqlSQL.MySQL:
		// comparisons between types are allowed but matches 0s if not convertible
		match := ts.createProduct("", 0, 0, false, nil)
		ts.createProduct("", 0, 2, false, nil)
		ts.createProduct("", 0, 2.3, false, nil)

		entities, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Float.IsUnsafe().Eq("not_convertible_to_float"),
		).Find()
		ts.Require().NoError(err)

		EqualList(&ts.Suite, []*models.Product{match}, entities)
	case cqlSQL.SQLServer:
		// returns an error
		_, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Float.IsUnsafe().Eq("not_convertible_to_float"),
		).Find()
		ts.ErrorContains(err, "mssql: Error converting data type nvarchar to float.")
	case cqlSQL.Postgres:
		// returns an error
		_, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Float.IsUnsafe().Eq("not_convertible_to_float"),
		).Find()
		ts.ErrorContains(err, "not_convertible_to_float")
	}
}

func (ts *OperatorsIntTestSuite) TestUnsafeOperatorInCaseFieldWithTypesNotMatch() {
	switch getDBDialector() {
	case cqlSQL.SQLite:
		// comparisons between fields with different types are allowed
		match1 := ts.createProduct("0", 0, 0, false, nil)
		match2 := ts.createProduct("1", 0, 1, false, nil)
		ts.createProduct("0", 0, 1, false, nil)
		ts.createProduct("not_convertible", 0, 0, false, nil)

		entities, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Float.IsUnsafe().Eq(conditions.Product.String.Value()),
		).Find()
		ts.Require().NoError(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)
	case cqlSQL.MySQL:
		// comparisons between fields with different types are allowed but matches 0s on not convertible
		match1 := ts.createProduct("0", 1, 0, false, nil)
		match2 := ts.createProduct("1", 2, 1, false, nil)
		match3 := ts.createProduct("not_convertible", 2, 0, false, nil)
		ts.createProduct("0.0", 2, 1.0, false, nil)

		entities, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Float.IsUnsafe().Eq(conditions.Product.String.Value()),
		).Find()
		ts.Require().NoError(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2, match3}, entities)
	case cqlSQL.SQLServer:
		// comparisons between fields with different types are allowed and returns error only if at least one is not convertible
		match1 := ts.createProduct("0", 1, 0, false, nil)
		match2 := ts.createProduct("1", 2, 1, false, nil)

		entities, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Float.IsUnsafe().Eq(conditions.Product.String.Value()),
		).Find()
		ts.Require().NoError(err)

		EqualList(&ts.Suite, []*models.Product{match1, match2}, entities)

		ts.createProduct("not_convertible", 3, 0, false, nil)
		ts.createProduct("0.0", 4, 1.0, false, nil)

		_, err = cql.Query[models.Product](
			ts.db,
			conditions.Product.Float.IsUnsafe().Eq(conditions.Product.String.Value()),
		).Find()
		ts.ErrorContains(err, "mssql: Error converting data type nvarchar to float.")
	case cqlSQL.Postgres:
		// returns an error
		_, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Float.IsUnsafe().Eq(conditions.Product.String.Value()),
		).Find()

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
	ts.Require().NoError(err)

	notMatch := ts.createProduct("chau", 0, 0.0, false, nil)
	notMatch.MultiString = models.MultiString{"hola", "chau"}
	err = ts.db.Save(notMatch).Error
	ts.Require().NoError(err)

	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.String.IsUnsafe().Eq(conditions.Product.MultiString.Value()),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{match}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericWithPlus() {
	int1 := 1
	product1 := ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 3, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().Plus(1)),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericWithMinus() {
	int1 := 3
	product1 := ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 3, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().Minus(1)),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericWithTimes() {
	int1 := 1
	product1 := ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 3, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().Times(2)),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericWithDivided() {
	int1 := 4
	product1 := ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 3, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().Divided(2)),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericWithModulo() {
	int1 := 5
	product1 := ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().Modulo(2)),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericWithPower() {
	switch getDBDialector() {
	case cqlSQL.SQLite:
		_, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().Power(2)),
		).Find()
		ts.ErrorContains(err, "no such function: POWER")
	default:
		int1 := 2
		product1 := ts.createProduct("", 4, 0.0, false, &int1)
		ts.createProduct("", 2, 0.0, false, &int1)
		ts.createProduct("", 0, 0.0, false, nil)

		entities, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().Power(2)),
		).Find()
		ts.Require().NoError(err)

		EqualList(&ts.Suite, []*models.Product{product1}, entities)
	}
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericWithSquareRoot() {
	switch getDBDialector() {
	case cqlSQL.SQLite:
		_, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().SquareRoot()),
		).Find()
		ts.ErrorContains(err, "no such function: SQRT")
	default:
		int1 := 4
		product1 := ts.createProduct("", 2, 0.0, false, &int1)
		ts.createProduct("", 4, 0.0, false, &int1)
		ts.createProduct("", 0, 0.0, false, nil)

		entities, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().SquareRoot()),
		).Find()
		ts.Require().NoError(err)

		EqualList(&ts.Suite, []*models.Product{product1}, entities)
	}
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericWithAbsolute() {
	int1 := -2
	product1 := ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", -2, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().Absolute()),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericWithAnd() {
	int1 := 7
	product1 := ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 7, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().And(1)),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericWithOr() {
	int1 := 5
	product1 := ts.createProduct("", 7, 0.0, false, &int1)
	ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().Or(3)),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericWithXor() {
	switch getDBDialector() {
	case cqlSQL.SQLite:
		_, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().Xor(3)),
		).Find()
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "function: Xor")
	default:
		int1 := 5
		product1 := ts.createProduct("", 6, 0.0, false, &int1)
		ts.createProduct("", 3, 0.0, false, &int1)
		ts.createProduct("", 0, 0.0, false, nil)

		entities, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().Xor(3)),
		).Find()
		ts.Require().NoError(err)

		EqualList(&ts.Suite, []*models.Product{product1}, entities)
	}
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericWithNot() {
	int1 := 1
	product1 := ts.createProduct("", 4, 0.0, false, &int1)
	ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().Not().And(5)),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericWithShiftLeft() {
	int1 := 1
	product1 := ts.createProduct("", 4, 0.0, false, &int1)
	ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().ShiftLeft(2)),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericWithShiftRight() {
	int1 := 4
	product1 := ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 4, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().ShiftRight(2)),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericWithMultipleFunction() {
	int1 := 1
	product1 := ts.createProduct("", 4, 0.0, false, &int1)
	ts.createProduct("", 3, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().Plus(1).Times(2)),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForNumericWithFunctionOfDifferentType() {
	int1 := 2
	product1 := ts.createProduct("", 3, 0.0, false, &int1)
	product2 := ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.IsDynamic().Eq(conditions.Product.IntPointer.Value().Times(1.5)),
	).Find()
	ts.Require().NoError(err)

	switch getDBDialector() {
	case cqlSQL.Postgres:
		EqualList(&ts.Suite, []*models.Product{product2}, entities)
	case cqlSQL.MySQL, cqlSQL.SQLServer, cqlSQL.SQLite:
		EqualList(&ts.Suite, []*models.Product{product1}, entities)
	}
}

func (ts *OperatorsIntTestSuite) TestDynamicOperatorForStringWithConcat() {
	product1 := ts.createProduct("asd123", 2, 0.0, false, nil)
	product1.String2 = "asd"
	err := ts.db.Save(product1).Error
	ts.Require().NoError(err)

	product2 := ts.createProduct("asd", 3, 0.0, false, nil)
	product2.String2 = "asd"
	err = ts.db.Save(product2).Error
	ts.Require().NoError(err)

	ts.createProduct("asd123", 3, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.String.IsDynamic().Eq(conditions.Product.String2.Value().Concat("123")),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}
