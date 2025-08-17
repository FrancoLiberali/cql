package appearance

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

func testSetNotNecessary() {
	cql.Update[models.Phone](
		db,
		conditions.Phone.Name.Is().Eq(cql.String("asd")),
	).Set(
		conditions.Phone.Name.Set().Eq(conditions.Phone.Name.Appearance(0)), // want "Appearance call not necessary, github.com/FrancoLiberali/cql/test/models.Phone appears only once"
	)
}

func testSetNotNecessaryWithFunction() {
	cql.Update[models.Phone](
		db,
		conditions.Phone.Name.Is().Eq(cql.String("asd")),
	).Set(
		conditions.Phone.Name.Set().Eq(conditions.Phone.Name.Appearance(0).Concat("asd")), // want "Appearance call not necessary, github.com/FrancoLiberali/cql/test/models.Phone appears only once"
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
		conditions.Child.Name.Set().Eq(conditions.ParentParent.Name), // want "github.com/FrancoLiberali/cql/test/models.ParentParent appears more than once, select which one you want to use with Appearance"
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
		conditions.Child.Name.Set().Eq(conditions.ParentParent.Name.Concat("asd")), // want "github.com/FrancoLiberali/cql/test/models.ParentParent appears more than once, select which one you want to use with Appearance"
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
		conditions.Child.Name.Set().Eq(conditions.ParentParent.Name.Appearance(0)),
	)
}

func testSetOutOfRange() {
	cql.Update[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
	).Set(
		conditions.Child.Name.Set().Eq(conditions.ParentParent.Name.Appearance(2)), // want "selected appearance is bigger than github.com/FrancoLiberali/cql/test/models.ParentParent's number of appearances"
	)
}
