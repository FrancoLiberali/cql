package preload

import (
	"github.com/ditrit/badaas/orm/errors"
	"github.com/ditrit/badaas/orm/model"
)

func VerifyStructLoaded[T model.Model](toVerify *T) (*T, error) {
	if toVerify == nil || !(*toVerify).IsLoaded() {
		return nil, errors.ErrRelationNotLoaded
	}

	return toVerify, nil
}

func VerifyPointerLoaded[TModel model.Model, TID model.ID](id *TID, toVerify *TModel) (*TModel, error) {
	// when the pointer to the object is nil
	// but the id pointer indicates that the relation is not nil
	if id != nil && toVerify == nil {
		return nil, errors.ErrRelationNotLoaded
	}

	return toVerify, nil
}

func VerifyPointerWithIDLoaded[TModel model.Model, TID model.ID](id TID, toVerify *TModel) (*TModel, error) {
	// when the pointer to the object is nil
	// but the id indicates that the relation is not nil
	if !id.IsNil() && toVerify == nil {
		return nil, errors.ErrRelationNotLoaded
	}

	return toVerify, nil
}

func VerifyCollectionLoaded[T model.Model](collection *[]T) ([]T, error) {
	if collection == nil {
		return nil, errors.ErrRelationNotLoaded
	}

	return *collection, nil
}
