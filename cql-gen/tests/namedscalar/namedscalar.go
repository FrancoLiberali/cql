package namedscalar

import "github.com/FrancoLiberali/cql/model"

// Color is a named scalar type: valid SQL (gorm stores it as an int) but not
// one the fast scanner can classify — it's not a model ID, time.Time, sql.Null*
// wrapper, or gorm custom type.
type Color int

// WithNamedScalar has a Color column the fast scanner can't handle. cql-gen must
// generate conditions for it but NOT a scanner, so its queries fall back to
// gorm's reflective scan rather than silently dropping the column.
type WithNamedScalar struct {
	model.UUIDModel

	Name     string
	Favorite Color
}
