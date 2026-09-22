package condition

import "errors"

// errIntoTypeMismatch is returned if the scanned value is not the strictly
// typed destination the Into selector expects. In practice this cannot happen
// through the generated API — Into is fully typed at compile time — so it only
// guards against a hand-built Selection wired incorrectly.
var errIntoTypeMismatch = errors.New("into: scanned value is not of the expected type")

// intoSelection is the strict, pointer-form Selection produced by the Into
// methods below.
//
// Unlike ValueIntoSelection (the closure form, func(value, *result)), it binds
// the selected value straight into a field of the per-row result through
// selector, which returns the ADDRESS of that field. The destination must be
// exactly TValue, so no conversion is ever needed: the database driver scans
// into a TValue scratch of the correct type — an *int for an int column, not
// the float64 the numeric value abstraction would otherwise impose — and it is
// copied into the destination.
type intoSelection[TValue any, TResults any] struct {
	value    IValue
	selector func(*TResults) *TValue
}

// ValueType returns a scratch of the destination's exact type, so the driver
// scans directly into it (no float64 detour).
func (s intoSelection[TValue, TResults]) ValueType() any {
	return new(TValue)
}

// Apply copies the strictly-typed scratch into the destination field.
func (s intoSelection[TValue, TResults]) Apply(scanned any, result *TResults) error {
	typed, ok := scanned.(*TValue)
	if !ok {
		return errIntoTypeMismatch
	}

	*s.selector(result) = *typed

	return nil
}

func (s intoSelection[TValue, TResults]) ToSQL(query *CQLQuery) (string, []any, error) {
	return s.value.ToSQL(query)
}

// Into selects this field into the result of a query, addressing the
// destination field through selector.
//
// It is strictly typed: the destination must be *TAttribute — the field's own
// Go type — so a plain int column lands in an int with no float64 detour and
// no int(value) conversion at the call site:
//
//	results, err := cql.Select(
//		cql.Query[models.Product](ctx, db),
//		conditions.Product.Int.Into(func(r *Result) *int { return &r.Int }),
//	)
//
// (With Go 1.27 this is a generic method: TResults is the method's own type
// parameter, which was impossible before.)
func (field Field[TModel, TAttribute]) Into[TResults any](
	selector func(*TResults) *TAttribute,
) Selection[TResults] {
	return intoSelection[TAttribute, TResults]{value: field, selector: selector}
}

// Into selects an arithmetic expression (the result of Plus/Minus/Times/
// Divided/...) into the result of a query.
//
// The destination is *float64, not the column's own type, because arithmetic
// follows SQL's type promotion: int + float is a float, and int / int can be
// fractional. float64 is SQL's safe superset for a mixed numeric expression,
// so this both matches the database and lets int.Plus(aFloat) work instead of
// failing when the row is scanned. (A plain column keeps its real type through
// Field.Into above; only expressions widen.)
//
// PRECISION: float64 represents integers exactly only up to 2^53
// (9,007,199,254,740,992), well below the int64 range (~9.2e18). An expression
// whose true integer value exceeds 2^53 is therefore rounded. If you need exact
// results for large integers, select the underlying column(s) with Field.Into
// (which keeps the real integer type) and do the arithmetic in Go, or model the
// column as a decimal type.
func (field NotUpdatableNumericField[TModel, TAttribute]) Into[TResults any](
	selector func(*TResults) *float64,
) Selection[TResults] {
	return intoSelection[float64, TResults]{value: field, selector: selector}
}

// Into selects this aggregation into the result of a query. The destination
// type is the aggregation's honest result type — float64 for Sum/Average of a
// numeric column, because SQL itself produces a fractional/widened value there;
// there is nothing to convert away.
//
// PRECISION: for numeric aggregations the destination is float64, exact for
// integer results only up to 2^53 (9,007,199,254,740,992). A SUM of large
// integers beyond that is rounded. There is no exact-typed alternative here yet
// (see docs); if you need it, aggregate in the database into a decimal column
// or fetch the rows and reduce in Go.
func (aggregation AggregationResult[T]) Into[TResults any](
	selector func(*TResults) *T,
) Selection[TResults] {
	return intoSelection[T, TResults]{value: aggregation, selector: selector}
}
