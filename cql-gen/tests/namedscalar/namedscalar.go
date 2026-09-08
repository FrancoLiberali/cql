package namedscalar

import "github.com/FrancoLiberali/cql/model"

// Color is a named scalar type: no custom Scan/Value, so gorm stores it as its
// underlying int. cql-gen classifies it via its underlying kind, so it takes
// the fast-scan path (scan into NullInt64, cast back to Color) and gets a
// condition field typed by the named type (NumericField[..., Color]).
type Color int

// WithNamedScalar has a Color column as a value and as a nullable pointer, so
// cql-gen must generate both conditions and a scanner covering both the value
// and pointer named-scalar paths.
type WithNamedScalar struct {
	model.UUIDModel

	Name     string
	Favorite Color
	Second   *Color
}
