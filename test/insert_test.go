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
	).Exec()
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
	).Exec()
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

	inserted, err := cql.Insert(
		ts.db,
		product1,
		product2,
	).ExecInBatches(1)
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

func (ts *InsertIntTestSuite) TestInsertOneOnConflictAnyDoNothingThatInserts() {
	product := &models.Product{
		Int: 1,
	}

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConflict().DoNothing().Exec()
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

func (ts *InsertIntTestSuite) TestInsertOneOnConflictIDDoNothingThatInserts() {
	product := &models.Product{
		Int: 1,
	}

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConflict(conditions.Product.ID).DoNothing().Exec()
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

func (ts *InsertIntTestSuite) TestInsertOneConflictReturnsError() {
	product := ts.createProduct("", 1, 0, false, nil)
	ts.NotEmpty(product.ID)

	inserted, err := cql.Insert(
		ts.db,
		product,
	).Exec()
	ts.ErrorContains(err, "UNIQUE constraint failed: products.id")
	ts.Equal(int64(0), inserted)
}

func (ts *InsertIntTestSuite) TestInsertOneOnConflictAnyDoNothingThatConflicts() {
	product := ts.createProduct("", 1, 0, false, nil)
	ts.NotEmpty(product.ID)

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConflict().DoNothing().Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(0), inserted)

	productsReturned, err := cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(1)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 1)
}

func (ts *InsertIntTestSuite) TestInsertOneOnConflictIDDoNothingThatConflicts() {
	product := ts.createProduct("", 1, 0, false, nil)
	ts.NotEmpty(product.ID)

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConflict(conditions.Product.ID).DoNothing().Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(0), inserted)
	ts.NotEmpty(product.ID)

	productsReturned, err := cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(1)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 1)
}

func (ts *InsertIntTestSuite) TestInsertOneOnConflictUpdateAllThatInserts() {
	product := &models.Product{
		Int: 1,
	}

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConflict().UpdateAll().Exec()
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

func (ts *InsertIntTestSuite) TestInsertOneOnConflictUpdateAllThatConflicts() {
	product := ts.createProduct("", 1, 0, false, nil)
	ts.NotEmpty(product.ID)

	product.Int = 2
	product.Float = 1

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConflict().UpdateAll().Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(1), inserted)

	productsReturned, err := cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(1)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 0)

	productsReturned, err = cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(2)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 1)
}

func (ts *InsertIntTestSuite) TestInsertOneOnConflictUpdateThatInserts() {
	product := &models.Product{
		Int: 1,
	}

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConflict().Update(conditions.Product.Int).Exec()
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

func (ts *InsertIntTestSuite) TestInsertOneOnConflictUpdateThatConflicts() {
	product := ts.createProduct("", 1, 0, false, nil)
	ts.NotEmpty(product.ID)

	product.Int = 2
	product.Float = 1

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConflict().Update(conditions.Product.Int).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(1), inserted)

	productsReturned, err := cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(1)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 0)

	productsReturned, err = cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(2)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 1)

	productsReturned, err = cql.Query(
		ts.db,
		conditions.Product.Float.Is().Eq(cql.Int(0)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 1)

	productsReturned, err = cql.Query(
		ts.db,
		conditions.Product.Float.Is().Eq(cql.Int(1)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 0)
}

func (ts *InsertIntTestSuite) TestInsertOneOnConflictSetThatInserts() {
	product := &models.Product{
		Int: 1,
	}

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConflict().Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	).Exec()
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

func (ts *InsertIntTestSuite) TestInsertOneOnConflictSetThatConflicts() {
	ts.createProduct("", 3, 0, false, nil)

	product := ts.createProduct("", 1, 0, false, nil)
	ts.NotEmpty(product.ID)

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConflict().Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(1), inserted)

	productsReturned, err := cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(1)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 0)

	productsReturned, err = cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(2)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 1)
}

func (ts *InsertIntTestSuite) TestInsertOneOnConflictSetThatConflictsMultiple() {
	product1 := ts.createProduct("", 3, 0, false, nil)
	product2 := ts.createProduct("", 1, 0, false, nil)

	inserted, err := cql.Insert(
		ts.db,
		product1,
		product2,
	).OnConflict().Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(2), inserted)

	productsReturned, err := cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(1)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 0)

	productsReturned, err = cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(3)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 0)

	productsReturned, err = cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(2)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 2)
}

func (ts *InsertIntTestSuite) TestInsertOneOnConflictSetThatConflictsMultipleWithWhere() {
	product1 := ts.createProduct("", 3, 0, false, nil)
	product2 := ts.createProduct("", 1, 0, false, nil)

	inserted, err := cql.Insert(
		ts.db,
		product1,
		product2,
	).OnConflict().Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	).Where(
		conditions.Product.Int.Is().Eq(cql.Int(1)),
	).Exec()
	ts.Require().NoError(err)
	ts.Equal(int64(1), inserted)

	productsReturned, err := cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(1)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 0)

	productsReturned, err = cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(3)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 1)

	productsReturned, err = cql.Query(
		ts.db,
		conditions.Product.Int.Is().Eq(cql.Int(2)),
	).Find()
	ts.Require().NoError(err)
	ts.Len(productsReturned, 1)
}

// create from map no
// create from sql expresion si puede ser, pero es lo mismo que gormValue, asi que no, pero igual es algo que no estoy manejando bien me parece en las queries
// upser / onconflict si interesante pero meter la logica de tipos
// tiene el update all, el do nothing y el update solo algunas columnas al valor de la query o a otro valor
// insert select es donde esta lo mas interesante
// insert returning no tiene mucho sentido para el que es por objetos pero si para el que es por select

// TODO inserts con relaciones test
// TODO multiple clauses test: que pasa si ponen varias iguales? o una sin nada y despues otras -> lint posible que no voy a hacer
