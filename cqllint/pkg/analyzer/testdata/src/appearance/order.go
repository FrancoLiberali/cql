package appearance

import (
	"github.com/FrancoLiberali/cql"
	"github.com/FrancoLiberali/cql/test/conditions"
	"github.com/FrancoLiberali/cql/test/models"
)

func testOrderNotNecessary() {
	cql.Query[models.Brand](
		db,
	).Descending(conditions.Brand.Name.Appearance(0)).Find() // want "Appearance call not necessary, github.com/FrancoLiberali/cql/test/models.Brand appears only once"
}

func testOrderNecessaryNotCalled() {
	cql.Query[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
	).Descending(conditions.ParentParent.ID).Find() // want "github.com/FrancoLiberali/cql/test/models.ParentParent appears more than once, select which one you want to use with Appearance"
}

func testOrderNecessaryCalled() {
	cql.Query[models.Child](
		db,
		conditions.Child.Parent1(
			conditions.Parent1.ParentParent(),
		),
		conditions.Child.Parent2(
			conditions.Parent2.ParentParent(),
		),
	).Descending(conditions.ParentParent.ID.Appearance(0)).Find()
}
