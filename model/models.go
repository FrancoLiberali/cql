package model

import (
	"time"

	"gorm.io/gorm"
)

// supported types for model identifier
type ID interface {
	UIntID | UUID

	IsNil() bool
}

type Model interface {
	IsLoaded() bool
	SoftDeleteColumnName() string
	UpdatedAtColumnName() string
}

// Base Model for cql with uuid as id and timestamps for creation, edition and deletion (soft-delete)
//
// Every model intended to be saved in the database must embed
// UUIDModel, UUIDModelWithTimestamps, UIntModel or UIntModelWithTimestamps
// reference: https://gorm.io/docs/models.html#gorm-Model
type UUIDModel struct {
	ID UUID `gorm:"primarykey;not null"`
}

func (model UUIDModel) IsLoaded() bool {
	return !model.ID.IsNil()
}

func (model UUIDModel) SoftDeleteColumnName() string {
	return ""
}

func (model UUIDModel) UpdatedAtColumnName() string {
	return ""
}

func (model *UUIDModel) BeforeCreate(_ *gorm.DB) (err error) {
	if model.ID == NilUUID {
		model.ID = NewUUID()
	}

	return nil
}

// Base Model for cql with uuid as id and timestamps for creation, edition and deletion (soft-delete)
//
// Every model intended to be saved in the database must embed
// UUIDModel, UUIDModelWithTimestamps, UIntModel or UIntModelWithTimestamps
// reference: https://gorm.io/docs/models.html#gorm-Model
type UUIDModelWithTimestamps struct {
	ID        UUID `gorm:"primarykey;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (model UUIDModelWithTimestamps) IsLoaded() bool {
	return !model.ID.IsNil()
}

func (model UUIDModelWithTimestamps) SoftDeleteColumnName() string {
	return "deleted_at"
}

func (model UUIDModelWithTimestamps) UpdatedAtColumnName() string {
	return "updated_at"
}

func (model *UUIDModelWithTimestamps) BeforeCreate(_ *gorm.DB) (err error) {
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

// Base Model for cql with uint as id
//
// Every model intended to be saved in the database must embed
// UUIDModel, UUIDModelWithTimestamps, UIntModel or UIntModelWithTimestamps
// reference: https://gorm.io/docs/models.html#gorm-Model
type UIntModel struct {
	ID UIntID `gorm:"primarykey;not null"`
}

func (model UIntModel) IsLoaded() bool {
	return !model.ID.IsNil()
}

func (model UIntModel) SoftDeleteColumnName() string {
	return ""
}

func (model UIntModel) UpdatedAtColumnName() string {
	return ""
}

// Base Model for cql with uint as id and timestamps for creation, edition and deletion (soft-delete)
//
// Every model intended to be saved in the database must embed
// UUIDModel, UUIDModelWithTimestamps, UIntModel or UIntModelWithTimestamps
// reference: https://gorm.io/docs/models.html#gorm-Model
type UIntModelWithTimestamps struct {
	ID        UIntID `gorm:"primarykey;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (model UIntModelWithTimestamps) IsLoaded() bool {
	return !model.ID.IsNil()
}

func (model UIntModelWithTimestamps) SoftDeleteColumnName() string {
	return "deleted_at"
}

func (model UIntModelWithTimestamps) UpdatedAtColumnName() string {
	return "updated_at"
}
