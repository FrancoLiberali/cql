package appearance

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

func testQueryNotNecessary() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name.Appearance(0)), // want "Appearance call not necessary, github.com/FrancoLiberali/cql/test/models.Phone appears only once"
		),
	).Find()
}

func testQueryNotNecessaryWithFunction() {
	cql.Query[models.Phone](
		db,
		conditions.Phone.Brand(
			conditions.Brand.Name.Is().Eq(conditions.Phone.Name.Appearance(0).Concat(cql.String("asd"))), // want "Appearance call not necessary, github.com/FrancoLiberali/cql/test/models.Phone appears only once"
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
		conditions.Child.ID.Is().Eq(conditions.ParentParent.ID), // want "github.com/FrancoLiberali/cql/test/models.ParentParent appears more than once, select which one you want to use with Appearance"
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
		conditions.Child.Number.Is().Eq(conditions.ParentParent.Number.Plus(cql.Int(1))), // want "github.com/FrancoLiberali/cql/test/models.ParentParent appears more than once, select which one you want to use with Appearance"
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
		conditions.Child.ID.Is().Eq(conditions.ParentParent.ID.Appearance(0)),
	).Find()
}

func testQueryOutOfRange() {
	cql.Query[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
		conditions.Child.ID.Is().Eq(conditions.ParentParent.ID.Appearance(2)), // want "selected appearance is bigger than github.com/FrancoLiberali/cql/test/models.ParentParent's number of appearances"
	).Find()
}
