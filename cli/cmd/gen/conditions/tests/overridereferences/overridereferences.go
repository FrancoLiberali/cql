package overridereferences

import "github.com/ditrit/badaas/orm"

type Brand struct {
	orm.UUIDModel

	Name string `gorm:"unique;type:VARCHAR(255)"`
}

type Phone struct {
	orm.UUIDModel

	// Bicycle BelongsTo Person (Bicycle 0..* -> 1 Person)
	Brand     Brand `gorm:"references:Name;foreignKey:BrandName"`
	BrandName string
}
