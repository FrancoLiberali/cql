package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base Model for gorm
//
// Every model intended to be saved in the database must embed this BaseModel
// reference: https://gorm.io/docs/models.html#gorm-Model
type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
