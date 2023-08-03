package orm

import "errors"

var ErrRelationNotLoaded = errors.New("relation not loaded")

func VerifyStructLoaded[T Model](toVerify *T) (*T, error) {
	if toVerify == nil || !(*toVerify).IsLoaded() {
		return nil, ErrRelationNotLoaded
	}

	return toVerify, nil
}

func VerifyPointerLoaded[TModel Model, TID ModelID](id *TID, toVerify *TModel) (*TModel, error) {
	// when the pointer to the object is nil
	// but the id pointer indicates that the relation is not nil
	if id != nil && toVerify == nil {
		return nil, ErrRelationNotLoaded
	}

	return toVerify, nil
}

func VerifyPointerWithIDLoaded[TModel Model, TID ModelID](id TID, toVerify *TModel) (*TModel, error) {
	// when the pointer to the object is nil
	// but the id indicates that the relation is not nil
	if !id.IsNil() && toVerify == nil {
		return nil, ErrRelationNotLoaded
	}

	return toVerify, nil
}
