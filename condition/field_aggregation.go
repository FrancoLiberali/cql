package condition

import (
	"fmt"

	"github.com/FrancoLiberali/cql/sql"
)

type FieldAggregation struct {
	field IField
}

type NumericFieldAggregation struct {
	FieldAggregation
}

// Sum TODO
func (fieldAggregation NumericFieldAggregation) Sum() Aggregation {
	return Aggregation{
		field:    fieldAggregation.field,
		function: sql.Sum,
	}
}

type Aggregation struct {
	field    IField
	function sql.FunctionByDialector
}

func (aggregation Aggregation) toSQL(query *GormQuery, table Table, as string) (string, error) {
	function, isPresent := aggregation.function.Get(query.Dialector())
	if !isPresent {
		return "", functionError(ErrUnsupportedByDatabase, aggregation.function)
	}

	return fmt.Sprintf(
		"%s AS %s",
		function.ApplyTo(aggregation.field.columnSQL(query, table), 0),
		as,
	), nil
}
