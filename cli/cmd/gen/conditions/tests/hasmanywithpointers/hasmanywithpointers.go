package hasmanywithpointers

import (
	"github.com/ditrit/badaas/orm"
)

type CompanyWithPointers struct {
	orm.UUIDModel

	Sellers *[]*SellerInPointers // CompanyWithPointers HasMany SellerInPointers
}

type SellerInPointers struct {
	orm.UUIDModel

	Company   *CompanyWithPointers
	CompanyID *orm.UUID // Company HasMany Seller
}
