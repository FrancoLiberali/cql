package a

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
		conditions.Brand.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "conditions.City is not joined by the query"
	).Find()
}

func testNotJoinedInDifferentLines() {
	cql.Query[models.Brand](
		db,
		conditions.Brand.Name.IsDynamic().Eq(
			conditions.City.Name.Value(), // want "conditions.City is not joined by the query"
		),
	).Find()
}

func testNotJoinedInsideJoinCondition() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "conditions.City is not joined by the query"
		),
	).Find()
}

func testJoinedWithMainModel() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.IsDynamic().Eq(conditions.Phone.Name.Value()),
		),
	).Find()
}

func testNotJoinedInsideNestedJoinCondition() {
	cql.Query[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(
				conditions.ParentParent.Name.IsDynamic().Eq(conditions.City.Name.Value()), // want "conditions.City is not joined by the query"
			),
		),
	).Find()
}

func testJoinedInsideNestedJoinConditionWithMainModel() {
	cql.Query[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(
				conditions.ParentParent.Name.IsDynamic().Eq(conditions.Child.Name.Value()),
			),
		),
	).Find()
}

func testJoinedInsideNestedJoinConditionWithPreviousJoin() {
	cql.Query[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(
				conditions.ParentParent.Name.IsDynamic().Eq(conditions.Parent1.Name.Value()),
			),
		),
	).Find()
}

func testJoinedWithJoinedWithCondition() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq("asd"),
		),
		conditions.Phone.Name.IsDynamic().Eq(conditions.Brand.Name.Value()),
	).Find()
}

func testJoinedWithJoinedWithPreload() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand().Preload(),
		conditions.Phone.Name.IsDynamic().Eq(conditions.Brand.Name.Value()),
	).Find()
}

func testNotJoinedWithJoinedWithConditionBefore() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Name.IsDynamic().Eq(conditions.Brand.Name.Value()), // want "conditions.Brand is not joined by the query"
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq("asd"),
		),
	).Find()
}
