package not_concerned

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/condition"
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
	)
}

func testGroupByJoinedModel() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		conditions.Phone.Name,
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}
func testGroupByJoinedWithJoinedWithCondition() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("asd")),
		),
	).GroupBy(
		conditions.Phone.Name,
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupByJoinedWithJoinedWithPreload() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand().Preload(),
	).GroupBy(
		conditions.Phone.Name,
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupByJoinedWithJoinedWithConditionsWithPreload() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("asd")),
		).Preload(),
	).GroupBy(
		conditions.Phone.Name,
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupByJoinedModelInVariable() {
	value := conditions.Brand.Name

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		value,
	)
}

func testGroupByNotJoined() {
	cql.Query[models.Brand](
		db,
	).GroupBy(
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)
}

func testGroupByNotJoinedInVariable() {
	value := conditions.City.Name // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		value,
	)
}

func testGroupByJoinedConditionInList() {
	values := []condition.IField{
		conditions.Phone.Name,
	}

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		values...,
	)
}

func testGroupByNotJoinedConditionInList() {
	values := []condition.IField{
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	}

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		values...,
	)
}

func testGroupByJoinedConditionInListWithAppend() {
	values := []condition.IField{}

	values = append(
		values,
		conditions.Phone.Name,
	)

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		values...,
	)
}

func testGroupByNotJoinedConditionInListWithAppend() {
	values := []condition.IField{}

	values = append(
		values,
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		values...,
	)
}

func testGroupByNotJoinedConditionInListWithAppendSecond() {
	values := []condition.IField{}

	values = append(
		values,
		conditions.Phone.Name,
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		values...,
	)
}

func testGroupByNotJoinedConditionInListWithAppendMultiple() {
	values := []condition.IField{}

	values = append(
		values,
		conditions.Phone.Name,
	)

	values = append(
		values,
		conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	)

	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
	).GroupBy(
		values...,
	)
}
