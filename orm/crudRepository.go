package orm

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm/condition"
	"github.com/ditrit/badaas/orm/errors"
	"github.com/ditrit/badaas/orm/model"
)

// Generic CRUD Repository
// T can be any model whose identifier attribute is of type ID
type CRUDRepository[T model.Model, ID model.ID] interface {
	// Create model "model" inside transaction "tx"
	Create(tx *gorm.DB, entity *T) error

	// ----- read -----
	// Get a model by its ID
	GetByID(tx *gorm.DB, id ID) (*T, error)

	// Get only one model that match "conditions" inside transaction "tx"
	// or returns error if 0 or more than 1 are found.
	QueryOne(tx *gorm.DB, conditions ...condition.Condition[T]) (*T, error)

	// Get the list of models that match "conditions" inside transaction "tx"
	Query(tx *gorm.DB, conditions ...condition.Condition[T]) ([]*T, error)

	// Save model "model" inside transaction "tx"
	Save(tx *gorm.DB, entity *T) error

	// Delete model "model" inside transaction "tx"
	Delete(tx *gorm.DB, entity *T) error
}

// Implementation of the Generic CRUD Repository
type crudRepositoryImpl[T model.Model, ID model.ID] struct {
	CRUDRepository[T, ID]
}

// Constructor of the Generic CRUD Repository
func NewCRUDRepository[T model.Model, ID model.ID]() CRUDRepository[T, ID] {
	return &crudRepositoryImpl[T, ID]{}
}

// Create model "model" inside transaction "tx"
func (repository *crudRepositoryImpl[T, ID]) Create(tx *gorm.DB, model *T) error {
	return tx.Create(model).Error
}

// Delete model "model" inside transaction "tx"
func (repository *crudRepositoryImpl[T, ID]) Delete(tx *gorm.DB, model *T) error {
	return tx.Delete(model).Error
}

// Save model "model" inside transaction "tx"
func (repository *crudRepositoryImpl[T, ID]) Save(tx *gorm.DB, model *T) error {
	return tx.Save(model).Error
}

// Get a model by its ID
func (repository *crudRepositoryImpl[T, ID]) GetByID(tx *gorm.DB, id ID) (*T, error) {
	var model T

	err := tx.First(&model, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &model, nil
}

// Get only one model that match "conditions" inside transaction "tx"
// or returns error if 0 or more than 1 are found.
func (repository *crudRepositoryImpl[T, ID]) QueryOne(tx *gorm.DB, conditions ...condition.Condition[T]) (*T, error) {
	entities, err := repository.Query(tx, conditions...)
	if err != nil {
		return nil, err
	}

	switch {
	case len(entities) == 1:
		return entities[0], nil
	case len(entities) == 0:
		return nil, errors.ErrObjectNotFound
	default:
		return nil, errors.ErrMoreThanOneObjectFound
	}
}

// Get the list of models that match "conditions" inside transaction "tx"
func (repository *crudRepositoryImpl[T, ID]) Query(tx *gorm.DB, conditions ...condition.Condition[T]) ([]*T, error) {
	query, err := condition.ApplyConditions(tx, conditions)
	if err != nil {
		return nil, err
	}

	// execute query
	var entities []*T
	err = query.Find(&entities)

	return entities, err
}
