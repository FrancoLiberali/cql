// Code generated by badaas-cli v0.0.0, DO NOT EDIT.
package models

import preload "github.com/ditrit/badaas/orm/preload"

func (m Bicycle) GetOwner() (*Person, error) {
	return preload.VerifyStructLoaded[Person](&m.Owner)
}
func (m Child) GetParent1() (*Parent1, error) {
	return preload.VerifyStructLoaded[Parent1](&m.Parent1)
}
func (m Child) GetParent2() (*Parent2, error) {
	return preload.VerifyStructLoaded[Parent2](&m.Parent2)
}
func (m City) GetCountry() (*Country, error) {
	return preload.VerifyPointerWithIDLoaded[Country](m.CountryID, m.Country)
}
func (m Company) GetSellers() ([]Seller, error) {
	return preload.VerifyCollectionLoaded[Seller](m.Sellers)
}
func (m Country) GetCapital() (*City, error) {
	return preload.VerifyStructLoaded[City](&m.Capital)
}
func (m Employee) GetBoss() (*Employee, error) {
	return preload.VerifyPointerLoaded[Employee](m.BossID, m.Boss)
}
func (m Parent1) GetParentParent() (*ParentParent, error) {
	return preload.VerifyStructLoaded[ParentParent](&m.ParentParent)
}
func (m Parent2) GetParentParent() (*ParentParent, error) {
	return preload.VerifyStructLoaded[ParentParent](&m.ParentParent)
}
func (m Phone) GetBrand() (*Brand, error) {
	return preload.VerifyStructLoaded[Brand](&m.Brand)
}
func (m Sale) GetProduct() (*Product, error) {
	return preload.VerifyStructLoaded[Product](&m.Product)
}
func (m Sale) GetSeller() (*Seller, error) {
	return preload.VerifyPointerLoaded[Seller](m.SellerID, m.Seller)
}
func (m Seller) GetCompany() (*Company, error) {
	return preload.VerifyPointerLoaded[Company](m.CompanyID, m.Company)
}
func (m Seller) GetUniversity() (*University, error) {
	return preload.VerifyPointerLoaded[University](m.UniversityID, m.University)
}
