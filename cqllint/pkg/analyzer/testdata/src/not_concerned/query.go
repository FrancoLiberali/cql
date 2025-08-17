package not_concerned

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

var db *gorm.DB

func testNotJoinedInSameLine() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testNotJoinedInDifferentLines() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.Is().Eq(
			conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testNotJoinedWithTrue() {
	cql.Query[models.Brand](
		db,
		cql.True[models.Brand](),
		conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testNotJoinedInsideConnector() {
	cql.Query[models.Brand](
		db,
		cql.And(
			conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testNotJoinedInsideJoinCondition() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testJoinedWithMainModel() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testNotJoinedInsideNestedJoinCondition() {
	cql.Query[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(
				conditions.ParentParent.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
			),
		),
	).Find()
}

func testJoinedInsideNestedJoinConditionWithMainModel() {
	cql.Query[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(
				conditions.ParentParent.Name.Is().Eq(conditions.Child.Name),
			),
		),
		conditions.Child.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedInsideNestedJoinConditionWithPreviousJoin() {
	cql.Query[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(
				conditions.ParentParent.Name.Is().Eq(conditions.Parent1.Name),
			),
		),
		conditions.Child.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedWithJoinedWithCondition() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("asd")),
		),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedWithJoinedWithoutCondition() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedWithJoinedWithPreload() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand().Preload(),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedWithJoinedWithConditionsWithPreload() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("asd")),
		).Preload(),
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name),
		conditions.Phone.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testNotJoinedWithJoinedWithConditionBefore() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Name.Is().Eq(conditions.Brand.Name), // want "github.com/FrancoLiberali/cql/test/models.Brand is not joined by the query"
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(cql.String("asd")),
		),
	).Find()
}

func testJoinedWithDifferentRelationNameWithConditionsUsesConditionName() {
	cql.Query[models.Bicycle](
		db,
		conditions.Bicycle.Owner(
			conditions.Person.Name.Is().Eq(cql.String("asd")),
		),
		conditions.Bicycle.Name.Is().Eq(conditions.Person.Name),
		conditions.Bicycle.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedWithDifferentRelationNameWithConditionsWithPreloadUsesConditionName() {
	cql.Query[models.Bicycle](
		db,
		conditions.Bicycle.Owner(
			conditions.Person.Name.Is().Eq(cql.String("asd")),
		).Preload(),
		conditions.Bicycle.Name.Is().Eq(conditions.Person.Name),
		conditions.Bicycle.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedWithDifferentRelationNameWithoutConditions() {
	cql.Query[models.Bicycle](
		db,
		conditions.Bicycle.Owner(),
		conditions.Bicycle.Name.Is().Eq(conditions.Person.Name),
		conditions.Bicycle.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testJoinedWithDifferentRelationNameWithoutConditionsWithPreload() {
	cql.Query[models.Bicycle](
		db,
		conditions.Bicycle.Owner().Preload(),
		conditions.Bicycle.Name.Is().Eq(conditions.Person.Name),
		conditions.Bicycle.Name.Is().Eq(conditions.City.Name), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
	).Find()
}

func testNotJoinedWithAppearance() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.City.Name.Appearance(0)), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testNotJoinedWithFunction() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.City.Name.Concat("asd")), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testNotJoinedWithTwoFunctions() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.City.Name.Concat("asd").Concat("asd")), // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}

func testMultipleArgumentsFirstNotJoined() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
		conditions.Phone.Name.Is().Between(
			conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
			conditions.Brand.Name,
		),
	).Find()
}

func testMultipleArgumentsSecondNotJoined() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name),
		),
		conditions.Phone.Name.Is().Between(
			conditions.Brand.Name,
			conditions.City.Name, // want "github.com/FrancoLiberali/cql/test/models.City is not joined by the query"
		),
	).Find()
}
