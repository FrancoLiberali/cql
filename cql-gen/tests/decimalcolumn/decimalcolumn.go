package decimalcolumn

import (
	"database/sql/driver"
	"fmt"

	"github.com/FrancoLiberali/cql/model"
)

// Money is a fixed-point decimal money type stored in a DECIMAL column. It is a
// plain Scanner/Valuer with NO Go-side arithmetic — cql-gen recognizes it as a
// decimal from the column's `type:decimal` tag, not from its methods.
type Money int64

func (m Money) Value() (driver.Value, error) { return int64(m), nil }

func (m *Money) Scan(src any) error {
	v, ok := src.(int64)
	if !ok {
		return fmt.Errorf("cannot scan %T into Money", src)
	}

	*m = Money(v)

	return nil
}

// Note is a custom Scanner/Valuer type WITHOUT a decimal
// tag: it must NOT become a DecimalField — it stays a plain (Updatable) Field.
type Note string

func (n Note) Value() (driver.Value, error) { return string(n), nil }

func (n *Note) Scan(src any) error {
	v, ok := src.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into Note", src)
	}

	*n = Note(v)

	return nil
}

// Wallet has a decimal column (Balance, tagged type:decimal) alongside a
// non-decimal custom column (Note) to exercise both the positive and negative
// detection paths.
type Wallet struct {
	model.UUIDModel

	Balance Money `gorm:"type:decimal(19,4)"`
	Note    Note
}
