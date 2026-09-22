package test

import (
	"context"
	"database/sql/driver"
	"fmt"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/test/models"
)

// Cents is a minimal fixed-point "decimal" money type: exact, backed by an
// integer number of cents, stored in / read from an INTEGER column. It plays
// the role github.com/shopspring/decimal.Decimal would in a real app — the
// point of the prototype is that DecimalField never routes it through float64.
type Cents int64

// Cents satisfies condition.Decimal by being a driver.Valuer — no Go-side
// arithmetic methods are needed; CQL evaluates the math in SQL.
func (c Cents) Value() (driver.Value, error) { return int64(c), nil }

func (c *Cents) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		*c = Cents(v)
	case int:
		*c = Cents(v)
	default:
		return fmt.Errorf("cannot scan %T into Cents", src)
	}

	return nil
}

type centsResult struct{ Total Cents }

// Amount views the existing integer "Int" column of Product as a decimal money
// column. In a real app cql-gen would emit this as conditions.Product.Amount.
func amountField() condition.DecimalField[models.Product, Cents] {
	return condition.NewDecimalField[models.Product, Cents]("Int", "", "")
}

// SUM over a DecimalField stays exact past 2^53 — the value float64 cannot
// represent — with no arbitrary destination type (compile-time safe, unlike
// the rejected IntoAs).
func (ts *SelectIntTestSuite) TestDecimalFieldSumExact() {
	// 2^53 + 1: float64 would round this to ...992.
	ts.createProduct("a", 9007199254740993, 0, false, nil)
	ts.createProduct("b", 0, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](context.Background(), ts.db),
		amountField().Aggregate().Sum().Into(func(r *centsResult) *Cents { return &r.Total }),
	)

	ts.Require().NoError(err)
	ts.Require().Len(results, 1)
	ts.Equal(Cents(9007199254740993), results[0].Total)
}

// Closed decimal arithmetic: decimal + decimal stays decimal, bound exactly.
func (ts *SelectIntTestSuite) TestDecimalFieldArithmetic() {
	ts.createProduct("a", 10, 0, false, nil)
	ts.createProduct("b", 20, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](context.Background(), ts.db).Ascending(amountField()),
		amountField().Plus(condition.Value[Cents]{Value: Cents(5)}).Into(func(r *centsResult) *Cents { return &r.Total }),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []centsResult{{Total: 15}, {Total: 25}}, results)
}

// Plain read of a decimal column binds exactly into the decimal type.
func (ts *SelectIntTestSuite) TestDecimalFieldExactRead() {
	ts.createProduct("a", 4200, 0, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](context.Background(), ts.db),
		amountField().Into(func(r *centsResult) *Cents { return &r.Total }),
	)

	ts.Require().NoError(err)
	ts.Require().Len(results, 1)
	ts.Equal(Cents(4200), results[0].Total)
}
