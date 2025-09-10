package cql

import (
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/sql"
)

// CountAll is an Aggregation that returns the number of values (including nulls)
//
// Example:
//
// cql.Query[models.Product](db).GroupBy(conditions.Product.Int).SelectValue(cql.CountAll(), "aggregation").Into(&results)
func CountAll() condition.AggregationResult[float64] {
	return condition.AggregationResult[float64]{
		Function: sql.CountAll,
	}
}
