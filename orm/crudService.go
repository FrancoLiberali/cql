package orm

import (
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm/condition"
	"github.com/ditrit/badaas/orm/model"
)

// T can be any model whose identifier attribute is of type ID
type CRUDService[T model.Model, ID model.ID] interface {
	// Get the model of type T that has the "id"
	GetByID(id ID) (*T, error)

	// Get only one model that match "conditions"
	// or return error if 0 or more than 1 are found.
	QueryOne(conditions ...condition.Condition[T]) (*T, error)

	// Get the list of models that match "conditions"
	Query(conditions ...condition.Condition[T]) ([]*T, error)
}

// check interface compliance
var _ CRUDService[model.UUIDModel, model.UUID] = (*crudServiceImpl[model.UUIDModel, model.UUID])(nil)

// Implementation of the CRUD Service
type crudServiceImpl[T model.Model, ID model.ID] struct {
	CRUDService[T, ID]
	db         *gorm.DB
	repository CRUDRepository[T, ID]
}

func NewCRUDService[T model.Model, ID model.ID](
	db *gorm.DB,
	repository CRUDRepository[T, ID],
) CRUDService[T, ID] {
	return &crudServiceImpl[T, ID]{
		db:         db,
		repository: repository,
	}
}

// Get the model of type T that has the "id"
func (service *crudServiceImpl[T, ID]) GetByID(id ID) (*T, error) {
	return service.repository.GetByID(service.db, id)
}

// Get only one model that match "conditions"
// or return error if 0 or more than 1 are found.
func (service *crudServiceImpl[T, ID]) QueryOne(conditions ...condition.Condition[T]) (*T, error) {
	return service.repository.QueryOne(service.db, conditions...)
}

// Get the list of models that match "conditions"
func (service *crudServiceImpl[T, ID]) Query(conditions ...condition.Condition[T]) ([]*T, error) {
	return service.repository.Query(service.db, conditions...)
}
