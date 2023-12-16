package hasone

import "github.com/ditrit/badaas/orm"

type Company struct {
	orm.UUIDModel

	Sellers []Seller // Company HasMany Sellers (Company 0..1 -> 0..* Seller)
}

type Seller struct {
	orm.UUIDModel

	CompanyID *orm.UUID // Company HasMany Sellers (Company 0..1 -> 0..* Seller)
}
