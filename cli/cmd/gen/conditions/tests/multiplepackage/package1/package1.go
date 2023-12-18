package package1

import (
	"github.com/FrancoLiberali/cql/orm/model"
	"github.com/ditrit/badaas-cli/cmd/gen/conditions/tests/multiplepackage/package2"
)

type Package1 struct {
	model.UUIDModel

	Package2 package2.Package2 // Package1 HasOne Package2 (Package1 1 -> 1 Package2)
}
