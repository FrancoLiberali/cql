package goembeddedscan

import "github.com/FrancoLiberali/cql/model"

// Inner is embedded anonymously (Go-style promotion, no gorm:"embedded" tag).
// gorm maps its fields to columns with no prefix, so they must not collide with
// the parent's own columns.
type Inner struct {
	InnerInt int
}

// GoEmbeddedScan embeds Inner anonymously with no column collision, so the
// scanner generator emits a complete switch that threads the Go access path to
// the promoted field. Positive companion to TestGoEmbedded, which covers the
// negative case (a promoted column colliding with a top-level one).
type GoEmbeddedScan struct {
	model.UIntModel
	Inner

	Top int
}
