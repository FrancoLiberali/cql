// Code generated by cql-gen v0.1.0, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	belongsto "github.com/FrancoLiberali/cql/cql-gen/cmd/gen/conditions/tests/belongsto"
	model "github.com/FrancoLiberali/cql/model"
	"time"
)

type ownerConditions struct {
	ID        condition.Field[belongsto.Owner, model.UUID]
	CreatedAt condition.Field[belongsto.Owner, time.Time]
	UpdatedAt condition.Field[belongsto.Owner, time.Time]
	DeletedAt condition.Field[belongsto.Owner, time.Time]
}

var Owner = ownerConditions{
	CreatedAt: condition.NewField[belongsto.Owner, time.Time]("CreatedAt", "", ""),
	DeletedAt: condition.NewField[belongsto.Owner, time.Time]("DeletedAt", "", ""),
	ID:        condition.NewField[belongsto.Owner, model.UUID]("ID", "", ""),
	UpdatedAt: condition.NewField[belongsto.Owner, time.Time]("UpdatedAt", "", ""),
}

// Preload allows preloading the Owner when doing a query
func (ownerConditions ownerConditions) preload() condition.Condition[belongsto.Owner] {
	return condition.NewPreloadCondition[belongsto.Owner](ownerConditions.ID, ownerConditions.CreatedAt, ownerConditions.UpdatedAt, ownerConditions.DeletedAt)
}
