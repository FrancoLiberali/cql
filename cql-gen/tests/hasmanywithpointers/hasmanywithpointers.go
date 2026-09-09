package hasmanywithpointers

import (
	"github.com/FrancoLiberali/cql/model"
)

type CompanyWithPointers struct {
	model.UUIDModel

	Sellers *[]*SellerInPointers // CompanyWithPointers HasMany SellerInPointers
}

type SellerInPointers struct {
	model.UUIDModel

	// FK named after the parent model, per gorm's has-many convention, so the
	// generated has-many loader and conditions agree (cql-gen derives the FK
	// as <ParentModel>ID).
	CompanyWithPointers   *CompanyWithPointers
	CompanyWithPointersID *model.UUID
}
