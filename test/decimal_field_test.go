package test

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

// The Bank model has two DECIMAL columns, exercising DecimalField end-to-end
// (generated conditions + scanner) with two backing types:
//   - Balance: models.Cents, a manually-defined decimal type (no library)
//   - Rate:    shopspring/decimal.Decimal, a real decimal library
//
// NOTE: sqlite has no true DECIMAL type (NUMERIC affinity stores as int/real),
// so the fractional Rate tests use values that are exact in float64 (halves)
// to stay reliable across every dialect. The integer-valued Cents column is
// exact everywhere.

func (ts *SelectIntTestSuite) createBank(balance models.Cents, rate decimal.Decimal) {
	ts.Require().NoError(ts.db.GormDB.Create(&models.Bank{Balance: balance, Rate: rate}).Error)
}

// --- manually-defined decimal type: models.Cents ---

type centsResult struct{ Total models.Cents }

func (ts *SelectIntTestSuite) TestDecimalFieldManualExactRead() {
	ts.createBank(4200, decimal.Zero)

	results, err := cql.Select(
		cql.Query[models.Bank](context.Background(), ts.db),
		conditions.Bank.Balance.Into(func(r *centsResult) *models.Cents { return &r.Total }),
	)

	ts.Require().NoError(err)
	ts.Require().Len(results, 1)
	ts.Equal(models.Cents(4200), results[0].Total)
}

// Sum stays exact past 2^53 — the value float64 cannot represent — because the
// column is NUMERIC and DecimalField never routes it through float64.
func (ts *SelectIntTestSuite) TestDecimalFieldManualSumExact() {
	ts.createBank(9007199254740993, decimal.Zero) // 2^53 + 1
	ts.createBank(0, decimal.Zero)

	results, err := cql.Select(
		cql.Query[models.Bank](context.Background(), ts.db),
		conditions.Bank.Balance.Aggregate().Sum().Into(func(r *centsResult) *models.Cents { return &r.Total }),
	)

	ts.Require().NoError(err)
	ts.Require().Len(results, 1)
	ts.Equal(models.Cents(9007199254740993), results[0].Total)
}

// Closed arithmetic: decimal + decimal stays decimal, bound exactly.
func (ts *SelectIntTestSuite) TestDecimalFieldManualArithmetic() {
	ts.createBank(10, decimal.Zero)
	ts.createBank(20, decimal.Zero)

	results, err := cql.Select(
		cql.Query[models.Bank](context.Background(), ts.db).Ascending(conditions.Bank.Balance),
		conditions.Bank.Balance.Plus(condition.Value[models.Cents]{Value: 5}).Into(
			func(r *centsResult) *models.Cents { return &r.Total },
		),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []centsResult{{Total: 15}, {Total: 25}}, results)
}

// --- decimal library: shopspring/decimal ---

type rateResult struct{ Total decimal.Decimal }

func (ts *SelectIntTestSuite) TestDecimalFieldLibraryExactRead() {
	ts.createBank(0, decimal.RequireFromString("12.5"))

	results, err := cql.Select(
		cql.Query[models.Bank](context.Background(), ts.db),
		conditions.Bank.Rate.Into(func(r *rateResult) *decimal.Decimal { return &r.Total }),
	)

	ts.Require().NoError(err)
	ts.Require().Len(results, 1)
	ts.True(results[0].Total.Equal(decimal.RequireFromString("12.5")))
}

func (ts *SelectIntTestSuite) TestDecimalFieldLibrarySum() {
	ts.createBank(0, decimal.RequireFromString("1.25"))
	ts.createBank(0, decimal.RequireFromString("2.25"))

	results, err := cql.Select(
		cql.Query[models.Bank](context.Background(), ts.db),
		conditions.Bank.Rate.Aggregate().Sum().Into(func(r *rateResult) *decimal.Decimal { return &r.Total }),
	)

	ts.Require().NoError(err)
	ts.Require().Len(results, 1)
	ts.True(results[0].Total.Equal(decimal.RequireFromString("3.5")))
}

// Column-to-column decimal arithmetic (Rate + Rate = 2*Rate), exact.
func (ts *SelectIntTestSuite) TestDecimalFieldLibraryArithmetic() {
	ts.createBank(0, decimal.RequireFromString("2.5"))

	results, err := cql.Select(
		cql.Query[models.Bank](context.Background(), ts.db),
		conditions.Bank.Rate.Plus(conditions.Bank.Rate).Into(func(r *rateResult) *decimal.Decimal { return &r.Total }),
	)

	ts.Require().NoError(err)
	ts.Require().Len(results, 1)
	ts.True(results[0].Total.Equal(decimal.RequireFromString("5")))
}
