package overrideforeignkey

import "github.com/ditrit/badaas/orm/model"

type Person struct {
	model.UUIDModel
}

type Bicycle struct {
	model.UUIDModel

	// Bicycle BelongsTo Person (Bicycle 0..* -> 1 Person)
	Owner            Person `gorm:"foreignKey:OwnerSomethingID"`
	OwnerSomethingID string
}
