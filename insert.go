package cql

import (
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/model"
)

// TODO docs
func Insert[T model.Model](tx *gorm.DB, models ...*T) (int64, error) {
	result := tx.Create(models)

	return result.RowsAffected, result.Error
}

// TODO docs
func InsertInBatches[T model.Model](tx *gorm.DB, batchSize int, models ...*T) (int64, error) {
	result := tx.CreateInBatches(models, batchSize)

	return result.RowsAffected, result.Error
}
