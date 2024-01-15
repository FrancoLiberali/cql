package appearance

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

var db *gorm.DB

func testQueryNotNecessary() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.IsDynamic().Eq(conditions.Phone.Name.Appearance(0).Value()), // want "Appearance call not necessary, github.com/FrancoLiberali/cql/test/models.Phone appears only once"
		),
	).Find()
}

func testQueryNotNecessaryWithFunction() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.IsDynamic().Eq(conditions.Phone.Name.Appearance(0).Value().Concat("asd")), // want "Appearance call not necessary, github.com/FrancoLiberali/cql/test/models.Phone appears only once"
		),
	).Find()
}

func testQueryNecessaryNotCalled() {
	cql.Query[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
		conditions.Child.ID.IsDynamic().Eq(conditions.ParentParent.ID.Value()), // want "github.com/FrancoLiberali/cql/test/models.ParentParent appears more than once, select which one you want to use with Appearance"
	).Find()
}

func testQueryNecessaryNotCalledWithFunction() {
	cql.Query[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
		conditions.Child.Number.IsDynamic().Eq(conditions.ParentParent.Number.Value().Plus(1)), // want "github.com/FrancoLiberali/cql/test/models.ParentParent appears more than once, select which one you want to use with Appearance"
	).Find()
}

func testQueryNecessaryCalled() {
	cql.Query[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
		conditions.Child.ID.IsDynamic().Eq(conditions.ParentParent.ID.Appearance(0).Value()),
	).Find()
}
