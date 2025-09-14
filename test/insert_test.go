package test

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
	"gorm.io/gorm"
)

type InsertIntTestSuite struct {
	testSuite
}

func NewInsertIntTestSuite(
	db *gorm.DB,
) *InsertIntTestSuite {
	return &InsertIntTestSuite{
		testSuite: testSuite{
			db: db,
		},
	}
}

func (ts *InsertIntTestSuite) TestInsertOne() {
	product := &models.Product{
		Int: 1,
	}

	inserted, err := cql.Insert(
		ts.db,
		product,
	)
	ts.Require().NoError(err)
	ts.Equal(int64(1), inserted)
	ts.NotEmpty(product.ID)

	productsReturned, err := cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(1)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 1)
}

func (ts *InsertIntTestSuite) TestInsertMultiple() {
	product1 := &models.Product{
		Int:    1,
		String: "1",
	}

	product2 := &models.Product{
		Int:    1,
		String: "2",
	}

	inserted, err := cql.Insert(
		ts.db,
		product1,
		product2,
	)
	ts.Require().NoError(err)
	ts.Equal(int64(2), inserted)
	ts.NotEmpty(product1.ID)
	ts.NotEmpty(product2.ID)

	productsReturned, err := cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(1)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 2)
}

func (ts *InsertIntTestSuite) TestInsertInBatches() {
	product1 := &models.Product{
		Int:    1,
		String: "1",
	}

	product2 := &models.Product{
		Int:    1,
		String: "2",
	}

	inserted, err := cql.InsertInBatches(
		ts.db,
		1,
		product1,
		product2,
	)
	ts.Require().NoError(err)
	ts.Equal(int64(2), inserted)
	ts.NotEmpty(product1.ID)
	ts.NotEmpty(product2.ID)

	productsReturned, err := cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(1)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 2)
}

// insert batch desde models me interesa
// create from map no
// create from sql expresion si puede ser, pero es lo mismo que gormValue, asi que no, pero igual es algo que no estoy manejando bien me parece en las queries
// upser / onconflict si interesante pero meter la logica de tipos
// tiene el update all, el do nothing y el update solo algunas columnas al valor de la query o a otro valor
