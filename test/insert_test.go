package test

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/sql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
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

func (ts *InsertIntTestSuite) TestInsertOneOnConflictDoNothingThatInserts() {
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
	).OnConflictOn(conditions.Product.ID).DoNothing().Exec()

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: OnConflictOn")
	default:
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
}

func (ts *InsertIntTestSuite) TestInsertOneConflictReturnsError() {
	product := ts.createProduct("", 1, 0, false, nil)
	ts.NotEmpty(product.ID)

	inserted, err := cql.Insert(
		ts.db,
		product,
	).Exec()
	ts.Equal(int64(0), inserted)

	switch getDBDialector() {
	case sql.Postgres:
		ts.ErrorContains(err, `duplicate key value violates unique constraint "products_pkey" (SQLSTATE 23505)`)
	case sql.MySQL:
		ts.ErrorContains(err, `Error 1062 (23000): Duplicate entry `)
		ts.ErrorContains(err, `for key 'products.PRIMARY'`)
	case sql.SQLServer:
		ts.ErrorContains(err, `mssql: Violation of PRIMARY KEY constraint 'PK__products__`)
		ts.ErrorContains(err, `'. Cannot insert duplicate key in object 'dbo.products'. The duplicate key value is `)
	case sql.SQLite:
		ts.ErrorContains(err, "UNIQUE constraint failed: products.id")
	}
}

func (ts *InsertIntTestSuite) TestInsertOneOnConflictDoNothingThatConflicts() {
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
	).OnConflictOn(conditions.Product.ID).DoNothing().Exec()

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: OnConflictOn")
	default:
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
}

func (ts *InsertIntTestSuite) TestInsertOneOnConstraintDoNothingThatInserts() {
	product := &models.Product{
		Int: 1,
	}

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConstraint("products_pkey").DoNothing().Exec()

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer, sql.SQLite:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: OnConstraint")
	case sql.Postgres:
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
}

func (ts *InsertIntTestSuite) TestInsertOneOnConstraintDoNothingThatConflicts() {
	product := ts.createProduct("", 1, 0, false, nil)
	ts.NotEmpty(product.ID)

	product.Int = 2
	product.Float = 1

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConstraint("products_pkey").DoNothing().Exec()

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer, sql.SQLite:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: OnConstraint")
	case sql.Postgres:
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
}

func (ts *InsertIntTestSuite) TestInsertOneOnConflictUpdateAllThatInserts() {
	product := &models.Product{
		Int: 1,
	}

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConflict().UpdateAll().Exec()

	switch getDBDialector() {
	case sql.Postgres:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: UpdateAll after OnConflict")
	default:
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
}

func (ts *InsertIntTestSuite) TestInsertOneOnConflictOnUpdateAllThatInserts() {
	product := &models.Product{
		Int: 1,
	}

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConflictOn(conditions.Product.ID).UpdateAll().Exec()

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: OnConflictOn")
	default:
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

	switch getDBDialector() {
	case sql.Postgres:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: UpdateAll after OnConflict")
	default:
		switch getDBDialector() {
		case sql.MySQL:
			ts.checkUpdateAllThatConflicts(int64(2), inserted, err)
		default:
			ts.checkUpdateAllThatConflicts(int64(1), inserted, err)
		}
	}
}

func (ts *InsertIntTestSuite) TestInsertOneOnConflictOnUpdateAllThatConflicts() {
	product := ts.createProduct("", 1, 0, false, nil)
	ts.NotEmpty(product.ID)

	product.Int = 2
	product.Float = 1

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConflictOn(conditions.Product.ID).UpdateAll().Exec()

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: OnConflictOn")
	default:
		switch getDBDialector() {
		case sql.MySQL:
			ts.checkUpdateAllThatConflicts(int64(2), inserted, err)
		default:
			ts.checkUpdateAllThatConflicts(int64(1), inserted, err)
		}
	}
}

func (ts *InsertIntTestSuite) TestInsertOneOnConstraintUpdateAllThatInserts() {
	product := &models.Product{
		Int: 1,
	}

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConstraint("products_pkey").UpdateAll().Exec()

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer, sql.SQLite:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: OnConstraint")
	case sql.Postgres:
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
}

func (ts *InsertIntTestSuite) TestInsertOneOnConstraintUpdateAllThatConflicts() {
	product := ts.createProduct("", 1, 0, false, nil)
	ts.NotEmpty(product.ID)

	product.Int = 2
	product.Float = 1

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConstraint("products_pkey").UpdateAll().Exec()

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer, sql.SQLite:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: OnConstraint")
	case sql.Postgres:
		ts.checkUpdateAllThatConflicts(int64(1), inserted, err)
	}
}

func (ts *InsertIntTestSuite) checkUpdateAllThatConflicts(expectedInserted, inserted int64, err error) {
	ts.T().Helper()

	ts.Require().NoError(err)
	ts.Equal(expectedInserted, inserted)

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

	var inserted int64

	var err error

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer:
		inserted, err = cql.Insert(
			ts.db,
			product,
		).OnConflict().Update(conditions.Product.Int).Exec()
	default:
		inserted, err = cql.Insert(
			ts.db,
			product,
		).OnConflictOn(conditions.Product.ID).Update(conditions.Product.Int).Exec()
	}

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

	var inserted int64

	var err error

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer:
		inserted, err = cql.Insert(
			ts.db,
			product,
		).OnConflict().Update(conditions.Product.Int).Exec()
	default:
		inserted, err = cql.Insert(
			ts.db,
			product,
		).OnConflictOn(conditions.Product.ID).Update(conditions.Product.Int).Exec()
	}

	switch getDBDialector() {
	case sql.MySQL:
		ts.checkUpdateThatConflicts(int64(2), inserted, err)
	default:
		ts.checkUpdateThatConflicts(int64(1), inserted, err)
	}
}

func (ts *InsertIntTestSuite) TestInsertOneOnConstraintUpdateThatInserts() {
	product := &models.Product{
		Int: 1,
	}

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConstraint("products_pkey").Update(conditions.Product.Int).Exec()

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer, sql.SQLite:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: OnConstraint")
	case sql.Postgres:
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
}

func (ts *InsertIntTestSuite) TestInsertOneOnConstraintUpdateThatConflicts() {
	product := ts.createProduct("", 1, 0, false, nil)
	ts.NotEmpty(product.ID)

	product.Int = 2
	product.Float = 1

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConstraint("products_pkey").Update(conditions.Product.Int).Exec()

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer, sql.SQLite:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: OnConstraint")
	case sql.Postgres:
		ts.checkUpdateThatConflicts(int64(1), inserted, err)
	}
}

func (ts *InsertIntTestSuite) checkUpdateThatConflicts(expectedInserted, inserted int64, err error) {
	ts.T().Helper()

	ts.Require().NoError(err)
	ts.Equal(expectedInserted, inserted)

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

	var inserted int64

	var err error

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer:
		inserted, err = cql.Insert(
			ts.db,
			product,
		).OnConflict().Set(
			conditions.Product.Int.Set().Eq(cql.Int(2)),
		).Exec()
	default:
		inserted, err = cql.Insert(
			ts.db,
			product,
		).OnConflictOn(conditions.Product.ID).Set(
			conditions.Product.Int.Set().Eq(cql.Int(2)),
		).Exec()
	}

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

	var inserted int64

	var err error

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer:
		inserted, err = cql.Insert(
			ts.db,
			product,
		).OnConflict().Set(
			conditions.Product.Int.Set().Eq(cql.Int(2)),
		).Exec()
	default:
		inserted, err = cql.Insert(
			ts.db,
			product,
		).OnConflictOn(conditions.Product.ID).Set(
			conditions.Product.Int.Set().Eq(cql.Int(2)),
		).Exec()
	}

	switch getDBDialector() {
	case sql.MySQL:
		ts.checkSetThatConflicts(int64(2), inserted, err)
	default:
		ts.checkSetThatConflicts(int64(1), inserted, err)
	}
}

func (ts *InsertIntTestSuite) TestInsertOneOnConstraintSetThatInserts() {
	product := &models.Product{
		Int: 1,
	}

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConstraint("products_pkey").Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	).Exec()

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer, sql.SQLite:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: OnConstraint")
	case sql.Postgres:
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
}

func (ts *InsertIntTestSuite) TestInsertOneOnConstraintSetThatConflicts() {
	product := ts.createProduct("", 1, 0, false, nil)
	ts.NotEmpty(product.ID)

	inserted, err := cql.Insert(
		ts.db,
		product,
	).OnConstraint("products_pkey").Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	).Exec()

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer, sql.SQLite:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: OnConstraint")
	case sql.Postgres:
		ts.checkSetThatConflicts(int64(1), inserted, err)
	}
}

func (ts *InsertIntTestSuite) checkSetThatConflicts(expectedInserted, inserted int64, err error) {
	ts.T().Helper()

	ts.Require().NoError(err)
	ts.Equal(expectedInserted, inserted)

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

func (ts *InsertIntTestSuite) TestInsertOneOnConflictSetThatConflictsDynamic() {
	ts.createProduct("", 3, 0, false, nil)

	product := ts.createProduct("", 1, 1, false, nil)
	ts.NotEmpty(product.ID)

	var inserted int64

	var err error

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer:
		inserted, err = cql.Insert(
			ts.db,
			product,
		).OnConflict().Set(
			// TODO aca tambien necesita linter aunque no seria necesario realmente
			conditions.Product.Int.Set().Eq(conditions.Product.Int.Plus(cql.Int(1))),
		).Exec()
	default:
		inserted, err = cql.Insert(
			ts.db,
			product,
		).OnConflictOn(conditions.Product.ID).Set(
			// TODO aca tambien necesita linter aunque no seria necesario realmente
			conditions.Product.Int.Set().Eq(conditions.Product.Int.Plus(cql.Int(1))),
		).Exec()
	}

	ts.Require().NoError(err)

	switch getDBDialector() {
	case sql.MySQL:
		ts.Equal(int64(2), inserted)
	default:
		ts.Equal(int64(1), inserted)
	}

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

	var inserted int64

	var err error

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer:
		inserted, err = cql.Insert(
			ts.db,
			product1,
			product2,
		).OnConflict().Set(
			conditions.Product.Int.Set().Eq(cql.Int(2)),
		).Exec()
	default:
		inserted, err = cql.Insert(
			ts.db,
			product1,
			product2,
		).OnConflictOn(conditions.Product.ID).Set(
			conditions.Product.Int.Set().Eq(cql.Int(2)),
		).Exec()
	}

	ts.Require().NoError(err)

	switch getDBDialector() {
	case sql.MySQL:
		ts.Equal(int64(4), inserted)
	default:
		ts.Equal(int64(2), inserted)
	}

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
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	).Where(
		conditions.Product.Int.Is().Eq(cql.Int(1)),
	).Exec()

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer:
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: Where")
	default:
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
}

func (ts *InsertIntTestSuite) TestInsertOneOnConflictSetThatConflictsMultipleWithWhereDynamic() {
	product1 := ts.createProduct("", 3, 0, false, nil)
	product2 := ts.createProduct("", 1, 0, false, nil)

	inserted, err := cql.Insert(
		ts.db,
		product1,
		product2,
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	).Where(
		// TODO aca tambien necesita linter aunque no seria necesario realmente
		conditions.Product.Int.Is().Eq(conditions.Product.Float.Plus(cql.Int(1))),
	).Exec()

	switch getDBDialector() {
	case sql.MySQL, sql.SQLServer:
		// Where is not supported by mysql
		ts.ErrorIs(err, cql.ErrUnsupportedByDatabase)
		ts.ErrorContains(err, "method: Where")
	default:
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
}
