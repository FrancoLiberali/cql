package selfreferential

import "github.com/ditrit/badaas/orm"

type Employee struct {
	orm.UUIDModel

	Boss   *Employee `gorm:"constraint:OnDelete:SET NULL;"` // Self-Referential Has One (Employee 0..* -> 0..1 Employee)
	BossID *orm.UUID
}
