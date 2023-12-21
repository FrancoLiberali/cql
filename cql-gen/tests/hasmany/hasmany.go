package hasmany

import (
	"github.com/FrancoLiberali/cql/model"
)

type Company struct {
	model.UUIDModel

	Sellers *[]Seller // Company HasMany Sellers (Company 0..1 -> 0..* Seller)
}

type Seller struct {
	model.UUIDModel

	Company   *Company
	CompanyID *model.UUID // Company HasMany Sellers (Company 0..1 -> 0..* Seller)
}
