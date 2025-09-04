package cql

import (
	"errors"

	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/model"
)

type ValueIntoSelection[TValue any, TResults any] struct {
	value    condition.ValueOfType[TValue]
	selector func(TValue, *TResults)
}

// TODO probablemente ya no sea necesaria
func (selection ValueIntoSelection[TValue, TResults]) ResultsType() TResults {
	return *new(TResults)
}

func (selection ValueIntoSelection[TValue, TResults]) Apply(value any, result *TResults) error {
	valueT, isT := value.(TValue)
	if !isT {
		// TODO definir bien el error
		return errors.New("not possible error")
	}

	selection.selector(valueT, result)

	return nil
}

func ValueInto[TValue any, TResults any](
	value condition.ValueOfType[TValue],
	selector func(TValue, *TResults),
) *ValueIntoSelection[TValue, TResults] {
	return &ValueIntoSelection[TValue, TResults]{
		value:    value,
		selector: selector,
	}
}

// TODO docs
func Select[TResults any, TModel model.Model](
	query *condition.Query[TModel],
	selections ...condition.Selection[TResults],
) ([]TResults, error) {
	return condition.Select(
		query,
		selections,
	)
}
