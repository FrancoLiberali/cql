// Code generated by cql-gen v0.1.0, DO NOT EDIT.
package conditions

import (
	condition "github.com/FrancoLiberali/cql/condition"
	hasmanywithpointers "github.com/FrancoLiberali/cql/cql-gen/cmd/gen/conditions/tests/hasmanywithpointers"
	model "github.com/FrancoLiberali/cql/model"
	"time"
)

type companyWithPointersConditions struct {
	ID        condition.Field[hasmanywithpointers.CompanyWithPointers, model.UUID]
	CreatedAt condition.Field[hasmanywithpointers.CompanyWithPointers, time.Time]
	UpdatedAt condition.Field[hasmanywithpointers.CompanyWithPointers, time.Time]
	DeletedAt condition.Field[hasmanywithpointers.CompanyWithPointers, time.Time]
	Sellers   condition.Collection[hasmanywithpointers.CompanyWithPointers, hasmanywithpointers.SellerInPointers]
}

var CompanyWithPointers = companyWithPointersConditions{
	CreatedAt: condition.NewField[hasmanywithpointers.CompanyWithPointers, time.Time]("CreatedAt", "", ""),
	DeletedAt: condition.NewField[hasmanywithpointers.CompanyWithPointers, time.Time]("DeletedAt", "", ""),
	ID:        condition.NewField[hasmanywithpointers.CompanyWithPointers, model.UUID]("ID", "", ""),
	Sellers:   condition.NewCollection[hasmanywithpointers.CompanyWithPointers, hasmanywithpointers.SellerInPointers]("Sellers", "ID", "CompanyWithPointersID"),
	UpdatedAt: condition.NewField[hasmanywithpointers.CompanyWithPointers, time.Time]("UpdatedAt", "", ""),
}

// Preload allows preloading the CompanyWithPointers when doing a query
func (companyWithPointersConditions companyWithPointersConditions) preload() condition.Condition[hasmanywithpointers.CompanyWithPointers] {
	return condition.NewPreloadCondition[hasmanywithpointers.CompanyWithPointers](companyWithPointersConditions.ID, companyWithPointersConditions.CreatedAt, companyWithPointersConditions.UpdatedAt, companyWithPointersConditions.DeletedAt)
}
