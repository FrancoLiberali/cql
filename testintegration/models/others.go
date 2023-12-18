//go:build !mysql
// +build !mysql

package models

import (
	"github.com/ditrit/badaas/orm/model"
)

type Employee struct {
	model.UUIDModel

	Name   string
	Boss   *Employee // Self-Referential Has One (Employee 0..* -> 0..1 Employee)
	BossID *model.UUID
}

func (m Employee) Equal(other Employee) bool {
	return m.Name == other.Name
}
