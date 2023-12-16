package package1

import (
	"github.com/ditrit/badaas-cli/cmd/gen/conditions/tests/multiplepackage/package2"
	"github.com/ditrit/badaas/orm"
)

type Package1 struct {
	orm.UUIDModel

	Package2 package2.Package2 // Package1 HasOne Package2 (Package1 1 -> 1 Package2)
}
