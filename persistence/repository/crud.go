package repository

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/condition"
	"github.com/ditrit/badaas/orm/model"
)

// Generic CRUD Repository
// T can be any model whose identifier attribute is of type ID
type CRUD[T model.Model, ID model.ID] interface {
	// Create model "model" inside transaction "tx"
	Create(tx *gorm.DB, entity *T) error

	// ----- read -----
	// Get a model by its ID
	GetByID(tx *gorm.DB, id ID) (*T, error)

	// Get the list of models that match "conditions" inside transaction "tx"
	Find(tx *gorm.DB, conditions ...condition.Condition[T]) ([]*T, error)

	// Get the only one model that match "conditions" inside transaction "tx"
	// or returns error if 0 or more than 1 are found.
	FindOne(tx *gorm.DB, conditions ...condition.Condition[T]) (*T, error)

	// Save model "model" inside transaction "tx"
	Save(tx *gorm.DB, entity *T) error

	// Delete model "model" inside transaction "tx"
	Delete(tx *gorm.DB, entity *T) error
}

// Implementation of the Generic CRUD Repository
type crudImpl[T model.Model, ID model.ID] struct {
	CRUD[T, ID]
}

// Constructor of the Generic CRUD Repository
func NewCRUD[T model.Model, ID model.ID]() CRUD[T, ID] {
	return &crudImpl[T, ID]{}
}

// Create model "model" inside transaction "tx"
func (repository *crudImpl[T, ID]) Create(tx *gorm.DB, model *T) error {
	return tx.Create(model).Error
}

// Delete model "model" inside transaction "tx"
func (repository *crudImpl[T, ID]) Delete(tx *gorm.DB, model *T) error {
	return tx.Delete(model).Error
}

// Save model "model" inside transaction "tx"
func (repository *crudImpl[T, ID]) Save(tx *gorm.DB, model *T) error {
	return tx.Save(model).Error
}

// Get a model by its ID
func (repository *crudImpl[T, ID]) GetByID(tx *gorm.DB, id ID) (*T, error) {
	var model T

	err := tx.First(&model, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &model, nil
}

// Get the list of models that match "conditions" inside transaction "tx"
func (repository *crudImpl[T, ID]) Find(tx *gorm.DB, conditions ...condition.Condition[T]) ([]*T, error) {
	return orm.NewQuery[T](tx, conditions...).Find()
}

// Get the only one model that match "conditions" inside transaction "tx"
// or returns error if 0 or more than 1 are found.
func (repository *crudImpl[T, ID]) FindOne(tx *gorm.DB, conditions ...condition.Condition[T]) (*T, error) {
	return orm.NewQuery[T](tx, conditions...).FindOne()
}
