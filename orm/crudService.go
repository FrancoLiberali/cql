package orm

import (
	"gorm.io/gorm"
)

// T can be any model whose identifier attribute is of type ID
type CRUDService[T any, ID BadaasID] interface {
	// Get the model of type T that has the "id"
	GetByID(id ID) (*T, error)

	// Get only one model that match "conditions"
	// or return error if 0 or more than 1 are found.
	QueryOne(conditions ...Condition[T]) (*T, error)

	// Get the list of models that match "conditions"
	Query(conditions ...Condition[T]) ([]*T, error)
}

// check interface compliance
var _ CRUDService[UUIDModel, UUID] = (*crudServiceImpl[UUIDModel, UUID])(nil)

// Implementation of the CRUD Service
type crudServiceImpl[T any, ID BadaasID] struct {
	CRUDService[T, ID]
	db         *gorm.DB
	repository CRUDRepository[T, ID]
}

func NewCRUDService[T any, ID BadaasID](
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
func (service *crudServiceImpl[T, ID]) QueryOne(conditions ...Condition[T]) (*T, error) {
	return service.repository.QueryOne(service.db, conditions...)
}

// Get the list of models that match "conditions"
func (service *crudServiceImpl[T, ID]) Query(conditions ...Condition[T]) ([]*T, error) {
	return service.repository.Query(service.db, conditions...)
}
