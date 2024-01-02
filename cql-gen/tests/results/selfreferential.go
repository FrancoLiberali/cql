// Code generated by cql-gen v0.0.5, DO NOT EDIT.
package conditions

import (
	selfreferential "github.com/FrancoLiberali/cql-gen/cmd/gen/conditions/tests/selfreferential"
	condition "github.com/FrancoLiberali/cql/condition"
	model "github.com/FrancoLiberali/cql/model"
	"time"
)

func (employeeConditions employeeConditions) Boss(conditions ...condition.Condition[selfreferential.Employee]) condition.JoinCondition[selfreferential.Employee] {
	return condition.NewJoinCondition[selfreferential.Employee, selfreferential.Employee](conditions, "Boss", "BossID", employeeConditions.preload(), "ID", Employee.preload())
}

type employeeConditions struct {
	ID        condition.Field[selfreferential.Employee, model.UUID]
	CreatedAt condition.Field[selfreferential.Employee, time.Time]
	UpdatedAt condition.Field[selfreferential.Employee, time.Time]
	DeletedAt condition.Field[selfreferential.Employee, time.Time]
	BossID    condition.NullableField[selfreferential.Employee, model.UUID]
}

var Employee = employeeConditions{
	BossID:    condition.NullableField[selfreferential.Employee, model.UUID]{UpdatableField: condition.UpdatableField[selfreferential.Employee, model.UUID]{Field: condition.Field[selfreferential.Employee, model.UUID]{Name: "BossID"}}},
	CreatedAt: condition.Field[selfreferential.Employee, time.Time]{Name: "CreatedAt"},
	DeletedAt: condition.Field[selfreferential.Employee, time.Time]{Name: "DeletedAt"},
	ID:        condition.Field[selfreferential.Employee, model.UUID]{Name: "ID"},
	UpdatedAt: condition.Field[selfreferential.Employee, time.Time]{Name: "UpdatedAt"},
}

// Preload allows preloading the Employee when doing a query
func (employeeConditions employeeConditions) preload() condition.Condition[selfreferential.Employee] {
	return condition.NewPreloadCondition[selfreferential.Employee](employeeConditions.ID, employeeConditions.CreatedAt, employeeConditions.UpdatedAt, employeeConditions.DeletedAt, employeeConditions.BossID)
}
