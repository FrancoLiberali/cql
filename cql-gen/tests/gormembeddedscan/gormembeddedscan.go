package gormembeddedscan

import "github.com/FrancoLiberali/cql/model"

// Inner is reused by two embedded fields below to verify the scanner
// generator threads the Go access path correctly (dest.A.Int vs dest.B.Int).
type Inner struct {
	Int int
}

// GormEmbeddedScan exercises the gorm:"embedded" tag path. Both embedded
// fields use distinct embeddedPrefix values so columns do NOT collide,
// which lets the scanner generator emit a complete switch with
// dest.<Parent>.<Inner> assignments. Companion to TestGormEmbedded, which
// covers the negative case (colliding columns triggering DuplicateColumnError).
type GormEmbeddedScan struct {
	model.UIntModel

	Top  int
	Foo  Inner `gorm:"embedded;embeddedPrefix:foo_"`
	Bar  Inner `gorm:"embedded;embeddedPrefix:bar_"`
}
