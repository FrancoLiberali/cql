package unsupportedcolumn

import "github.com/FrancoLiberali/cql/model"

// WithUnsupportedColumn has a []string persisted via gorm's json serializer:
// there's no static Scan/Value and the underlying type is a slice, so the fast
// scanner can't classify it. cql-gen must still generate conditions for it but
// NOT a scanner, so its queries fall back to gorm rather than silently drop the
// column.
type WithUnsupportedColumn struct {
	model.UUIDModel

	Name string
	Tags []string `gorm:"serializer:json"`
}
