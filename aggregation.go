package cql

import (
	"github.com/FrancoLiberali/cql/condition"
	"github.com/FrancoLiberali/cql/sql"
)

// CountAll is an Aggregation that returns the number of values (including nulls)
//
// Example:
//
// cql.Query[models.Product](db).GroupBy(conditions.Product.Int).Select(cql.CountAll(), "aggregation").Into(&results)
func CountAll() condition.Aggregation {
	return condition.Aggregation{
		Function: sql.CountAll,
	}
}
