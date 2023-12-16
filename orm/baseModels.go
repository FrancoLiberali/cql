package orm

import (
	"time"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

// supported types for model identifier
type BadaasID interface {
	uint | UUID
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

func (model *UUIDModel) BeforeCreate(tx *gorm.DB) (err error) {
	if model.ID == UUID(uuid.Nil) {
		model.ID = UUID(uuid.New())
	}
	return nil
}

type UIntModel gorm.Model
