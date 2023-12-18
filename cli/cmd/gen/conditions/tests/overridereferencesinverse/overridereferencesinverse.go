package overridereferencesinverse

import "github.com/ditrit/badaas/orm/model"

type Computer struct {
	model.UUIDModel
	Name      string
	Processor Processor `gorm:"foreignKey:ComputerName;references:Name"`
}

type Processor struct {
	model.UUIDModel
	ComputerName string
}
