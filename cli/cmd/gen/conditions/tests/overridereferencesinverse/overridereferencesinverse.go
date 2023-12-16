package overridereferencesinverse

import "github.com/ditrit/badaas/orm"

type Computer struct {
	orm.UUIDModel
	Name      string
	Processor Processor `gorm:"foreignKey:ComputerName;references:Name"`
}

type Processor struct {
	orm.UUIDModel
	ComputerName string
}
