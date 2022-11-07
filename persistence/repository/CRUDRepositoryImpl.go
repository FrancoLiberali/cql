package repository

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/gormdatabase"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/pagination"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Return a database error
func DatabaseError(message string, golangError error) httperrors.HTTPError {
	return httperrors.NewInternalServerError(
		"database error",
		message,
		golangError,
	)
}

// Implementation of the Generic CRUD Repository
type CRUDRepositoryImpl[T models.Tabler] struct {
	CRUDRepository[T]
	gormDatabase *gorm.DB
	logger       *zap.Logger
}

// Contructor of the Generic CRUD Repository
func NewCRUDRepository[T models.Tabler](database *gorm.DB, logger *zap.Logger) CRUDRepository[T] {
	return &CRUDRepositoryImpl[T]{gormDatabase: database, logger: logger}
}

// Run the function passed as parameter, if it returns the error and rollback the transaction.
// If no error is returned, it commits the transaction and return the interface{} value.
func (repository *CRUDRepositoryImpl[T]) Transaction(transactionFunction func(CRUDRepository[T]) (any, error)) (any, error) {
	transaction := repository.gormDatabase.Begin()
	defer func() {
		if recoveredError := recover(); recoveredError != nil {
			transaction.Rollback()
		}
	}()
	returnValue, err := transactionFunction(&CRUDRepositoryImpl[T]{gormDatabase: transaction})
	if err != nil {
		transaction.Rollback()
		return nil, err
	}
	return returnValue, transaction.Commit().Error
}

// Create an entity of a Model
func (repository *CRUDRepositoryImpl[T]) Create(entity *T) httperrors.HTTPError {
	err := repository.gormDatabase.Create(entity).Error
	if err != nil {
		if gormdatabase.IsDuplicateKeyError(err) {
			return httperrors.NewHTTPError(
				http.StatusConflict,
				fmt.Sprintf("%T already exist in database", entity),
				"",
				nil, false)
		}
		return DatabaseError(
			fmt.Sprintf("could not create  %v in %s", entity, (*entity).TableName()),
			err,
		)

	}
	return nil
}

// Delete an entity of a Model
func (repository *CRUDRepositoryImpl[T]) Delete(entity *T) httperrors.HTTPError {
	err := repository.gormDatabase.Delete(entity).Error
	if err != nil {
		return DatabaseError(
			fmt.Sprintf("could not delete %v in %s", entity, (*entity).TableName()),
			err,
		)
	}
	return nil
}

// Save an entity of a Model
func (repository *CRUDRepositoryImpl[T]) Save(entity *T) httperrors.HTTPError {
	err := repository.gormDatabase.Save(entity).Error
	if err != nil {
		return DatabaseError(
			fmt.Sprintf("could not save user %v in %s", entity, (*entity).TableName()),
			err,
		)
	}
	return nil
}

// Get an entity of a Model By ID
func (repository *CRUDRepositoryImpl[T]) GetByID(id uint) (*T, httperrors.HTTPError) {
	var entity T
	transaction := repository.gormDatabase.First(&entity, "id = ?", id)
	if transaction.Error != nil {
		return nil, DatabaseError(
			fmt.Sprintf("could not get %s by id %d", entity.TableName(), id),
			transaction.Error,
		)
	}
	return &entity, nil
}

// Get all entities of a Model
func (repository *CRUDRepositoryImpl[T]) GetAll(sortOptions ...pagination.SortOption) ([]*T, httperrors.HTTPError) {
	var entities []*T
	transaction := repository.gormDatabase
	for _, sortOption := range sortOptions {
		transaction = transaction.Order(sortOption.ToClause())
	}
	transaction.Find(&entities)
	if transaction.Error != nil {
		var emptyInstanceForError T
		return nil, DatabaseError(
			fmt.Sprintf("could not get %s", emptyInstanceForError.TableName()),
			transaction.Error,
		)
	}
	return entities, nil
}

// Count entities of a models
func (repository *CRUDRepositoryImpl[T]) Count(filters squirrel.Sqlizer) (uint, httperrors.HTTPError) {
	whereClause, values, httpError := repository.compileSQL(filters)
	if httpError != nil {
		return 0, httpError
	}
	return repository.count(whereClause, values)
}

// Count the number of record that match the where clause with the provided values on the db
func (repository *CRUDRepositoryImpl[T]) count(whereClause string, values []interface{}) (uint, httperrors.HTTPError) {
	var entity *T
	var count int64
	transaction := repository.gormDatabase.Model(entity).Where(whereClause, values).Count(&count)
	if transaction.Error != nil {
		var emptyInstanceForError T
		return 0, DatabaseError(
			fmt.Sprintf("could not count data from %s with condition %q", emptyInstanceForError.TableName(), whereClause),
			transaction.Error,
		)
	}
	return uint(count), nil
}

// Find entities of a Model
func (repository *CRUDRepositoryImpl[T]) Find(
	filters squirrel.Sqlizer,
	page pagination.Paginator,
	sortOptions ...pagination.SortOption,
) (*pagination.Page[T], httperrors.HTTPError) {
	transaction := repository.gormDatabase.Begin()
	defer func() {
		if recoveredError := recover(); recoveredError != nil {
			transaction.Rollback()

		}
	}()
	var instances []*T
	whereClause, values, httpError := repository.compileSQL(filters)

	if httpError != nil {
		return nil, httpError
	}
	if page == nil {
		transaction = transaction.Where(whereClause, values...).Find(&instances)
	} else {
		transaction = transaction.
			Offset(
				int((page.Offset() - 1) * page.Limit()),
			).
			Limit(
				int(page.Limit()),
			)
	}
	for _, sortOption := range sortOptions {
		transaction = transaction.Order(sortOption.ToClause())
	}
	transaction = transaction.Where(whereClause, values...).Find(&instances)
	if transaction.Error != nil {
		transaction.Rollback()
		var emptyInstanceForError T
		return nil, DatabaseError(
			fmt.Sprintf("could not get data from %s with condition %q", emptyInstanceForError.TableName(), whereClause),
			transaction.Error,
		)
	}
	// Get Count
	nbElem, httpError := repository.count(whereClause, values)
	if httpError != nil {
		transaction.Rollback()
		return nil, httpError
	}
	err := transaction.Commit().Error
	if err != nil {
		return nil, DatabaseError(
			"transaction failed to commit", err)
	}
	return pagination.NewPage(instances, page.Offset(), page.Limit(), nbElem), nil
}

// compile the sql where clause
func (repository *CRUDRepositoryImpl[T]) compileSQL(filters squirrel.Sqlizer) (string, []interface{}, httperrors.HTTPError) {
	compiledSQLString, values, err := filters.ToSql()
	if err != nil {
		return "", []interface{}{}, httperrors.NewInternalServerError(
			"sql error",
			fmt.Sprintf("Failed to build the sql request (condition=%v)", filters),
			err,
		)
	}
	return compiledSQLString, values, nil
}
