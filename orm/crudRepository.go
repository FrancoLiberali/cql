package orm

import (
	"errors"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Generic CRUD Repository
// T can be any model whose identifier attribute is of type ID
type CRUDRepository[T Model, ID ModelID] interface {
	// Create model "model" inside transaction "tx"
	Create(tx *gorm.DB, entity *T) error

	// ----- read -----
	// Get a model by its ID
	GetByID(tx *gorm.DB, id ID) (*T, error)

	// Get only one model that match "conditions" inside transaction "tx"
	// or returns error if 0 or more than 1 are found.
	QueryOne(tx *gorm.DB, conditions ...Condition[T]) (*T, error)

	// Get the list of models that match "conditions" inside transaction "tx"
	Query(tx *gorm.DB, conditions ...Condition[T]) ([]*T, error)

	// Save model "model" inside transaction "tx"
	Save(tx *gorm.DB, entity *T) error

	// Delete model "model" inside transaction "tx"
	Delete(tx *gorm.DB, entity *T) error
}

var (
	ErrMoreThanOneObjectFound = errors.New("found more that one object that meet the requested conditions")
	ErrObjectNotFound         = errors.New("no object exists that meets the requested conditions")
)

// Implementation of the Generic CRUD Repository
type CRUDRepositoryImpl[T Model, ID ModelID] struct {
	CRUDRepository[T, ID]
}

// Constructor of the Generic CRUD Repository
func NewCRUDRepository[T Model, ID ModelID]() CRUDRepository[T, ID] {
	return &CRUDRepositoryImpl[T, ID]{}
}

// Create object "entity" inside transaction "tx"
func (repository *CRUDRepositoryImpl[T, ID]) Create(tx *gorm.DB, entity *T) error {
	return tx.Create(entity).Error
}

// Delete object "entity" inside transaction "tx"
func (repository *CRUDRepositoryImpl[T, ID]) Delete(tx *gorm.DB, entity *T) error {
	return tx.Delete(entity).Error
}

// Save object "entity" inside transaction "tx"
func (repository *CRUDRepositoryImpl[T, ID]) Save(tx *gorm.DB, entity *T) error {
	return tx.Save(entity).Error
}

// Get a model by its ID
func (repository *CRUDRepositoryImpl[T, ID]) GetByID(tx *gorm.DB, id ID) (*T, error) {
	var model T

	err := tx.First(&model, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &model, nil
}

// Get only one model that match "conditions" inside transaction "tx"
// or returns error if 0 or more than 1 are found.
func (repository *CRUDRepositoryImpl[T, ID]) QueryOne(tx *gorm.DB, conditions ...Condition[T]) (*T, error) {
	entities, err := repository.Query(tx, conditions...)
	if err != nil {
		return nil, err
	}

	switch {
	case len(entities) == 1:
		return entities[0], nil
	case len(entities) == 0:
		return nil, ErrObjectNotFound
	default:
		return nil, ErrMoreThanOneObjectFound
	}
}

// Get the list of models that match "conditions" inside transaction "tx"
func (repository *CRUDRepositoryImpl[T, ID]) Query(tx *gorm.DB, conditions ...Condition[T]) ([]*T, error) {
	query, err := applyConditionsToQuery(tx, conditions)
	if err != nil {
		return nil, err
	}

	// execute query
	var entities []*T
	err = query.Find(&entities).Error

	return entities, err
}

func applyConditionsToQuery[T Model](query *gorm.DB, conditions []Condition[T]) (*gorm.DB, error) {
	initialTableName, err := getTableName(query, *new(T))
	if err != nil {
		return nil, err
	}

	initialTable := Table{
		Name:    initialTableName,
		Alias:   initialTableName,
		Initial: true,
	}

	query = query.Select(initialTableName + ".*")
	for _, condition := range conditions {
		query, err = condition.ApplyTo(query, initialTable)
		if err != nil {
			return nil, err
		}
	}

	return query, nil
}

// Get the name of the table in "db" in which the data for "entity" is saved
// returns error is table name can not be found by gorm,
// probably because the type of "entity" is not registered using AddModel
func getTableName(db *gorm.DB, entity any) (string, error) {
	schemaName, err := schema.Parse(entity, &sync.Map{}, db.NamingStrategy)
	if err != nil {
		return "", err
	}

	return schemaName.Table, nil
}
