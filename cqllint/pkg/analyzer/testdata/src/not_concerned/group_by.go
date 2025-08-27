package not_concerned

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

// GroupBy
// Select
// Having
// Verificar que tenga los campos seleccionados? seria otro tipo de test, que esta bueno pero complicado
// TODO groupby repetidos

var results = []struct{}{}

func testGroupBySameModel() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.Brand.Name),
	).GroupBy(
		conditions.Brand.Name,
	).Select(
		conditions.Brand.Name.Aggregate().Count(), "aggregation1",
	).Into(&results)
}

func testGroupByNotJoined() {
	cql.Query[models.Brand](
		db,
	).GroupBy(
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Select(
		conditions.Brand.Name.Aggregate().Count(), "aggregation1",
	).Into(&results)
}
