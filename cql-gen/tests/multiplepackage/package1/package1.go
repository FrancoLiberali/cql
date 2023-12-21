package package1

import (
	"github.com/FrancoLiberali/cql/cql-gen/cmd/gen/conditions/tests/multiplepackage/package2"
	"github.com/FrancoLiberali/cql/model"
)

type Package1 struct {
	model.UUIDModel

	Package2 package2.Package2 // Package1 HasOne Package2 (Package1 1 -> 1 Package2)
}
