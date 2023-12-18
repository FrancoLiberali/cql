package package2

import "github.com/ditrit/badaas/orm"

type Package2 struct {
	orm.UUIDModel

	Package1ID orm.UUID // Package1 HasOne Package2 (Package1 1 -> 1 Package2)
}
