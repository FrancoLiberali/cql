package cql

import (
	"errors"

	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/model"
)

var errValueTypeIsNotExpectedType = errors.New("type of value is not the expected type")

type ValueIntoSelection[TValue any, TResults any] struct {
	value    condition.ValueOfType[TValue]
	selector func(TValue, *TResults)
}

func (selection ValueIntoSelection[TValue, TResults]) Apply(value any, result *TResults) error {
	valueT, isTPointer := value.(*TValue)
	if !isTPointer {
		return errValueTypeIsNotExpectedType
	}

	selection.selector(*valueT, result)

	return nil
}

func (selection ValueIntoSelection[TValue, TResults]) ValueType() any {
	return new(TValue)
}

func (selection ValueIntoSelection[TValue, TResults]) ToSQL(query *condition.CQLQuery) (string, []any, error) {
	return selection.value.ToSQL(query)
}

// ValueInto allows the definition of a selection of a value into an attribute of the results
//
// For example, to select sale.Code into result.Code:
//
//	cql.ValueInto(conditions.Sale.Code, func(value float64, result *Result) {
//		result.Code = int(value)
//	})
func ValueInto[TValue any, TResults any](
	value condition.ValueOfType[TValue],
	selector func(TValue, *TResults),
) *ValueIntoSelection[TValue, TResults] {
	return &ValueIntoSelection[TValue, TResults]{
		value:    value,
		selector: selector,
	}
}

// Select specify fields that you want when querying.
//
// # Use Select when you only want a subset of the fields, not all the fields of a model
//
// Use cql.ValueInto to generate the selections, for example:
//
//	// Select only sale.Code into a []Result
//	results, err := cql.Select(
//		cql.Query[models.Sale](ts.db),
//		cql.ValueInto(conditions.Sale.Code, func(value float64, result *Result) {
//			result.Code = int(value)
//		}),
//	)
func Select[TResults any, TModel model.Model](
	query *condition.Query[TModel],
	selections ...condition.Selection[TResults],
) ([]TResults, error) {
	return condition.Select(
		query,
		selections,
	)
}
