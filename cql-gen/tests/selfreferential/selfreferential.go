package selfreferential

import (
	"github.com/FrancoLiberali/cql/model"
)

type Employee struct {
	model.UUIDModel

	Boss   *Employee `gorm:"constraint:OnDelete:SET NULL;"` // Self-Referential Has One (Employee 0..* -> 0..1 Employee)
	BossID *model.UUID
}
