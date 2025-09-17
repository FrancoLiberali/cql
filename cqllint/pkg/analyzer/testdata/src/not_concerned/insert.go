package not_concerned

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

func testOnConflictSetStatic() {
	cql.Insert(
		db,
		&models.Product{},
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	).Exec()
}

func testOnConflictSetSameModel() {
	cql.Insert(
		db,
		&models.Product{},
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.Int.Set().Eq(conditions.Product.Float),
	).Exec()
}

func testOnConflictSetDifferentModel() {
	cql.Insert(
		db,
		&models.Product{},
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.String.Set().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}
