package overrideforeignkey

import "github.com/ditrit/badaas/orm"

type Person struct {
	orm.UUIDModel
}

type Bicycle struct {
	orm.UUIDModel

	// Bicycle BelongsTo Person (Bicycle 0..* -> 1 Person)
	Owner            Person `gorm:"foreignKey:OwnerSomethingID"`
	OwnerSomethingID string
}
