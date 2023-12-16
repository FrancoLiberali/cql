package hasmany

import "github.com/ditrit/badaas/orm"

type Company struct {
	orm.UUIDModel

	Sellers *[]Seller // Company HasMany Sellers (Company 0..1 -> 0..* Seller)
}

type Seller struct {
	orm.UUIDModel

	Company   *Company
	CompanyID *orm.UUID // Company HasMany Sellers (Company 0..1 -> 0..* Seller)
}
