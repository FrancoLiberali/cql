package hasmanywithpointers

import (
	"github.com/FrancoLiberali/cql/orm/model"
)

type CompanyWithPointers struct {
	model.UUIDModel

	Sellers *[]*SellerInPointers // CompanyWithPointers HasMany SellerInPointers
}

type SellerInPointers struct {
	model.UUIDModel

	Company   *CompanyWithPointers
	CompanyID *model.UUID // Company HasMany Seller
}
