package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/pagination"
	"github.com/ditrit/badaas/services/httperrors"
)

var (
	// HTTP Error record already exists
	ErrAlreadyExists = httperrors.NewErrorNotFound("model instance", "the model instance is not in the database")
)

// Generic CRUD Repository
type CRUDRepository[T models.Tabler] interface {
	Create(*T) httperrors.HTTPError
	Delete(*T) httperrors.HTTPError
	Save(*T) httperrors.HTTPError
	GetByID(uint) (*T, httperrors.HTTPError)
	GetAll(...pagination.SortOption) ([]*T, httperrors.HTTPError)
	Count(squirrel.Sqlizer) (uint, httperrors.HTTPError)
	Find(squirrel.Sqlizer, pagination.Paginator, ...pagination.SortOption) (*pagination.Page[T], httperrors.HTTPError)
	Transaction(fn func(CRUDRepository[T]) (any, error)) (any, error)
}
