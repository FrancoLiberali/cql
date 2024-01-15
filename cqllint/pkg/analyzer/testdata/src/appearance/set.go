package appearance

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

func testSetNotNecessary() {
	cql.Update[models.Phone](
		db,
		conditions.Phone.Name.Is().Eq("asd"),
	).Set(
		conditions.Phone.Name.Set().Dynamic(conditions.Phone.Name.Appearance(0).Value()), // want "Appearance call not necessary, github.com/FrancoLiberali/cql/test/models.Phone appears only once"
	)
}

func testSetNotNecessaryWithFunction() {
	cql.Update[models.Phone](
		db,
		conditions.Phone.Name.Is().Eq("asd"),
	).Set(
		conditions.Phone.Name.Set().Dynamic(conditions.Phone.Name.Appearance(0).Value().Concat("asd")), // want "Appearance call not necessary, github.com/FrancoLiberali/cql/test/models.Phone appears only once"
	)
}

func testSetNecessaryNotCalled() {
	cql.Update[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
	).Set(
		conditions.Child.Name.Set().Dynamic(conditions.ParentParent.Name.Value()), // want "github.com/FrancoLiberali/cql/test/models.ParentParent appears more than once, select which one you want to use with Appearance"
	)
}

func testSetNecessaryNotCalledWithFunction() {
	cql.Update[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
	).Set(
		conditions.Child.Name.Set().Dynamic(conditions.ParentParent.Name.Value().Concat("asd")), // want "github.com/FrancoLiberali/cql/test/models.ParentParent appears more than once, select which one you want to use with Appearance"
	)
}

func testSetNecessaryCalled() {
	cql.Update[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
	).Set(
		conditions.Child.Name.Set().Dynamic(conditions.ParentParent.Name.Appearance(0).Value()),
	)
}
