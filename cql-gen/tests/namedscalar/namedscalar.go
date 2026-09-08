package namedscalar

import "github.com/FrancoLiberali/cql/model"

// Color / Mood / Switch are named scalar types over int / string / bool: none
// has a custom Scan/Value, so gorm stores each as its underlying kind. cql-gen
// classifies them via the underlying kind (scan into the matching sql.Null*
// wrapper, cast back to the named type) and gives each a condition field typed
// by the named type where the field kind allows it.
type Color int

type Mood string

type Switch bool

// WithNamedScalar has each named scalar as a value and as a nullable pointer,
// so cql-gen must generate conditions and a scanner covering the int, string
// and bool named-scalar paths in both the value and pointer forms.
type WithNamedScalar struct {
	model.UUIDModel

	Name string

	Favorite Color
	Second   *Color

	Mood      Mood
	AltMood   *Mood
	Active    Switch
	AltActive *Switch
}
