// Code generated by cql-gen v0.0.8, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	hasmany "github.com/FrancoLiberali/cql/cql-gen/cmd/gen/conditions/tests/hasmany"
	model "github.com/FrancoLiberali/cql/model"
	"time"
)

func (sellerConditions sellerConditions) Company(conditions ...condition.Condition[hasmany.Company]) condition.JoinCondition[hasmany.Seller] {
	return condition.NewJoinCondition[hasmany.Seller, hasmany.Company](conditions, "Company", "CompanyID", sellerConditions.preload(), "ID", Company.preload())
}

type sellerConditions struct {
	ID        condition.Field[hasmany.Seller, model.UUID]
	CreatedAt condition.Field[hasmany.Seller, time.Time]
	UpdatedAt condition.Field[hasmany.Seller, time.Time]
	DeletedAt condition.Field[hasmany.Seller, time.Time]
	CompanyID condition.NullableField[hasmany.Seller, model.UUID]
}

var Seller = sellerConditions{
	CompanyID: condition.NewNullableField[hasmany.Seller, model.UUID]("CompanyID", "", ""),
	CreatedAt: condition.NewField[hasmany.Seller, time.Time]("CreatedAt", "", ""),
	DeletedAt: condition.NewField[hasmany.Seller, time.Time]("DeletedAt", "", ""),
	ID:        condition.NewField[hasmany.Seller, model.UUID]("ID", "", ""),
	UpdatedAt: condition.NewField[hasmany.Seller, time.Time]("UpdatedAt", "", ""),
}

// Preload allows preloading the Seller when doing a query
func (sellerConditions sellerConditions) preload() condition.Condition[hasmany.Seller] {
	return condition.NewPreloadCondition[hasmany.Seller](sellerConditions.ID, sellerConditions.CreatedAt, sellerConditions.UpdatedAt, sellerConditions.DeletedAt, sellerConditions.CompanyID)
}
