// Code generated by cql-gen v0.0.10, DO NOT EDIT.
package hasmanywithpointers

import preload "github.com/FrancoLiberali/cql/preload"

func (m CompanyWithPointers) GetSellers() ([]*SellerInPointers, error) {
	return preload.VerifyCollectionLoaded[*SellerInPointers](m.Sellers)
}
func (m SellerInPointers) GetCompany() (*CompanyWithPointers, error) {
	return preload.VerifyPointerLoaded[CompanyWithPointers](m.CompanyID, m.Company)
}
