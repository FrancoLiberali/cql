package orm

import (
	"time"

	"gorm.io/gorm"
)

// supported types for model identifier
type ModelID interface {
	UIntID | UUID

	IsNil() bool
}

type Model interface {
	IsLoaded() bool
}

// Base Model for gorm
//
// Every model intended to be saved in the database must embed this UUIDModel or UIntModel
// reference: https://gorm.io/docs/models.html#gorm-Model
type UUIDModel struct {
	ID        UUID `gorm:"primarykey;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (model UUIDModel) IsLoaded() bool {
	return !model.ID.IsNil()
}

func (model *UUIDModel) BeforeCreate(_ *gorm.DB) (err error) {
	if model.ID == NilUUID {
		model.ID = NewUUID()
	}

	return nil
}

type UIntID uint

const NilUIntID = 0

func (id UIntID) IsNil() bool {
	return id == NilUIntID
}

type UIntModel struct {
	ID        UIntID `gorm:"primarykey;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (model UIntModel) IsLoaded() bool {
	return !model.ID.IsNil()
}
