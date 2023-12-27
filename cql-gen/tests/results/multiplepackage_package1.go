// Code generated by cql-gen v0.0.5, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	package1 "github.com/FrancoLiberali/cql/cql-gen/cmd/gen/conditions/tests/multiplepackage/package1"
	package2 "github.com/FrancoLiberali/cql/cql-gen/cmd/gen/conditions/tests/multiplepackage/package2"
	model "github.com/FrancoLiberali/cql/model"
	"time"
)

func (package1Conditions package1Conditions) Package2(conditions ...condition.Condition[package2.Package2]) condition.JoinCondition[package1.Package1] {
	return condition.NewJoinCondition[package1.Package1, package2.Package2](conditions, "Package2", "ID", package1Conditions.preload(), "Package1ID", Package2.preload())
}

type package1Conditions struct {
	ID        condition.Field[package1.Package1, model.UUID]
	CreatedAt condition.Field[package1.Package1, time.Time]
	UpdatedAt condition.Field[package1.Package1, time.Time]
	DeletedAt condition.Field[package1.Package1, time.Time]
}

var Package1 = package1Conditions{
	CreatedAt: condition.Field[package1.Package1, time.Time]{Name: "CreatedAt"},
	DeletedAt: condition.Field[package1.Package1, time.Time]{Name: "DeletedAt"},
	ID:        condition.Field[package1.Package1, model.UUID]{Name: "ID"},
	UpdatedAt: condition.Field[package1.Package1, time.Time]{Name: "UpdatedAt"},
}

// Preload allows preloading the Package1 when doing a query
func (package1Conditions package1Conditions) preload() condition.Condition[package1.Package1] {
	return condition.NewPreloadCondition[package1.Package1](package1Conditions.ID, package1Conditions.CreatedAt, package1Conditions.UpdatedAt, package1Conditions.DeletedAt)
}
