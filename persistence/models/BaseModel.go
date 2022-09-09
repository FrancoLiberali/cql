package models

import (
	"time"

	"gorm.io/gorm"
)

// Base Model for gorm
//
// Every model intended to be saved in the database must embed this BaseModel
// reference: https://gorm.io/docs/models.html#gorm-Model
type BaseModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
