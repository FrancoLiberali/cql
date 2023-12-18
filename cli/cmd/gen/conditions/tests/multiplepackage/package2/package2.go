package package2

import (
	"github.com/FrancoLiberali/cql/orm/model"
)

type Package2 struct {
	model.UUIDModel

	Package1ID model.UUID // Package1 HasOne Package2 (Package1 1 -> 1 Package2)
}
