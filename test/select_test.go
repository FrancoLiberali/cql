package test

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

type SelectIntTestSuite struct {
	testSuite
}

func NewSelectIntTestSuite(
	db *gorm.DB,
) *SelectIntTestSuite {
	return &SelectIntTestSuite{
		testSuite: testSuite{
			db: db,
		},
	}
}

// TODO hacer lo mismo para los selects del groupby

func (ts *SelectIntTestSuite) TestSelect() {
	ts.createProduct("1", 1, 0, false, nil)
	ts.createProduct("2", 1, 1, false, nil)
	ts.createProduct("5", 0, 2, false, nil)

	results, err := cql.Select(
		cql.Query[models.Product](
			ts.db,
		).Descending(conditions.Product.Int),
		// TODO test de compilacion de esto
		// value tiene que ser del tipo correcto
		// todos los result del mismo tipo
		cql.ValueInto(conditions.Product.Int, func(value float64, result *ResultInt) {
			result.Int = int(value)
		}),
		cql.ValueInto(conditions.Product.Int.Plus(1), func(value float64, result *ResultInt) {
			result.Aggregation1 = int(value)
		}),
	)

	ts.Require().NoError(err)
	EqualList(&ts.Suite, []ResultInt{
		{Int: 0, Aggregation1: 1},
		{Int: 1, Aggregation1: 2},
		{Int: 1, Aggregation1: 2},
	}, results)
}

// TODO
// function
// multiple
// joined
