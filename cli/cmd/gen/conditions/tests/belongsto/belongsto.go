package belongsto

import "github.com/ditrit/badaas/orm"

type Owner struct {
	orm.UUIDModel
}
type Owned struct {
	orm.UUIDModel

	// Owned belongsTo Owner (Owned 0..* -> 1 Owner)
	Owner   Owner
	OwnerID orm.UUID
}
