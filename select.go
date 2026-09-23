package cql

import (
	"github.com/FrancoLiberali/cql/condition"
)

// Select specify fields that you want when querying.
//
// # Use Select when you only want a subset of the fields, not all the fields of a model
//
// Use the Into method of a field or aggregation to generate the selections,
// binding each selected value straight into an attribute of the results. For
// example, to select sale.Code into result.Code:
//
//	results, err := cql.Select(
//		cql.Query[models.Sale](ctx, db),
//		conditions.Sale.Code.Into(func(result *Result) *int { return &result.Code }),
//	)
func Select[TResults any](
	query condition.IQuery,
	selections ...condition.Selection[TResults],
) ([]TResults, error) {
	return condition.Select(
		query,
		selections,
	)
}
