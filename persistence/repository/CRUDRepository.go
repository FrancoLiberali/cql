package repository

import (
	"github.com/Masterminds/squirrel"
	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/pagination"
)

// Generic CRUD Repository
type CRUDRepository[T models.Tabler, ID any] interface {
	Create(*T) httperrors.HTTPError
	Delete(*T) httperrors.HTTPError
	Save(*T) httperrors.HTTPError
	GetByID(ID) (*T, httperrors.HTTPError)
	GetAll(...pagination.SortOption) ([]*T, httperrors.HTTPError)
	Count(squirrel.Sqlizer) (uint, httperrors.HTTPError)
	Find(squirrel.Sqlizer, pagination.Paginator, ...pagination.SortOption) (*pagination.Page[T], httperrors.HTTPError)
	Transaction(fn func(CRUDRepository[T, ID]) (any, error)) (any, error)
}
