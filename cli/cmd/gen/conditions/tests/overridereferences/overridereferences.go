package overridereferences

import "github.com/FrancoLiberali/cql/orm/model"

type Brand struct {
	model.UUIDModel

	Name string `gorm:"unique;type:VARCHAR(255)"`
}

type Phone struct {
	model.UUIDModel

	// Bicycle BelongsTo Person (Bicycle 0..* -> 1 Person)
	Brand     Brand `gorm:"references:Name;foreignKey:BrandName"`
	BrandName string
}
