package overridereferencesinverse

import "github.com/FrancoLiberali/cql/model"

type Computer struct {
	model.UUIDModel
	Name      string
	Processor Processor `gorm:"foreignKey:ComputerName;references:Name"`
}

type Processor struct {
	model.UUIDModel
	ComputerName string
}
