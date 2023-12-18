package preload

import (
	"github.com/FrancoLiberali/cql/model"
)

func VerifyStructLoaded[T model.Model](toVerify *T) (*T, error) {
	if toVerify == nil || !(*toVerify).IsLoaded() {
		return nil, ErrRelationNotLoaded
	}

	return toVerify, nil
}

func VerifyPointerLoaded[TModel model.Model, TID model.ID](id *TID, toVerify *TModel) (*TModel, error) {
	// when the pointer to the object is nil
	// but the id pointer indicates that the relation is not nil
	if id != nil && toVerify == nil {
		return nil, ErrRelationNotLoaded
	}

	return toVerify, nil
}

func VerifyPointerWithIDLoaded[TModel model.Model, TID model.ID](id TID, toVerify *TModel) (*TModel, error) {
	// when the pointer to the object is nil
	// but the id indicates that the relation is not nil
	if !id.IsNil() && toVerify == nil {
		return nil, ErrRelationNotLoaded
	}

	return toVerify, nil
}

func VerifyCollectionLoaded[T model.Model](collection *[]T) ([]T, error) {
	if collection == nil {
		return nil, ErrRelationNotLoaded
	}

	return *collection, nil
}
