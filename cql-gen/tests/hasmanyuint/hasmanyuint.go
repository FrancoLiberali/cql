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

	Company   *CompanyUint
	CompanyID *model.UIntID // nullable FK -> childFKExtractor's NilUIntID path
}
