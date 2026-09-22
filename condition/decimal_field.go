package condition

import (
	"database/sql/driver"

	"github.com/FrancoLiberali/cql/model"
	"github.com/FrancoLiberali/cql/sql"
)

// Decimal is the constraint for the Go type backing a decimal / NUMERIC column
// — e.g. github.com/shopspring/decimal.Decimal, cockroachdb/apd, or a
// fixed-point money type. It only requires driver.Valuer, so an operand value
// can be bound to SQL; reading back is handled by *TDecimal implementing
// sql.Scanner at runtime, as for any custom column type.
//
// It is intentionally loose: DecimalField is only ever emitted by cql-gen,
// which decides a column is a decimal from its DB type (a `type:decimal` /
// `type:numeric` gorm tag, or a known-decimal-type allowlist) — NOT from the Go
// type's shape. So the constraint does not need to prove "is a decimal"; that
// guarantee lives in the generator. This keeps any decimal library usable with
// no wrapper, regardless of its arithmetic method names — CQL evaluates the
// arithmetic in SQL and never calls Go-side Add/Sub/Mul/Div.
//
// (Structs can never satisfy the Numeric constraint — constraints.Integer |
// constraints.Float admits only basic types — which is why decimal columns
// cannot use NumericField and get their own field type here.)
type Decimal interface {
	driver.Valuer
}

// DecimalField identifies a decimal / NUMERIC column.
//
// Unlike NumericField it never goes through float64: decimal arithmetic is
// closed (decimal ∘ decimal = decimal) and its aggregations keep the decimal
// type, so results are exact — no 2^53 rounding. Comparisons (Is), updates
// (Set) and exact selection (Into) are inherited unchanged from Field, which is
// already generic over any attribute type; only the arithmetic and Sum/Average
// aggregations — the pieces the Numeric constraint used to gate — are added.
type DecimalField[TModel model.Model, TDecimal Decimal] struct {
	UpdatableField[TModel, TDecimal]
}

func NewDecimalField[TModel model.Model, TDecimal Decimal](
	name, column, columnPrefix string,
) DecimalField[TModel, TDecimal] {
	return DecimalField[TModel, TDecimal]{
		UpdatableField: NewUpdatableField[TModel, TDecimal](name, column, columnPrefix),
	}
}

// Aggregate allows applying decimal-preserving aggregation functions inside a
// group by. The result stays TDecimal, so Into binds it exactly.
func (field DecimalField[TModel, TDecimal]) Aggregate() DecimalFieldAggregation[TDecimal] {
	return DecimalFieldAggregation[TDecimal]{field: field}
}

// Plus adds other to the value. The operand is ValueOfType[TDecimal], so the
// arithmetic is type-closed: a decimal can only be combined with another
// decimal of the same type, never silently widened to float.
func (field DecimalField[TModel, TDecimal]) Plus(other ValueOfType[TDecimal]) NotUpdatableDecimalField[TModel, TDecimal] {
	return newNotUpdatableDecimalField(field.Field.addFunction(sql.Plus, other))
}

// Minus subtracts other from the value.
func (field DecimalField[TModel, TDecimal]) Minus(other ValueOfType[TDecimal]) NotUpdatableDecimalField[TModel, TDecimal] {
	return newNotUpdatableDecimalField(field.Field.addFunction(sql.Minus, other))
}

// Times multiplies the value by other.
func (field DecimalField[TModel, TDecimal]) Times(other ValueOfType[TDecimal]) NotUpdatableDecimalField[TModel, TDecimal] {
	return newNotUpdatableDecimalField(field.Field.addFunction(sql.Times, other))
}

// Divided divides the value by other.
func (field DecimalField[TModel, TDecimal]) Divided(other ValueOfType[TDecimal]) NotUpdatableDecimalField[TModel, TDecimal] {
	return newNotUpdatableDecimalField(field.Field.addFunction(sql.Divided, other))
}

// NotUpdatableDecimalField is a decimal arithmetic expression. It is not
// settable (you cannot UPDATE an expression), but it keeps the decimal type so
// it can be aggregated, combined further, or selected exactly with Into.
type NotUpdatableDecimalField[TModel model.Model, TDecimal Decimal] struct {
	Field[TModel, TDecimal]
}

func newNotUpdatableDecimalField[TModel model.Model, TDecimal Decimal](
	field Field[TModel, TDecimal],
) NotUpdatableDecimalField[TModel, TDecimal] {
	return NotUpdatableDecimalField[TModel, TDecimal]{Field: field}
}

func (field NotUpdatableDecimalField[TModel, TDecimal]) Plus(other ValueOfType[TDecimal]) NotUpdatableDecimalField[TModel, TDecimal] {
	return newNotUpdatableDecimalField(field.Field.addFunction(sql.Plus, other))
}

func (field NotUpdatableDecimalField[TModel, TDecimal]) Minus(other ValueOfType[TDecimal]) NotUpdatableDecimalField[TModel, TDecimal] {
	return newNotUpdatableDecimalField(field.Field.addFunction(sql.Minus, other))
}

func (field NotUpdatableDecimalField[TModel, TDecimal]) Times(other ValueOfType[TDecimal]) NotUpdatableDecimalField[TModel, TDecimal] {
	return newNotUpdatableDecimalField(field.Field.addFunction(sql.Times, other))
}

func (field NotUpdatableDecimalField[TModel, TDecimal]) Divided(other ValueOfType[TDecimal]) NotUpdatableDecimalField[TModel, TDecimal] {
	return newNotUpdatableDecimalField(field.Field.addFunction(sql.Divided, other))
}

// Aggregate allows aggregating the arithmetic expression, still exactly.
func (field NotUpdatableDecimalField[TModel, TDecimal]) Aggregate() DecimalFieldAggregation[TDecimal] {
	return DecimalFieldAggregation[TDecimal]{field: field}
}

// DecimalFieldAggregation exposes the decimal-preserving aggregations. Every
// result is an AggregationResult[TDecimal], so AggregationResult.Into binds the
// exact decimal type with no float64 in between.
type DecimalFieldAggregation[TDecimal Decimal] struct {
	field IField
}

// Sum calculates the exact summation of all values.
func (aggregation DecimalFieldAggregation[TDecimal]) Sum() AggregationResult[TDecimal] {
	return AggregationResult[TDecimal]{field: aggregation.field, Function: sql.Sum}
}

// Average calculates the exact average of all values.
func (aggregation DecimalFieldAggregation[TDecimal]) Average() AggregationResult[TDecimal] {
	return AggregationResult[TDecimal]{field: aggregation.field, Function: sql.Average}
}

// Min returns the minimum value of all values.
func (aggregation DecimalFieldAggregation[TDecimal]) Min() AggregationResult[TDecimal] {
	return AggregationResult[TDecimal]{field: aggregation.field, Function: sql.Min}
}

// Max returns the maximum value of all values.
func (aggregation DecimalFieldAggregation[TDecimal]) Max() AggregationResult[TDecimal] {
	return AggregationResult[TDecimal]{field: aggregation.field, Function: sql.Max}
}
