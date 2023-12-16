package repository

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/squirrel"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/gormdatabase"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/pagination"
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
type CRUDRepositoryImpl[T models.Tabler, ID any] struct {
	CRUDRepository[T, ID]
	gormDatabase            *gorm.DB
	logger                  *zap.Logger
	paginationConfiguration configuration.PaginationConfiguration
}

// Constructor of the Generic CRUD Repository
func NewCRUDRepository[T models.Tabler, ID any](
	database *gorm.DB,
	logger *zap.Logger,
	paginationConfiguration configuration.PaginationConfiguration,
) CRUDRepository[T, ID] {
	return &CRUDRepositoryImpl[T, ID]{
		gormDatabase:            database,
		logger:                  logger,
		paginationConfiguration: paginationConfiguration,
	}
}

// Run the function passed as parameter, if it returns the error and rollback the transaction.
// If no error is returned, it commits the transaction and return the interface{} value.
func (repository *CRUDRepositoryImpl[T, ID]) Transaction(transactionFunction func(CRUDRepository[T, ID]) (any, error)) (any, httperrors.HTTPError) {
	transaction := repository.gormDatabase.Begin()
	defer func() {
		if recoveredError := recover(); recoveredError != nil {
			transaction.Rollback()
		}
	}()
	returnValue, err := transactionFunction(&CRUDRepositoryImpl[T, ID]{gormDatabase: transaction})
	if err != nil {
		transaction.Rollback()
		return nil, DatabaseError("transaction failed", err)
	}
	err = transaction.Commit().Error
	if err != nil {
		return nil, DatabaseError("transaction failed to commit", err)
	}
	return returnValue, nil
}

// Create an entity of a Model
func (repository *CRUDRepositoryImpl[T, ID]) Create(entity *T) httperrors.HTTPError {
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
func (repository *CRUDRepositoryImpl[T, ID]) Delete(entity *T) httperrors.HTTPError {
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
func (repository *CRUDRepositoryImpl[T, ID]) Save(entity *T) httperrors.HTTPError {
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
func (repository *CRUDRepositoryImpl[T, ID]) GetByID(id ID) (*T, httperrors.HTTPError) {
	var entity T
	transaction := repository.gormDatabase.First(&entity, "id = ?", id)
	if transaction.Error != nil {
		return nil, DatabaseError(
			fmt.Sprintf("could not get %s by id %v", entity.TableName(), id),
			transaction.Error,
		)
	}
	return &entity, nil
}

// Get all entities of a Model
func (repository *CRUDRepositoryImpl[T, ID]) GetAll(sortOption SortOption) ([]*T, httperrors.HTTPError) {
	var entities []*T
	transaction := repository.gormDatabase
	if sortOption != nil {
		transaction = transaction.Order(buildClauseFromSortOption(sortOption))
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

// Build a gorm order clause from a SortOption
func buildClauseFromSortOption(sortOption SortOption) clause.OrderByColumn {
	return clause.OrderByColumn{Column: clause.Column{Name: sortOption.Column()}, Desc: sortOption.Desc()}
}

// Count entities of a models
func (repository *CRUDRepositoryImpl[T, ID]) Count(filters squirrel.Sqlizer) (uint, httperrors.HTTPError) {
	whereClause, values, httpError := repository.compileSQL(filters)
	if httpError != nil {
		return 0, httpError
	}
	return repository.count(whereClause, values)
}

// Count the number of record that match the where clause with the provided values on the db
func (repository *CRUDRepositoryImpl[T, ID]) count(whereClause string, values []interface{}) (uint, httperrors.HTTPError) {
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
func (repository *CRUDRepositoryImpl[T, ID]) Find(
	filters squirrel.Sqlizer,
	page pagination.Paginator,
	sortOption SortOption,
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
	if page != nil {
		transaction = transaction.
			Offset(
				int((page.Offset() - 1) * page.Limit()),
			).
			Limit(
				int(page.Limit()),
			)
	} else {
		page = pagination.NewPaginator(0, repository.paginationConfiguration.GetMaxElemPerPage())
	}
	if sortOption != nil {
		transaction = transaction.Order(buildClauseFromSortOption(sortOption))
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
func (repository *CRUDRepositoryImpl[T, ID]) compileSQL(filters squirrel.Sqlizer) (string, []interface{}, httperrors.HTTPError) {
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
