package belongsto

import (
	"github.com/FrancoLiberali/cql/orm/model"
)

type Owner struct {
	model.UUIDModel
}
type Owned struct {
	model.UUIDModel

	// Owned belongsTo Owner (Owned 0..* -> 1 Owner)
	Owner   Owner
	OwnerID model.UUID
}
