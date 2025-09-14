package condition

import (
	"github.com/FrancoLiberali/cql/model"
	"gorm.io/gorm"
)

type Insert[T model.Model] struct {
	tx     *gorm.DB
	models []*T
}

func NewInsert[T model.Model](tx *gorm.DB, models []*T) *Insert[T] {
	return &Insert[T]{
		tx:     tx,
		models: models,
	}
}

// TODO docs
func (insert *Insert[T]) Exec() (int64, error) {
	result := insert.tx.Create(insert.models)

	return result.RowsAffected, result.Error
}

// TODO docs
func (insert *Insert[T]) ExecInBatches(batchSize int) (int64, error) {
	result := insert.tx.CreateInBatches(insert.models, batchSize)

	return result.RowsAffected, result.Error
}
