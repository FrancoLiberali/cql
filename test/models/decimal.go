package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/FrancoLiberali/cql/model"
)

var errUnsupportedCentsScan = errors.New("unsupported scan source type for Cents")

// Cents is a manually-defined decimal type — an integer number of cents, with
// no external library. It is a driver.Valuer / sql.Scanner that maps to a
// NUMERIC column, so cql-gen recognizes it as a decimal purely from the
// column's `type:numeric` tag and exposes it as a DecimalField.
type Cents int64

func (c Cents) Value() (driver.Value, error) { return int64(c), nil }

func (c *Cents) Scan(src any) error {
	switch v := src.(type) {
	case int64:
		*c = Cents(v)
	case int:
		*c = Cents(v)
	case []byte:
		// postgres/mysql deliver NUMERIC as text.
		return c.scanText(string(v))
	case string:
		return c.scanText(v)
	default:
		return fmt.Errorf("%w: %T", errUnsupportedCentsScan, src)
	}

	return nil
}

func (c *Cents) scanText(s string) error {
	// A DECIMAL column can come back with a fractional part — e.g. mysql renders
	// "balance + 5" as "15.000...0". Cents is integer-valued, so accept an
	// all-zero fraction and reject a real one.
	if dot := strings.IndexByte(s, '.'); dot >= 0 {
		if strings.Trim(s[dot+1:], "0") != "" {
			return fmt.Errorf("%w: non-integer %q", errUnsupportedCentsScan, s)
		}

		s = s[:dot]
	}

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("%w: %q", errUnsupportedCentsScan, s)
	}

	*c = Cents(n)

	return nil
}

// Bank exercises DecimalField with both a manually-defined decimal type
// (Balance, type Cents) and the shopspring/decimal library (Rate). Both columns
// are tagged decimal/numeric, so cql-gen emits DecimalField for each.
type Bank struct {
	model.UUIDModel

	Balance Cents           `gorm:"type:numeric(19,0)"`
	Rate    decimal.Decimal `gorm:"type:decimal(19,4)"`
}

func (m Bank) Equal(other Bank) bool {
	return m.ID == other.ID
}
