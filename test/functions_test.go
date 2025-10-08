package test

import (
	"github.com/FrancoLiberali/cql"
	cqlSQL "github.com/FrancoLiberali/cql/sql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

type FunctionsIntTestSuite struct {
	testSuite
}

func NewFunctionsIntTestSuite(
	db *cql.DB,
) *FunctionsIntTestSuite {
	return &FunctionsIntTestSuite{
		testSuite: testSuite{
			db: db,
		},
	}
}

func (ts *FunctionsIntTestSuite) TestFunctionOnLeftSideWithStaticValue() {
	product1 := ts.createProduct("", 2, 0.0, false, nil)
	ts.createProduct("", 3, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Plus(cql.Int(1)).Is().Eq(cql.Int(3)),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestFunctionOnRightSideWithStaticValue() {
	int1 := 1
	product1 := ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 3, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.Plus(cql.Int(1))),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestFunctionOnBothSidesWithStaticValue() {
	int1 := 1
	product1 := ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 3, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Plus(cql.Int(1)).Is().Eq(conditions.Product.IntPointer.Plus(cql.Int(2))),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestFunctionOnLeftSideWithDynamicValue() {
	product1 := ts.createProduct("", 2, 0.0, false, nil)
	ts.createProduct("", 3, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Plus(conditions.Product.Int).Is().Eq(cql.Int(4)),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestFunctionOnRightSideWithDynamicValue() {
	int1 := 1
	product1 := ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 3, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.Plus(conditions.Product.IntPointer)),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestFunctionOnBothSidesWithDynamicValue() {
	int1 := 1
	product1 := ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 3, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Plus(conditions.Product.Int).Is().Eq(
			conditions.Product.IntPointer.Plus(conditions.Product.IntPointer),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestFunctionOnBothSidesWithDynamicValueWithFunction() {
	int1 := 1
	product1 := ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 3, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Plus(conditions.Product.Int.Minus(cql.Int(1))).Is().Eq(
			conditions.Product.IntPointer.Plus(conditions.Product.IntPointer.Minus(cql.Int(1))),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestFunctionOnBothSidesWithDynamicValueWithFunctions() {
	int1 := 1
	product1 := ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 3, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Plus(conditions.Product.Int.Minus(cql.Int(1))).Minus(cql.Int(1)).Minus(cql.Int(2)).Is().Eq(
			conditions.Product.IntPointer.Plus(conditions.Product.IntPointer.Minus(cql.Int(1))).Minus(cql.Int(3)),
		),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForNumericWithMinus() {
	int1 := 3
	product1 := ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 3, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.Minus(cql.Int(1))),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForNumericWithTimes() {
	int1 := 1
	product1 := ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 3, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.Times(cql.Int(2))),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForNumericWithDivided() {
	int1 := 4
	product1 := ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 3, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.Divided(cql.Int(2))),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForNumericWithModulo() {
	int1 := 5
	product1 := ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.Modulo(cql.Int(2))),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForNumericWithPower() {
	switch getDBDialector() {
	case cqlSQL.SQLite:
		_, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.Power(cql.Int(2))),
		).Find()
		ts.ErrorContains(err, "no such function: POWER")
	default:
		int1 := 2
		product1 := ts.createProduct("", 4, 0.0, false, &int1)
		ts.createProduct("", 2, 0.0, false, &int1)
		ts.createProduct("", 0, 0.0, false, nil)

		entities, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.Power(cql.Int(2))),
		).Find()
		ts.Require().NoError(err)

		EqualList(&ts.Suite, []*models.Product{product1}, entities)
	}
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForNumericWithSquareRoot() {
	switch getDBDialector() {
	case cqlSQL.SQLite:
		_, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.SquareRoot()),
		).Find()
		ts.ErrorContains(err, "no such function: SQRT")
	default:
		int1 := 4
		product1 := ts.createProduct("", 2, 0.0, false, &int1)
		ts.createProduct("", 4, 0.0, false, &int1)
		ts.createProduct("", 0, 0.0, false, nil)

		entities, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.SquareRoot()),
		).Find()

		if getDBDialector() == cqlSQL.Postgres && err != nil {
			// cockroach
			ts.ErrorContains(err, "unknown signature: sqrt(int) (SQLSTATE 42883)")
		} else {
			ts.Require().NoError(err)
			EqualList(&ts.Suite, []*models.Product{product1}, entities)
		}
	}
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForNumericWithAbsolute() {
	int1 := -2
	product1 := ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", -2, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.Absolute()),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForNumericWithAnd() {
	int1 := 7
	product1 := ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 7, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.And(cql.Int(1))),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForNumericWithAndDynamic() {
	int1 := 7
	product1 := ts.createProduct("", 1, 0.0, false, &int1)
	product2 := ts.createProduct("", 7, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.And(conditions.Product.Int)),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1, product2}, entities)
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForNumericWithOr() {
	int1 := 5
	product1 := ts.createProduct("", 7, 0.0, false, &int1)
	ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.Or(cql.Int(3))),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForNumericWithXor() {
	switch getDBDialector() {
	case cqlSQL.SQLite:
		_, err := cql.Query[models.Product](
			ts.db,
			conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.Xor(cql.Int(3))),
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
			conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.Xor(cql.Int(3))),
		).Find()
		ts.Require().NoError(err)

		EqualList(&ts.Suite, []*models.Product{product1}, entities)
	}
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForNumericWithNot() {
	int1 := 1
	product1 := ts.createProduct("", 4, 0.0, false, &int1)
	ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.Not().And(cql.Int(5))),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForNumericWithShiftLeft() {
	int1 := 1
	product1 := ts.createProduct("", 4, 0.0, false, &int1)
	ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.ShiftLeft(cql.Int(2))),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForNumericWithShiftRight() {
	int1 := 4
	product1 := ts.createProduct("", 1, 0.0, false, &int1)
	ts.createProduct("", 4, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.ShiftRight(cql.Int(2))),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForNumericWithMultipleFunction() {
	int1 := 1
	product1 := ts.createProduct("", 4, 0.0, false, &int1)
	ts.createProduct("", 3, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.Plus(cql.Int(1)).Times(cql.Int(2))),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForNumericWithFunctionOfDifferentType() {
	int1 := 2
	product1 := ts.createProduct("", 3, 0.0, false, &int1)
	product2 := ts.createProduct("", 2, 0.0, false, &int1)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.Int.Is().Eq(conditions.Product.IntPointer.Times(cql.Float64(1.5))),
	).Find()
	ts.Require().NoError(err)

	switch getDBDialector() {
	case cqlSQL.Postgres:
		EqualList(&ts.Suite, []*models.Product{product2}, entities)
	case cqlSQL.MySQL, cqlSQL.SQLServer, cqlSQL.SQLite:
		EqualList(&ts.Suite, []*models.Product{product1}, entities)
	}
}

func (ts *FunctionsIntTestSuite) TestDynamicOperatorForStringWithConcat() {
	product1 := ts.createProduct("asd123", 2, 0.0, false, nil)
	product1.String2 = "asd"
	err := ts.db.GormDB.Save(product1).Error
	ts.Require().NoError(err)

	product2 := ts.createProduct("asd", 3, 0.0, false, nil)
	product2.String2 = "asd"
	err = ts.db.GormDB.Save(product2).Error
	ts.Require().NoError(err)

	ts.createProduct("asd123", 3, 0.0, false, nil)
	ts.createProduct("", 0, 0.0, false, nil)

	entities, err := cql.Query[models.Product](
		ts.db,
		conditions.Product.String.Is().Eq(conditions.Product.String2.Concat(cql.String("123"))),
	).Find()
	ts.Require().NoError(err)

	EqualList(&ts.Suite, []*models.Product{product1}, entities)
}
