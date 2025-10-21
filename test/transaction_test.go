package test

import (
	"context"

	"github.com/stretchr/testify/assert"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

type TransactionIntTestSuite struct {
	testSuite
}

func NewTransactionIntTestSuite(
	db *cql.DB,
) *TransactionIntTestSuite {
	return &TransactionIntTestSuite{
		testSuite: testSuite{
			db: db,
		},
	}
}

func (ts *TransactionIntTestSuite) TestTransactionOK() {
	err := ts.db.Transaction(context.Background(), func(tx *cql.DB) error {
		inserted, err := cql.Insert(context.Background(), tx, &models.Product{
			String: "test",
		}).Exec()
		ts.Require().NoError(err)
		ts.Equal(int64(1), inserted)

		updated, err := cql.Update(
			context.Background(),
			tx,
			cql.True[models.Product](),
		).Set(
			conditions.Product.Int.Set().Eq(cql.Int(2)),
		)
		ts.Require().NoError(err)
		ts.Equal(int64(1), updated)

		return nil
	})
	ts.Require().NoError(err)

	productsReturned, err := cql.Query(
		context.Background(),
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(2)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 1)
}

func (ts *TransactionIntTestSuite) TestTransactionFail() {
	err := ts.db.Transaction(context.Background(), func(tx *cql.DB) error {
		inserted, err := cql.Insert(context.Background(), tx, &models.Product{
			String: "test",
		}).Exec()
		ts.Require().NoError(err)
		ts.Equal(int64(1), inserted)

		updated, err := cql.Update(
			context.Background(),
			tx,
			cql.True[models.Product](),
		).Set(
			conditions.Product.Int.Set().Eq(cql.Int(2)),
		)
		ts.Require().NoError(err)
		ts.Equal(int64(1), updated)

		return assert.AnError
	})
	ts.Require().ErrorIs(err, assert.AnError)

	productsReturned, err := cql.Query[models.Product](
		context.Background(),
		ts.db,
	).Find()
	ts.Require().NoError(err)
	ts.Empty(productsReturned)
}
