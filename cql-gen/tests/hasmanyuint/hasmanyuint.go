package hasmanyuint

import "github.com/FrancoLiberali/cql/model"

// CompanyUint is a UInt-keyed (model.UIntModel) HasMany parent. The UUID-keyed
// hasmany fixtures don't reach the parentPKIsUUID=false / parentPKTypeName
// (UIntID) / childFKExtractor NilUIntID branches of the has-many loader
// generator; this fixture does.
type CompanyUint struct {
	model.UIntModel

	Sellers []SellerUint
}

type SellerUint struct {
	model.UIntModel

	// FK named after the parent model (CompanyUint -> CompanyUintID), per
	// gorm's has-many convention, so the generated loader and conditions agree.
	CompanyUint   *CompanyUint
	CompanyUintID *model.UIntID // nullable FK -> childFKExtractor's NilUIntID path
}
