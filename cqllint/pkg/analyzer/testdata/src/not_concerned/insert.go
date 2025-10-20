package not_concerned

import (
	"context"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

func testOnConflictSetStatic() {
	cql.Insert(
		context.Background(),
		db,
		&models.Product{},
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	).Exec()
}

func testOnConflictSetStaticIndex() {
	cql.Insert[models.Product](
		context.Background(),
		db,
		&models.Product{},
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	).Exec()
}

func testOnConflictSetSameModel() {
	cql.Insert(
		context.Background(),
		db,
		&models.Product{},
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.Int.Set().Eq(conditions.Product.Float),
	).Exec()
}

func testOnConflictSetSameModelVar() {
	product := &models.Product{}

	cql.Insert(
		context.Background(),
		db,
		product,
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.Int.Set().Eq(conditions.Product.Float),
	).Exec()
}

func testOnConflictSetSameModelIndex() {
	cql.Insert[models.Product](
		context.Background(),
		db,
		&models.Product{},
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.Int.Set().Eq(conditions.Product.Float),
	).Exec()
}

func testOnConflictSetDifferentModel() {
	cql.Insert(
		context.Background(),
		db,
		&models.Product{},
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.String.Set().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testOnConflictSetDifferentModelVar() {
	product := &models.Product{}

	cql.Insert(
		context.Background(),
		db,
		product,
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.String.Set().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testOnConflictSetDifferentModelIndex() {
	cql.Insert[models.Product](
		context.Background(),
		db,
		&models.Product{},
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.String.Set().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testOnConflictSetWhereStatic() {
	cql.Insert(
		context.Background(),
		db,
		&models.Product{},
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	).Where(
		conditions.Product.Int.Is().Eq(cql.Int(2)),
	).Exec()
}

func testOnConflictSetWhereStaticIndex() {
	cql.Insert[models.Product](
		context.Background(),
		db,
		&models.Product{},
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	).Where(
		conditions.Product.Int.Is().Eq(cql.Int(2)),
	).Exec()
}

func testOnConflictSetWhereSameModel() {
	cql.Insert(
		context.Background(),
		db,
		&models.Product{},
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	).Where(
		conditions.Product.Int.Is().Eq(conditions.Product.Float),
	).Exec()
}

func testOnConflictSetWhereSameModelVar() {
	product := &models.Product{}

	cql.Insert(
		context.Background(),
		db,
		product,
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	).Where(
		conditions.Product.Int.Is().Eq(conditions.Product.Float),
	).Exec()
}

func testOnConflictSetWhereSameModelIndex() {
	cql.Insert[models.Product](
		context.Background(),
		db,
		&models.Product{},
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.Int.Set().Eq(cql.Int(2)),
	).Where(
		conditions.Product.Int.Is().Eq(conditions.Product.Float),
	).Exec()
}

func testOnConflictSetWhereDifferentModel() {
	cql.Insert(
		context.Background(),
		db,
		&models.Product{},
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.String.Set().Eq(cql.String("asd")),
	).Where(
		conditions.Product.String.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testOnConflictSetWhereDifferentModelVar() {
	product := &models.Product{}

	cql.Insert(
		context.Background(),
		db,
		product,
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.String.Set().Eq(cql.String("asd")),
	).Where(
		conditions.Product.String.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}

func testOnConflictSetWhereDifferentModelIndex() {
	cql.Insert[models.Product](
		context.Background(),
		db,
		&models.Product{},
	).OnConflictOn(conditions.Product.ID).Set(
		conditions.Product.String.Set().Eq(cql.String("asd")),
	).Where(
		conditions.Product.String.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Exec()
}
